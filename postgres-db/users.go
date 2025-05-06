package postgresdb

import (
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Online   bool   `json:"online"`
	Created  int64  `json:"created"`
	Updated  int64  `json:"updated"`
}

func CreateUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT UNIQUE NOT NULL PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		online BOOL DEFAULT false,
		created BIGINT DEFAULT (EXTRACT(EPOCH FROM now())),
    	updated BIGINT DEFAULT (EXTRACT(EPOCH FROM now()))
	);`

	_, err := db.Exec(query)
	return err
}

func GenerateUserID(email string) string {
	hash := sha256.Sum256([]byte(email))
	return fmt.Sprintf("user_%d", binary.BigEndian.Uint64(hash[:8]))
}

func HashedPassword(password string) (string, error) {
	hashedPassword, error := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashedPassword), error
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RegisterUser(db *sql.DB, c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := GenerateUserID(user.Email)

	hashedPassword, err := HashedPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	query := `INSERT INTO users (id, name, email, password, online)
		VALUES ($1, $2, $3, $4, $5)`
	_, err = db.ExecContext(c, query, userId, user.Name, user.Email, hashedPassword, user.Online)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created!"})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(db *sql.DB, c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var user User
	query := `SELECT id, name, email, password, online, created, updated FROM users WHERE email = $1`
	err := db.QueryRowContext(c, query, req.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Online, &user.Created, &user.Updated)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password Verification Failed"})
		return
	}

	userClaims := UserClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		},
	}

	token, err := NewAccessToken(userClaims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})
		return
	}

	refreshToken, err := NewRefreshToken(userClaims.StandardClaims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Login Success",
		"token":        token,
		"refreshToken": refreshToken,
		"user":         user,
	})
}

func RefreshTokenHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		const userIDKey ContextKey = "userID"

		id, ok := c.Request.Context().Value(userIDKey).(string)
		if !ok || id == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ID not found in context"})
			return
		}

		var user User
		query := `SELECT id, name, email, password, online, created, updated FROM users WHERE id = $1`
		err := db.QueryRowContext(c, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Online, &user.Created, &user.Updated)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		userClaims := UserClaims{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			StandardClaims: jwt.StandardClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			},
		}

		token, err := NewAccessToken(userClaims)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})
			return
		}

		refreshToken, err := NewRefreshToken(userClaims.StandardClaims)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      "Token Refreshed",
			"token":        token,
			"refreshToken": refreshToken,
		})
	}
}

func GetUsers(db *sql.DB, c *gin.Context) {
	rows, err := db.QueryContext(c, "SELECT id, name, email, password, online, created, updated FROM users;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Online, &user.Created, &user.Updated); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetUserByID(db *sql.DB, c *gin.Context) {
	id := c.Param("id")
	var user User
	query := `SELECT id, name, email, password, online, created, updated FROM users WHERE id = $1`
	err := db.QueryRowContext(c, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Online, &user.Created, &user.Updated)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUserByID(db *sql.DB, c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE users SET name=$1, email=$2, password=$3, online=$4, updated=EXTRACT(EPOCH FROM now()) WHERE id=$5`
	result, err := db.ExecContext(c, query, user.Name, user.Email, user.Password, user.Online, id)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated!"})
}

func DeleteUserByID(db *sql.DB, c *gin.Context) {
	id := c.Param("id")
	query := `DELETE FROM users WHERE id = $1`
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
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted!"})
}
