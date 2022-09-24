package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
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
		db.SetConnMaxLifetime(5 * time.Minute)
		myDB = &DB{db}
	})

	return myDB
}

func getConnectionString(cfg DBConfig) (string, error) {
	if cfg.Server == "" || cfg.Port == "" || cfg.DbName == "" || cfg.UserID == "" || cfg.Password == "" {
		return "", fmt.Errorf("getConnectionString: One or more required db configuration  missing: %v", cfg)
	}

	myCfg := mysql.Config{
		User:                 cfg.UserID,
		Passwd:               cfg.Password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", cfg.Server, cfg.Port),
		DBName:               cfg.DbName,
		AllowNativePasswords: true,
	}

	s := myCfg.FormatDSN()
	fmt.Printf("dsn: %s\n", s)

	return s, nil
}

func (d DB) FindAny(ctx context.Context, objects interface{}, query string, args []interface{}) error {
	var stmt *sqlx.Stmt

	stmt, err := d.db.PreparexContext(ctx, query)
	if err != nil {
		return err
	}

	return stmt.Select(objects, args...)
}
