package handlers

import (
	"context"
	"errors"
	"net/http"
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
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetTokenFromRequest extracts the JWT token from the Authorization header
func GetTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	// Check if the header starts with "Bearer "
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", errors.New("invalid authorization header format")
	}

	return authHeader[7:], nil
}

// AuthMiddleware is a middleware that validates JWT tokens
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := GetTokenFromRequest(r)
		if err != nil {
			SendError(w, http.StatusUnauthorized, err.Error())
			return
		}

		claims, err := ValidateToken(tokenString)
		if err != nil {
			SendError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Add the user ID to the request context
		ctx := context.WithValue(r.Context(), "userId", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
} 