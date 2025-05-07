package postgresdb

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Message struct {
	ID      string   `json:"id"`
	Sender  string   `json:"sender"`
	Text    string   `json:"text"`
	Images  []string `json:"images"`
	Created int64    `json:"created"`
	Updated int64    `json:"updated"`
}

func CreateMessagesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS messages (
		id TEXT UNIQUE NOT NULL PRIMARY KEY,
		sender TEXT NOT NULL,
		text TEXT,
		images TEXT[], 
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);`

	_, err := db.Exec(query)
	return err
}

func GenerateMessageID() string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("message_%d", timestamp)
}

func CreateMessage(db *sql.DB, c *gin.Context) {
	var message Message

	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messageId := GenerateMessageID()

	query := `INSERT INTO messages (id, sender, text, images)
		VALUES ($1, $2, $3, $4)`
	_, err := db.ExecContext(c, query, messageId, message.Sender, message.Text, pq.Array(message.Images))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Message sent!"})
}

func GetMessages(db *sql.DB, c *gin.Context) {
	rows, err := db.QueryContext(c, "SELECT id, sender, text, images, created, updated FROM messages;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.Sender, &message.Text, &message.Images, &message.Created, &message.Updated); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func GetMessageById(db *sql.DB, c *gin.Context) {
	id := c.Param("id")
	var message Message
	query := `SELECT id, sender, text, images, created, updated FROM messages WHERE id = $1`

	err := db.QueryRowContext(c, query, id).Scan(&message.ID, &message.Sender, &message.Text, &message.Images, &message.Created, &message.Updated)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

func UpdateMessageByID(db *sql.DB, c *gin.Context) {
	id := c.Param("id")
	var message Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE messages SET sender=$1, text=$2, images=$3, updated=EXTRACT(EPOCH FROM now()) WHERE id=$5`
	result, err := db.ExecContext(c, query, message.Sender, message.Text, message.Images, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message updated!"})
}

func DeleteMessageByID(db *sql.DB, c *gin.Context) {
	id := c.Param("id")
	query := `DELETE FROM messages WHERE id = $1`
	result, err := db.ExecContext(c, query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted!"})
}
