package global

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var PostgresUser string
var PostgresPassword string
var PostgresDBName string
var PostgresHost string
var PostgresPort string

var OpenAIKey string

var AwsAccessKey string
var AwsSecretKey string
var AwsRegion string
var AwsBucket string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	getPostgresEnvs()
	getAWSEnvs()
	getOpenAIEnvs()

}

func getPostgresEnvs() {

	PostgresPassword = os.Getenv("PSQL_PASSWORD")
	if PostgresPassword == "" {
		log.Fatal("PSQL_PASSWORD is not set in .env file")
	}
	PostgresUser = os.Getenv("PSQL_USER")
	if PostgresUser == "" {
		log.Fatal("PSQL_USER is not set in .env file")
	}
	PostgresDBName = os.Getenv("PSQL_DBNAME")
	if PostgresDBName == "" {
		log.Fatal("PSQL_DBNAME is not set in .env file")
	}
	PostgresHost = os.Getenv("PSQL_HOST")
	if PostgresHost == "" {
		log.Fatal("PSQL_HOST is not set in .env file")
	}
	PostgresPort = os.Getenv("PSQL_PORT")
	if PostgresPort == "" {
		log.Fatal("PSQL_PORT is not set in .env file")
	}

	log.Println("Postgres Environment Variables Loaded")

}

func getAWSEnvs() {

	AwsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	if AwsAccessKey == "" {
		log.Fatal("AWS_ACCESS_KEY_ID is not set in .env file")
	}
	AwsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	if AwsSecretKey == "" {
		log.Fatal("AWS_SECRET_ACCESS_KEY is not set in .env file")
	}
	AwsRegion := os.Getenv("AWS_REGION")
	if AwsRegion == "" {
		log.Fatal("AWS_REGION is not set in .env file")
	}
	AwsBucket := os.Getenv("AWS_BUCKET")
	if AwsBucket == "" {
		log.Fatal("AWS_BUCKET is not set in .env file")
	}

	log.Println("AWS Environment Variables Loaded")

}

func getOpenAIEnvs() {

	OpenAIKey := os.Getenv("OPENAI_API_KEY")
	if OpenAIKey == "" {
		log.Fatal("OPENAI_API_KEY is not set in .env file")
	}

	log.Println("Postgres Environment Variables Loaded")

}
