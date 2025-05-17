package postgresdb

import (
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	Password string         `json:"password"`
	Online   bool           `json:"online"`
	Channels pq.StringArray `json:"channels" sql:"type:text[]"`
	Created  int64          `json:"created"`
	Updated  int64          `json:"updated"`
}

func CreateUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT UNIQUE NOT NULL PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		online BOOL DEFAULT false,
		channels TEXT[],
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
		fmt.Println("Failed to bind JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("Registering user with email:", user.Email)

	userId := GenerateUserID(user.Email)
	fmt.Println("Generated user ID:", userId)

	hashedPassword, err := HashedPassword(user.Password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	fmt.Println("Password hashed successfully")

	query := `INSERT INTO users (id, name, email, password, online, channels)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = db.ExecContext(c, query, userId, user.Name, user.Email, hashedPassword, user.Online, pq.Array(user.Channels))
	if err != nil {
		fmt.Println("Database insert error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("User registered successfully:", userId)
	c.JSON(http.StatusCreated, gin.H{"message": "User created!"})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(db *sql.DB, c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Failed to bind JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	fmt.Println("Login attempt for email:", req.Email)

	var user User
	var channels pq.StringArray
	query := `SELECT id, name, email, password, online, channels, created, updated FROM users WHERE email = $1`
	err := db.QueryRowContext(c, query, req.Email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password,
		&user.Online, &channels, &user.Created, &user.Updated,
	)
	user.Channels = channels
	if user.Channels == nil {
		user.Channels = []string{}
	}
	if err == sql.ErrNoRows {
		fmt.Println("No user found with email:", req.Email)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		fmt.Println("Database query error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("User found:", user.ID)

	if !CheckPasswordHash(req.Password, user.Password) {
		fmt.Println("Password verification failed for user:", user.ID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password Verification Failed"})
		return
	}
	fmt.Println("Password verified for user:", user.ID)

	access, refresh, err := GenerateTokens(user.ID)

	fmt.Println("Login successful for user:", user.ID)
	c.JSON(http.StatusOK, gin.H{
		"message":      "Login Success",
		"token":        access,
		"refreshToken": refresh,
		"user":         user,
	})
}
func Refresh(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing refresh token"})
		return
	}

	claims, err := ValidateToken(body.RefreshToken, true)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	newAccess, _, err := GenerateTokens(claims.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newAccess})
}

func GetUsers(db *sql.DB, c *gin.Context) {
	fmt.Println("Fetching all users")
	rows, err := db.QueryContext(c, "SELECT id, name, email, password, online, channels, created, updated FROM users;")
	if err != nil {
		fmt.Println("Query failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Online, &user.Channels, &user.Created, &user.Updated); err != nil {
			fmt.Println("Row scan failed:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if user.Channels == nil {
			user.Channels = []string{}
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Row iteration error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("Total users fetched:", len(users))
	c.JSON(http.StatusOK, users)
}

func GetUserByID(db *sql.DB, c *gin.Context) {
	id := c.Param("id")
	fmt.Println("Fetching user with ID:", id)

	var user User
	query := `SELECT id, name, email, password, online, channels, created, updated FROM users WHERE id = $1`
	err := db.QueryRowContext(c, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Online, &user.Channels, &user.Created, &user.Updated)
	if err == sql.ErrNoRows {
		fmt.Println("User not found:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		fmt.Println("Query error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("User fetched:", user.ID)
	c.JSON(http.StatusOK, user)
}

func UpdateUserByID(db *sql.DB, c *gin.Context) {
	id := c.Param("id")
	fmt.Println("Updating user with ID:", id)

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println("Failed to bind request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE users SET name=$1, email=$2, password=$3, online=$4, channels=$5, updated=EXTRACT(EPOCH FROM now()) WHERE id=$6`
	result, err := db.ExecContext(c, query, user.Name, user.Email, user.Password, user.Online, user.Channels, id)
	if err != nil {
		fmt.Println("Update query failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Failed to retrieve rows affected:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		fmt.Println("No rows updated for ID:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	fmt.Println("User updated successfully:", id)
	c.JSON(http.StatusOK, gin.H{"message": "User updated!"})
}

func DeleteUserByID(db *sql.DB, c *gin.Context) {
	id := c.Param("id")
	fmt.Println("Deleting user with ID:", id)

	query := `DELETE FROM users WHERE id = $1`
	result, err := db.ExecContext(c, query, id)
	if err != nil {
		fmt.Println("Delete query failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Failed to retrieve rows affected:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		fmt.Println("No user found to delete with ID:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	fmt.Println("User deleted successfully:", id)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted!"})
}
