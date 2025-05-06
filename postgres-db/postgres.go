package postgresdb

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func ConnectPSQL(db *sql.DB) *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//host := "localhost"
	host := "postgres"
	port := 5432
	user := os.Getenv("PSQL_USER")
	password := os.Getenv("PSQL_PASSWORD")
	dbname := os.Getenv("PSQL_DBNAME")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// fmt.Println("Connecting with:", psqlInfo)

	mydb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = mydb.Ping()
	if err != nil {
		panic(err)
	}

	log.Printf("[CONNECTED] to Postgres on :%d", port)
	return mydb
}
