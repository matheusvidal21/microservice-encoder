package database

import (
	"errors"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
	"github.com/matheusvidal21/microservice-encoder/domain"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type Database struct {
	Db            *gorm.DB
	Dsn           string
	DsnTest       string
	DbType        string
	DbTypeTest    string
	Debug         bool
	AutoMigrateDb bool
	Env           string
}

func NewDatabase() *Database {
	return &Database{}
}

func NewDatabaseTest() *gorm.DB {
	dbInstance := NewDatabase()
	dbInstance.Env = "test"
	dbInstance.DbTypeTest = "sqlite3"
	dbInstance.DsnTest = ":memory:"
	dbInstance.AutoMigrateDb = true
	dbInstance.Debug = true

	connection, err := dbInstance.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}
	return connection
}

func (db *Database) Connect() (*gorm.DB, error) {
	var err error
	var dsn string
	var dbType string

	if db.Env == "test" {
		dsn = db.DsnTest
		dbType = db.DbTypeTest
	} else {
		dsn = db.Dsn
		dbType = db.DbType
	}

	switch dbType {
	case "sqlite3":
		db.Db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	case "postgres":
		db.Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		return nil, errors.New("Unsupported database type")
	}

	if err != nil {
		return nil, err
	}

	if db.Debug {
		db.Db.Logger.LogMode(4)
	}

	if db.AutoMigrateDb {
		db.Db.AutoMigrate(&domain.Job{}, &domain.Video{})
	}
	return db.Db, nil
}
