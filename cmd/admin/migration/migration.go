package adminmigration

import (
	"context"
	"fmt"
	appconfig "goapp/config"
	"goapp/internal/store"
	migrationentity "goapp/internal/store/migration"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	SQLMigrationDir = "./internal/store/migration/sql"
)

func Migration(cnf appconfig.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Migrate database",
		RunE: func(cmd *cobra.Command, args []string) error {
			files, err := os.ReadDir(SQLMigrationDir)
			if err != nil {
				return errors.Wrap(err, "failed to read sql migration directory")
			}

			newStore := store.NewStore(cnf)

			type Mig struct {
				Name   string
				SqlQry string
			}
			var queries []Mig
			var migIndexTable []Mig

			var exists bool
			migDefTable := "db_migrations"
			migDefFileName := "4715896587__migrations.sql"

			err = newStore.Postgres.DB.QueryRow(context.Background(), `SELECT EXISTS ( SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename  = $1 )`, migDefTable).Scan(&exists)
			if err != nil {
				return errors.Wrap(err, "failed to check migration table exists")
			}

			if len(files) > 0 {
				var migrations []migrationentity.MigrationModels
				if exists {
					migrations, err = newStore.Postgres.Migration.GetAll(context.Background())
					if err != nil {
						return errors.Wrap(err, "failed to get all migrations")
					}
				}
				for _, f := range files {
					ext := filepath.Ext(f.Name())
					if ext == ".sql" {
						found := false
						if len(migrations) > 0 {
							for _, a := range migrations {
								if a.Name == f.Name() {
									found = true
									break
								}
							}
						}

						if !found {
							filePath := fmt.Sprintf("%s/%s", SQLMigrationDir, f.Name())
							query, err := os.ReadFile(filePath)
							if err != nil {
								return errors.Wrap(err, "failed reading sql migration")
							}
							if len(query) > 0 {
								item := Mig{
									Name:   f.Name(),
									SqlQry: string(query),
								}
								if !exists && f.Name() == migDefFileName {
									migIndexTable = append(migIndexTable, item)
								} else {
									queries = append(queries, item)
								}
							}
						}
					}
				}
			}

			runQueries := func(newStore *store.Base, queryItems []Mig) error {
				for _, q := range queryItems {
					created := newStore.Postgres.Migration.CreateOneByOne(context.Background(), q.SqlQry)
					if created {
						mig, err := newStore.Postgres.Migration.Insert(context.Background(), migrationentity.MigrationModels{
							Name:      q.Name,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						})
						if err != nil {
							return errors.Wrap(err, "failed postgres sql migration insertions")
						}
						fmt.Println("postgres sql migration insertions completed", "migration", mig.Name)
					} else {
						fmt.Println("postgres sql migration failed", "query", q.Name)
					}
				}
				fmt.Println("postgres sql migration completed")
				return nil
			}
			if len(migIndexTable) > 0 {
				fmt.Println("migrations setup started")
				err := runQueries(newStore, migIndexTable)
				if err != nil {
					return errors.Wrap(err, "failed postgres sql migration insertions")
				}
				fmt.Println("migrations setup completed\n\r")
			}
			if len(queries) > 0 {
				fmt.Println("db migrations started")
				err := runQueries(newStore, queries)
				if err != nil {
					return errors.Wrap(err, "failed postgres sql migration insertions")
				}
				fmt.Println("db migrations completed\n\r")
			}
			if len(migIndexTable) <= 0 && len(queries) <= 0 {
				fmt.Println("all migrations are up to date")
			}
			os.Exit(1)
			return nil
		},
	}
}

func CreateMigration(cnf appconfig.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "create_migration",
		Short: "Generate empty migration file with the give <file_name>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fname := strconv.FormatInt(time.Now().Unix(), 10) + "__" + args[0] + ".sql"
			filepath := fmt.Sprintf("%s/%s", SQLMigrationDir, fname)
			file, err := os.Create(filepath)
			if err != nil {
				return errors.Wrap(err, "failed to create sql migration file")
			}
			err = file.Chmod(0777)
			if err != nil {
				return errors.Wrap(err, "failed to change sql migration file permission")
			}
			defer func(file *os.File) {
				if err := file.Close(); err != nil {
					fmt.Println("failed to close sql migration file", "err", err)
				}
			}(file)
			fmt.Println("sql migration file created", "file_path", filepath, "err", err)
			return nil
		},
	}
}
