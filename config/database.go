package config

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	postgresDNS = "postgres://%s:%s@%s:%s/%s?sslmode=%s"

	errorOpenDBConnection = "error open db connection %v"
	errorPingDB           = "error ping db %v"
	errorValidateConnDB   = "error validate conn db %v"
)

type Database struct {
	Driver          string `yaml:"driver"`
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	Name            string `yaml:"name"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	SSLMode         string `yaml:"sslMode"`
	MaxOpenConn     int    `yaml:"defaultMaxConn"`
	MaxIdleConn     int    `yaml:"defaultIdleConn"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime"`
	ConnMaxIdleTime int    `yaml:"connMaxIdleTime"`
}

func (d *Database) SetupConnection() (*sql.DB, error) {
	dsn := fmt.Sprintf(postgresDNS, d.Username, d.Password, d.Host, d.Port, d.Name, d.SSLMode)
	db, err := sql.Open(d.Driver, dsn)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(errorOpenDBConnection, err))
	}

	db.SetConnMaxLifetime(time.Duration(d.ConnMaxLifetime) * time.Second)
	db.SetMaxIdleConns(d.MaxIdleConn)
	db.SetMaxOpenConns(d.MaxOpenConn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(d.ConnMaxLifetime)*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, errors.New(fmt.Sprintf(errorPingDB, err))
	}

	if err := validateConnectionPool(ctx, db); err != nil {
		return nil, errors.New(fmt.Sprintf(errorValidateConnDB, err))
	}

	return db, nil
}

func validateConnectionPool(ctx context.Context, db *sql.DB) error {
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}

	defer conn.Close()

	// Execute a simple query to verify connection works
	var result int
	err = conn.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return err
	}

	if result != 1 {
		return errors.New("connection test query failed")
	}

	return nil
}
