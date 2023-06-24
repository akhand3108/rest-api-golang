package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/akhand3108/restgo/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type ContextString string

// const userIdString ContextString = "userId"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid token")
)

type AuthService struct {
	DB             *sql.DB
	TokenSecretKey string
}

func (as *AuthService) Signup(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)

	_, err = as.DB.Exec("INSERT INTO users (username, passwordhash) VALUES ($1, $2)", user.Username, user.PasswordHash)

	return err
}

func (as *AuthService) Signin(credentials *models.Credentials) (string, error) {
	var user models.User
	row := as.DB.QueryRow("SELECT * FROM users WHERE username = $1", credentials.Username)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrUserNotFound
		}
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password))
	if err != nil {
		return "", ErrInvalidCredentials
	}
	token, err := generateToken(user.ID, as.TokenSecretKey)

	if err != nil {
		return "", err
	}

	return token, nil
}

func generateToken(userID int, tokenSecretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
		"expiresAt": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(tokenSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func parseToken(tokenString string, tokenSecretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(tokenSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (as *AuthService) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println()
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Println(98, tokenString)
		if tokenString == authHeader {
			fmt.Println(92, tokenString == authHeader)
			http.Error(w, "Invalid token", http.StatusUnauthorized)

			return
		}

		claims, err := parseToken(tokenString, as.TokenSecretKey)
		if err != nil {
			fmt.Println(119, err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userID := int(claims[string("userID")].(float64))
		if userID == 0 {
			fmt.Println(119, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

	})
}
