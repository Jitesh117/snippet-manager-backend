package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var JWTKey = []byte("your_secret_key")

func ExtractUserIDFromToken(r *http.Request) (uuid.UUID, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return uuid.UUID{}, fmt.Errorf("Authorization header missing")
	}

	tokenString := strings.Split(authHeader, " ")[1] // "Bearer <token>"

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return JWTKey, nil
		},
	)
	if err != nil || !token.Valid {
		return uuid.UUID{}, fmt.Errorf("Invalid token")
	}

	userIDStr := (*claims)["user_id"].(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.UUID{}, err
	}
	return userID, nil
}

type contextKey string

const UserContextKey = contextKey("user_id")

func GenerateJWT(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := ExtractUserIDFromToken(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// // Helper function to extract user ID from the request context
// func GetUserIDFromContext(r *http.Request) (uuid.UUID, error) {
// 	userID, ok := r.Context().Value(UserContextKey).(uuid.UUID)
// 	if !ok {
// 		return uuid.UUID{}, fmt.Errorf("could not retrieve user ID from context")
// 	}
// 	return userID, nil
// }
