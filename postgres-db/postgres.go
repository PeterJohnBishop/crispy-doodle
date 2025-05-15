package postgresdb

import (
	"database/sql"
	"fmt"
	"log"

	"crispy-doodle/main.go/global"
)

func ConnectPSQL(db *sql.DB) *sql.DB {

	host := global.PostgresHost
	port := global.PostgresPort
	user := global.PostgresUser
	password := global.PostgresPassword
	dbname := global.PostgresDBName
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
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

	log.Printf("[CONNECTED] to Postgres on :%s", port)
	return mydb
}
