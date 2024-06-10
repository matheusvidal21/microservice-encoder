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

var db *database.Database

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	autoMigrateDb, err := strconv.ParseBool(os.Getenv("AUTO_MIGRATE_DB"))
	if err != nil {
		log.Fatalf("Error parsing boolean env var")
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Fatalf("Error parsing boolean env var")
	}

	db = &database.Database{
		AutoMigrateDb: autoMigrateDb,
		Debug:         debug,
		DsnTest:       os.Getenv("DSN_TEST"),
		Dsn:           os.Getenv("DSN"),
		DbTypeTest:    os.Getenv("DB_TYPE_TEST"),
		DbType:        os.Getenv("DB_TYPE"),
		Env:           os.Getenv("ENV"),
	}
	log.Printf("Database configuration: %+v\n", db)

}

func main() {
	messageChannel := make(chan amqp.Delivery)
	jobReturnChannel := make(chan services.JobWorkerResult)

	dbConnection, err := db.Connect()
	if err != nil {
		log.Fatalf("Error connecting to the database, error: %v", err)
	}
	defer dbConnection.Close()

	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()
	defer ch.Close()

	rabbitMQ.Consume(messageChannel)

	jobManager := services.NewJobManager(dbConnection, rabbitMQ, messageChannel, jobReturnChannel)
	jobManager.Start(ch)
}
