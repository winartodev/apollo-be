package config

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	pg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/winartodev/apollo-be/core/helper"
)

const (
	migrationsDir = "core/migrations"
	osWindows     = "windows"
)

type SchemaMigration struct {
	version int
	dirty   bool
}

type AutoMigration struct {
	migrate *migrate.Migrate
}

func NewAutoMigration(databaseName string, db *sql.DB) (*AutoMigration, error) {
	sourceURL, err := generateSourceURL()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to generate source URL: %v", err))
	}

	x, err := pg.WithInstance(db, &pg.Config{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create database driver instance: %v", err))
	}

	m, err := migrate.NewWithDatabaseInstance(*sourceURL, databaseName, x)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to get schema version after migration: %v", err))
	}

	log.Printf("AutoMigration initialized successfully for database: %s", databaseName)
	return &AutoMigration{
		migrate: m,
	}, nil
}

func (am *AutoMigration) Start() error {
	schemaUpErr := am.migrate.Up()
	if schemaUpErr == nil {
		schema, err := am.getSchemaMigration()
		if err != nil {
			return errors.New(fmt.Sprintf("failed to get schema version after migration error: %v", err))
		}

		log.Printf("Migration successful. Current schema version: %d", schema.version)
		return nil
	}

	if isErrorNoChange(schemaUpErr) {
		log.Printf("No migrations to run. Database is already at the latest version.")
		return nil
	}

	// Handle if migrate failed
	schema, err := am.getSchemaMigration()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get schema version after migration error: %v", err))
	}

	log.Printf("Migration failed error {%v}.\nCurrent schema version: %d, dirty: %v", schemaUpErr.Error(), schema.version, schema.dirty)

	err = am.Fix(schema.version)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to fix schema version: %v", err))
	}

	err = am.Rollback()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to rollback migration: %v", err))
	}

	log.Printf("Migration recovery completed successfully.")

	return nil
}

func (am *AutoMigration) Fix(version int) error {
	log.Printf("Fixing schema version to: %d", version)

	err := am.migrate.Force(version)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to force schema version: %v", err))
	}

	schema, err := am.getSchemaMigration()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get schema version after fix: %v", err))
	}

	log.Printf("Schema fix successful. Current version: %d, dirty: %v", schema.version, schema.dirty)
	return nil
}

func (am *AutoMigration) Rollback() error {
	log.Printf("Rolling back schema by 1 step.")

	err := am.migrate.Steps(-1)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to rollback schema: %v", err))
	}

	schema, err := am.getSchemaMigration()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get schema version after rollback: %v", err))
	}

	log.Printf("Rollback successful. Current version: %d, dirty: %v", schema.version, schema.dirty)

	return nil
}

func (am *AutoMigration) getSchemaMigration() (*SchemaMigration, error) {
	version, dirty, err := am.migrate.Version()
	if isNonNilAndNotExpectedMigrationError(err) {
		return nil, errors.New(fmt.Sprintf("failed to retrieve schema version: %v", err))
	}

	return &SchemaMigration{
		version: int(version),
		dirty:   dirty,
	}, nil
}

func generateSourceURL() (*string, error) {
	filePath, err := helper.GetCompletePath(migrationsDir)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	var completePath string
	if helper.CurrentOS(osWindows) {
		completePath = fmt.Sprintf("file:%s", filePath)
	} else {
		completePath = fmt.Sprintf("file://%s", filePath)
	}

	return &completePath, nil
}

func isNonNilAndNotExpectedMigrationError(err error) bool {
	return err != nil && isErrorNoMigration(err) && isErrorNoChange(err)
}

func isErrorNoChange(err error) bool {
	return err != nil && errors.Is(err, migrate.ErrNoChange)
}

func isErrorNoMigration(err error) bool {
	return err != nil && errors.Is(err, migrate.ErrNilVersion)
}
