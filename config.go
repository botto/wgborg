package main

import (
	"log"

	"github.com/caarlos0/env/v6"
)

var config *configDef

type configDef struct {
	DBHost     string `env:"DBHOST" envDefault:"127.0.0.1"`
	DBPort     uint32 `env:"DBPORT" envDefault:"5432"`
	DBUser     string `env:"DBUSER" envDefault:"wgmgr"`
	DBPassword string `env:"DBPASSWORD" envDefault:"password"`
	DBName     string `env:"DBNAME" envDefault:"wgmgr"`
}

func initConfig() {
	config = &configDef{}
	err := env.Parse(config)
	if err != nil {
		log.Fatalf("Could not parse config: %s\n", err)
	}
}
