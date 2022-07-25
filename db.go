package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	dbOnce sync.Once
	myDB   *DB
)

type DBConfig struct {
	Server   string
	Port     string
	DbName   string
	UserID   string
	Password string
}

type DB struct {
	db *sqlx.DB
}

func GetDB() *DB {
	return myDB
}

func LoadAuthDB(cfg DBConfig) *DB {
	dbOnce.Do(func() {
		datasource, err := getConnectionString(cfg)
		if err != nil {
			log.Fatal(err)
		}

		db, err := sqlx.Connect("mysql", datasource)
		if err != nil {
			log.Fatal(err)
		}
		myDB = &DB{db}
	})

	return myDB
}

func getConnectionString(cfg DBConfig) (string, error) {
	if cfg.Server == "" || cfg.Port == "" || cfg.DbName == "" || cfg.UserID == "" || cfg.Password == "" {
		return "", fmt.Errorf("getConnectionString: One or more required db configuration  missing")
	}

	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.UserID, cfg.Password, cfg.Server, cfg.Port, cfg.DbName)
	return connString, nil
}

func (d DB) FindAny(ctx context.Context, objects interface{}, query string, args []interface{}) error {
	var stmt *sqlx.Stmt

	stmt, err := d.db.PreparexContext(ctx, query)
	if err != nil {
		return err
	}

	return stmt.Select(objects, args...)
}
