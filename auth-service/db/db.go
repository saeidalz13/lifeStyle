package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/saeidalz13/lifestyle/auth-service/internal/dberr"
)

const (
	maxOpenConns = 300
	maxIdleConns = 100
	connMaxLife  = time.Minute * 15
	// there is a 'SchemeFromURL' function that splits the migrationDir by ':', so db/migration will be the URL
	migrationDir = "file:db/migration"
)

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func mustMigrate(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{DatabaseName: "lifestyle_auth"})
	checkErr(err)

	m, err := migrate.NewWithDatabaseInstance(migrationDir, "lifestyle_auth", driver)
	checkErr(err)

	version, dirty, err := m.Version()
	if err != nil {
		if err.Error() == dberr.ErrNoMigration.Error() {
			log.Println(err)
		} else {
			log.Fatalln(err)
		}
	}
	if dirty {
		log.Fatalf("%v; version: %d", dberr.ErrDirtyDb, version)
	}

	if err := m.Up(); err != nil {
		if err.Error() == dberr.ErrMigrationNoChange.Error() || err.Error() == dberr.ErrFileNotExists.Error() {
			return
		}

		log.Fatalln(err)
	}

	log.Println("migration successful")
}

func MustConnectDb(dbUrl string) *sql.DB {
	db, err := sql.Open("postgres", dbUrl)
	checkErr(err)

	err = db.Ping()
	checkErr(err)

	mustMigrate(db)

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLife)

	return db
}
