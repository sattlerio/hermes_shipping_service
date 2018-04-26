package config

import "os"

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	Dialect 	string
	Username 	string
	Password	string
	Name		string
	Charset 	string
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Dialect: "postgres",
			Username: os.Getenv("PSQL_USER"),
			Password: os.Getenv("PSQL_PASSWORD"),
			Name: os.Getenv("PSQL_DB"),
			Charset: "utf-8",
		},
	}
}
