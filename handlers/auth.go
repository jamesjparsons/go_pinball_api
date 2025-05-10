package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"backend/models"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// @Summary     Login user
// @Description Authenticate a user and return a JWT token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       credentials body LoginRequest true "Login credentials"
// @Success     200 {object} LoginResponse
// @Failure     400 {object} ErrorResponse
// @Failure     401 {object} ErrorResponse
// @Router      /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Model.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: tokenString,
		User:  user,
	})
}

// @Summary     Register new user
// @Description Create a new user account
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       user body SignupRequest true "User registration data"
// @Success     200 {object} SignupResponse
// @Failure     400 {object} ErrorResponse
// @Router      /auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := user.SetPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	c.JSON(http.StatusOK, SignupResponse{
		User: user,
	})
}

// @Summary     Get current user
// @Description Get the currently authenticated user's information
// @Tags        auth
// @Produce     json
// @Security    Bearer
// @Success     200 {object} models.User
// @Failure     401 {object} ErrorResponse
// @Router      /auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type LoginResponse struct {
	Token string             `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  models.SwaggerUser `json:"user"`
}

type SignupRequest struct {
	Email     string `json:"email" binding:"required,email" example:"user@example.com"`
	Password  string `json:"password" binding:"required,min=6" example:"password123"`
	FirstName string `json:"firstName" binding:"required" example:"John"`
	LastName  string `json:"lastName" binding:"required" example:"Doe"`
}

type SignupResponse struct {
	User models.SwaggerUser `json:"user"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
