package utils

import (
	"context"
	"net/http"
	"strings"
    "log"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("Lf8+5T+RHNskcLbs6D5q/kEr8Y6Q2F6ovhmf1A2jmdk=")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Malformed token", http.StatusUnauthorized)
			return
		}

		claims := &struct {
			UserID string `json:"user_id"`
			jwt.RegisteredClaims
		}{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WebSocketAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("WebSocketAuthMiddleware invoked")

		tokenString := r.URL.Query().Get("token")
		if tokenString == "" {
			log.Println("Missing token in query parameters")
			http.Error(w, "Missing token in query parameters", http.StatusUnauthorized)
			return
		}
		log.Println("Token received:", tokenString)

		claims := &struct {
			UserID string `json:"user_id"`
			jwt.RegisteredClaims
		}{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			log.Printf("Error parsing token: %v", err)
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			log.Println("Token is invalid or expired")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		log.Printf("Token valid for userID: %s", claims.UserID)
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}