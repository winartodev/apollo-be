package main

import (
	config2 "github.com/winartodev/apollo-be/config"
)

func main() {
	cfg, err := config2.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := cfg.Database.SetupConnection()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	autoMigration, err := config2.NewAutoMigration(cfg.Database.Name, db)
	if err != nil {
		panic(err)
	}

	if autoMigration == nil {
		panic("autoMigration is nil")
	}

	if err := autoMigration.Start(); err != nil {
		panic(err)
	}
}
