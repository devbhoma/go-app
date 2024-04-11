package appconfig

import (
	"github.com/joho/godotenv"
	"os"
)

type PgDatabase struct {
	Driver string `json:"driver"`
	Host   string `json:"host"`
	Port   string `json:"port"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Name   string `json:"name"`
}

type Config struct {
	Env      string `json:"env"`
	Port     string `json:"port"`
	Database PgDatabase
}

func Get() Config {

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	cnf := Config{
		Port: os.Getenv("APP_PORT"),
		Env:  os.Getenv("APP_ENV"),
		Database: PgDatabase{
			Driver: os.Getenv("DB_DRIVER"),
			Host:   os.Getenv("DB_HOST"),
			Port:   os.Getenv("DB_PORT"),
			User:   os.Getenv("DB_USER"),
			Pass:   os.Getenv("DB_PASS"),
			Name:   os.Getenv("DB_NAME"),
		},
	}

	return cnf
}
