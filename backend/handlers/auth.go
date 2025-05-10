package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"backend/models"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// Signup handles user registration
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object true "User registration details" {"email": "string", "password": "string", "firstName": "string", "lastName": "string"}
// @Success 201 {object} AuthResponse "User created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 409 {object} ErrorResponse "User already exists"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	var req struct {
		Email     string `json:"email" binding:"required"`
		Password  string `json:"password" binding:"required"`
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Signup error - Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if user already exists
	var existingUser models.User
	result := h.db.Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		log.Printf("Signup error - User already exists: %s", req.Email)
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Signup error - Password hashing failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	user := models.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  string(hashedPassword),
	}

	// Create user
	if err := h.db.Create(&user).Error; err != nil {
		log.Printf("Signup error - Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := GenerateToken(user.ID)
	if err != nil {
		log.Printf("Signup error - Token generation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	log.Printf("Signup success - User created: %s", user.Email)
	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
		},
		"token": token,
	})
}

// Login handles user authentication
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object true "Login credentials" {"email": "string", "password": "string"}
// @Success 200 {object} AuthResponse "Login successful"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Login error - Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	log.Printf("Login attempt for email: %s", req.Email)

	// Find user
	var user models.User
	result := h.db.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("Login error - User not found: %s", req.Email)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			log.Printf("Login error - Database error: %v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		}
		return
	}

	log.Printf("User found, checking password for: %s", req.Email)

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("Login error - Invalid password for user: %s", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	log.Printf("Password verified for user: %s", req.Email)

	// Generate JWT token
	token, err := GenerateToken(user.ID)
	if err != nil {
		log.Printf("Login error - Token generation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	log.Printf("Login success - User logged in: %s", req.Email)
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
		},
		"token": token,
	})
}

// GetCurrentUser retrieves the current authenticated user's information
// @Summary Get current user
// @Description Get information about the currently authenticated user
// @Tags auth
// @Produce json
// @Security Bearer
// @Success 200 {object} UserResponse "User information"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "User not found"
// @Router /auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		log.Printf("GetCurrentUser error - No user ID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		log.Printf("GetCurrentUser error - User not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	log.Printf("GetCurrentUser success - User retrieved: %s", user.Email)
	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"email":     user.Email,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
	})
}
