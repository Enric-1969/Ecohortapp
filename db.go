package main

import (
	"database/sql"
	"ecohortapp/repository"
	"log"
	"os"

	_ "github.com/glebarez/go-sqlite"
)

// connectSQL gestiona la connexió amb la base de dades SQLite
func (app *Config) connectSQL() (*sql.DB, error) {
	path := ""

	if os.Getenv("DB_PATH") != "" {
		path = os.Getenv("DB_PATH")
	} else {
		path = app.App.Storage().RootURI().Path() + "/sql.db"
		app.InfoLog.Println("db in:", path)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// setupDB inicialitza el repositori i executa les migracions
func (app *Config) setupDB(sqlDB *sql.DB) {
	app.DB = repository.NewSQLiteRepository(sqlDB)

	err := app.DB.Migrate()
	if err != nil {
		app.ErrorLog.Println(err)
		log.Panic(err)
	}
}
