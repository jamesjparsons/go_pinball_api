package handlers

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var (
	// JWTSecret is loaded from environment variables
	JWTSecret string
	// TokenExpiration is loaded from environment variables
	TokenExpiration time.Duration
)

func init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	// Load JWT secret
	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		panic("JWT_SECRET environment variable is required")
	}

	// Load token expiration
	expHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	if err != nil {
		panic("JWT_EXPIRATION_HOURS must be a valid integer")
	}
	TokenExpiration = time.Duration(expHours) * time.Hour
}

// Claims represents the JWT claims
type Claims struct {
	UserID uint `json:"userId"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(userID uint) (string, error) {
	// Get the JWT secret from environment variables
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Printf("JWT_SECRET not set in environment variables")
		return "", fmt.Errorf("JWT_SECRET not set")
	}

	log.Printf("Generating token for user ID: %d", userID)

	// Create the claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":     time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return "", err
	}

	log.Printf("Token generated successfully for user ID: %d", userID)
	return tokenString, nil
}

// ValidateToken verifies a JWT token and returns the user ID
func ValidateToken(tokenString string) (uint, error) {
	// Get the JWT secret from environment variables
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Printf("JWT_SECRET not set in environment variables")
		return 0, fmt.Errorf("JWT_SECRET not set")
	}

	log.Printf("Validating token: %s", tokenString)

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return 0, err
	}

	if !token.Valid {
		log.Printf("Token is invalid")
		return 0, fmt.Errorf("invalid token")
	}

	// Get the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("Invalid token claims")
		return 0, fmt.Errorf("invalid token claims")
	}

	// Get the user ID
	userID, ok := claims["user_id"].(float64)
	if !ok {
		log.Printf("Invalid user_id claim")
		return 0, fmt.Errorf("invalid user_id claim")
	}

	log.Printf("Token validated successfully for user ID: %d", uint(userID))
	return uint(userID), nil
}

