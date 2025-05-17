package ginserver

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"crispy-doodle/main.go/awservice"
	openai "crispy-doodle/main.go/open-ai"
	postgresdb "crispy-doodle/main.go/postgres-db"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func StartGinServer() {

	// connect to AWS
	cfg := awservice.StartAws()

	// connect to AWS S3
	s3Client := awservice.ConnectS3(cfg)

	// connecting to OpenAI
	ai := openai.OpenAI()
	if ai == nil {
		log.Fatal("Error connecting to OpenAI")
	}

	// connecting to Postgres
	db := postgresdb.ConnectPSQL(db)
	err := db.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	postgresdb.CreateUsersTable(db)

	// creating gin server
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(postgresdb.JWTMiddleware())
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "crispy-doodle",
		})
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"aws_s3":   "AWS S3 connected",
			"ai":       "OpenAI connected",
			"database": "Postgres connected",
			"server":   "Gin server running",
		})
	})
	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 0,
	}

	addUserRoutes(router, db)
	addChannelRoutes(router, db)
	addMessageRoutes(router, db)
	addAWSRoutes(router, s3Client)

	log.Println("[CONNECTED] Gin server on :8080")
	s.ListenAndServe()
}
