package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"erastos-backend/internal/config"

	"erastos-backend/internal/models"
	"erastos-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request CreateUserRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if request.Username == "" || request.Email == "" || request.Password == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "username, email, and password are required"})
			return
		}


		if len(request.Password) < 8 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "password must be at least 8 characters long"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error hashing password " + err.Error()})
			return
		}

		user := &models.User{
			Username: request.Username,
			Email:    request.Email,
			Password: string(hashedPassword),
		}

		createdUser, err := repository.CreateUser(pool, user)

		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				c.IndentedJSON(http.StatusConflict, gin.H{"message": "user with this email or username already exists"})
				return
			}
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error creating user " + err.Error()})
			return
		}

		c.IndentedJSON(http.StatusCreated, createdUser)
	}
}

func LoginHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request LoginRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		user, err := repository.GetUserByEmail(pool, request.Email)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error retrieving user " + err.Error()})
			return
		}

		if user == nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "invalid email or password"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))

		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "invalid email or password"})
			return
		}

		claims := jwt.MapClaims{
			"user_id": user.ID,
			"email":   user.Email,
			"exp":     time.Now().Add(time.Hour * 72).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error generating token " + err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, LoginResponse{Token: tokenString})
	}
}