package ginserver

import (
	postgresdb "crispy-doodle/main.go/postgres-db"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func addUserRoutes(r *gin.Engine, db *sql.DB) {
	r.POST("/login", func(c *gin.Context) {
		postgresdb.Login(db, c)
	})
	r.POST("/register", func(c *gin.Context) {
		postgresdb.RegisterUser(db, c)
	})
	r.GET("/users", func(c *gin.Context) {
		postgresdb.GetUsers(db, c)
	})
	r.GET("/users/:id", func(c *gin.Context) {
		postgresdb.GetUserByID(db, c)
	})
	r.PUT("/users/:id", func(c *gin.Context) {
		postgresdb.UpdateUserByID(db, c)
	})
	r.DELETE("/users/:id", func(c *gin.Context) {
		postgresdb.DeleteUserByID(db, c)
	})
}

func addMessageRoutes(r *gin.Engine, db *sql.DB) {
	r.POST("/messages", func(c *gin.Context) {
		postgresdb.CreateMessage(db, c)
	})
	r.GET("/messages", func(c *gin.Context) {
		postgresdb.GetMessages(db, c)
	})
	r.GET("/messages/:id", func(c *gin.Context) {
		postgresdb.GetMessageById(db, c)
	})
	r.PUT("/messages/:id", func(c *gin.Context) {
		postgresdb.UpdateMessageByID(db, c)
	})
	r.DELETE("/messages/:id", func(c *gin.Context) {
		postgresdb.DeleteMessageByID(db, c)
	})
}
