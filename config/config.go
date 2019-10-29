package config

import "os"

type dbConfig struct {
	DatabaseHost     string
	DatabaseUsername string
	DatabasePassword string
}

//DBConfig config
var DBConfig dbConfig

//Load app config
func Load() {
	DBConfig = dbConfig{
		DatabaseHost:     os.Getenv("DB_HOST"),
		DatabaseUsername: os.Getenv("DB_USERNAME"),
		DatabasePassword: os.Getenv("DB_PASSWORD"),
	}
}
