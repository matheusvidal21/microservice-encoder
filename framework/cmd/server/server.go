package main

import (
	"github.com/joho/godotenv"
	"github.com/matheusvidal21/microservice-encoder/application/services"
	"github.com/matheusvidal21/microservice-encoder/framework/database"
	"github.com/matheusvidal21/microservice-encoder/framework/queue"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strconv"
)

var db database.Database

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	autoMigrateDb, err := strconv.ParseBool(os.Getenv("AUTO_MIGRATE_DB"))
	if err != nil {
		log.Fatalf("Error parsing AUTO_MIGRATE_DB")
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Fatalf("Error parsing DEBUG")
	}

	db.AutoMigrateDb = autoMigrateDb
	db.Debug = debug
	db.DsnTest = os.Getenv("DSN_TEST")
	db.Dsn = os.Getenv("DSN")
	db.DbTypeTest = os.Getenv("DB_TYPE_TEST")
	db.DbType = os.Getenv("DB_TYPE")
	db.Env = os.Getenv("ENV")
}

func main() {
	messageChannel := make(chan amqp.Delivery)
	jobReturnChannel := make(chan services.JobWorkerResult)

	dbConnection, err := db.Connect()
	if err != nil {
		log.Fatalf("Error connecting to the database")
	}
	dbSql, _ := dbConnection.DB()
	defer dbSql.Close()

	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()
	defer ch.Close()

	rabbitMQ.Consume(messageChannel)

	jobManager := services.NewJobManager(dbConnection, rabbitMQ, messageChannel, jobReturnChannel)
	jobManager.Start(ch)
}
