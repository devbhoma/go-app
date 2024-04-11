package coreadmin

import (
	adminmigration "goapp/cmd/admin/migration"
	"goapp/cmd/cli"
	appconfig "goapp/config"
)

func Boot(cli *cli.Base, cnf appconfig.Config) {

	cli.AddCommand(adminmigration.Migration(cnf))
	cli.AddCommand(adminmigration.CreateMigration(cnf))

}
