package postgresdb

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Channel struct {
	ID       string   `json:"id"`
	Title    string   `json:"text"`
	Messages []string `json:"messages"`
	Created  int64    `json:"created"`
	Updated  int64    `json:"updated"`
}

func CreateChannelsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS channels (
		id TEXT UNIQUE NOT NULL PRIMARY KEY,
		Title TEXT,
		Messages TEXT[], 
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);`

	_, err := db.Exec(query)
	return err
}

func GenerateChannelID() string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("channel_%d", timestamp)
}

func CreateChannel(db *sql.DB, c *gin.Context) {
	var channel Channel

	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	channelID := GenerateChannelID()

	query := `INSERT INTO channels (id, title, messages)
		VALUES ($1, $2, $3)`
	_, err := db.ExecContext(c, query, channelID, channel.Title, pq.Array(channel.Messages))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Message sent!"})
}

func GetChannels(db *sql.DB, c *gin.Context) {
	rows, err := db.QueryContext(c, "SELECT id, title, messages, created, updated FROM channels;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var channels []Channel
	for rows.Next() {
		var channel Channel
		if err := rows.Scan(&channel.ID, &channel.Title, &channel.Messages, &channel.Created, &channel.Updated); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		channels = append(channels, channel)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, channels)
}

func GetChannelByID(db *sql.DB, c *gin.Context) {
	id := c.Param("id")
	var channel Channel
	query := `SELECT id, title, messages, created, updated FROM channels WHERE id = $1`

	err := db.QueryRowContext(c, query, id).Scan(&channel.ID, &channel.Title, &channel.Messages, &channel.Created, &channel.Updated)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, channel)
}

func UpdateChannelByID(db *sql.DB, c *gin.Context) {
	id := c.Param("id")
	var channel Channel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE channels SET title=$1, messages=$2, updated=EXTRACT(EPOCH FROM now()) WHERE id=$3`
	result, err := db.ExecContext(c, query, channel.Title, channel.Messages, id)
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

	c.JSON(http.StatusOK, gin.H{"message": "Channel updated!"})
}

func DeleteChannelByID(db *sql.DB, c *gin.Context) {
	id := c.Param("id")
	query := `DELETE FROM channels WHERE id = $1`
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Channel deleted!"})
}
