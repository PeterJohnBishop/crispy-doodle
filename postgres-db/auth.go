package postgresdb

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var AccessTokenSecret = []byte(os.Getenv("TOKEN_SECRET"))
var RefreshTokenSecret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))
var AccessTokenTTL = time.Minute * 15
var RefreshTokenTTL = time.Hour * 24 * 7

type UserClaims struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func NewAccessToken(claims UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func NewRefreshToken(claims jwt.StandardClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func ParseAccessToken(accessToken string) *UserClaims {
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil || !parsedAccessToken.Valid {
		return nil
	}

	return parsedAccessToken.Claims.(*UserClaims)
}

func ParseRefreshToken(refreshToken string) *jwt.StandardClaims {
	parsedRefreshToken, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil || !parsedRefreshToken.Valid {
		return nil
	}

	return parsedRefreshToken.Claims.(*jwt.StandardClaims)
}

type ContextKey string

const UserIDKey ContextKey = "userID"

type VerifyRefreshRequest struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func VerifyJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/register") || strings.HasPrefix(c.Request.URL.Path, "/login") || strings.HasPrefix(c.Request.URL.Path, "/health") {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication Header is missing!"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		userClaims := ParseAccessToken(token)
		if userClaims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to verify token!"})
			c.Abort()
			return
		}

		c.Set("userClaims", userClaims)

		c.Next()
	}
}

func VerifyRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VerifyRefreshRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			c.Abort()
			return
		}

		claims := ParseRefreshToken(req.Token)
		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to verify token!"})
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), UserIDKey, req.ID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
