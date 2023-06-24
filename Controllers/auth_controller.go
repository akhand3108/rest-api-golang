package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akhand3108/restgo/models"
	"github.com/akhand3108/restgo/services"
)

type AuthController struct {
	AuthService *services.AuthService
}

func NewAuthController(db *sql.DB, secretkey string) *AuthController {

	authService := &services.AuthService{
		DB:             db,
		TokenSecretKey: secretkey,
	}
	return &AuthController{
		AuthService: authService,
	}
}

func (ac *AuthController) Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = ac.AuthService.Signup(&user)
	if err != nil {
		http.Error(w, "Failed to signup user", http.StatusInternalServerError)
		fmt.Printf("Error: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"userCreated": "true"})

}

func (ac *AuthController) Signin(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := ac.AuthService.Signin(&credentials)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		fmt.Printf("Error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (ac *AuthController) Signout(w http.ResponseWriter, r *http.Request) {
	// TODO: Perform signout logic here
	w.WriteHeader(http.StatusOK)
}
