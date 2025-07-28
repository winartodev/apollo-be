package main

import "github.com/winartodev/apollo-be/core/config"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := cfg.Database.SetupConnection()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	autoMigration, err := config.NewAutoMigration(cfg.Database.Name, db)
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
