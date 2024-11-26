package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Token string `json:"token"`
}

type Claims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}

var jwtKey = []byte("Lf8+5T+RHNskcLbs6D5q/kEr8Y6Q2F6ovhmf1A2jmdk=")

type LoginHandler struct {
    DB *gorm.DB
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    var user struct {
        ID           string
        PasswordHash string
    }
    err := h.DB.Table("users").Where("username = ?", req.Username).First(&user).Error
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    expirationTime := time.Now().Add(1 * time.Hour)
    claims := &Claims{
        UserID: user.ID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "Could not create token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})
}