package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"github.com/3007-Ani-Sha/distributed-storage-system/db"
)

// SignUpRequest represents the signup request body
type SignUpRequest struct {
	Email    string `json:"email"`
	OTP      string `json:"otp"`
	Password string `json:"password"`
}

// SignUp handles user registration
func SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.OTP == "" || req.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Verify if the user is not already registered
	var storedUser string
	err = db.DB.QueryRow("SELECT user FROM users_app WHERE email = $1", req.Email).Scan(&storedUser)
	if err == nil {
		http.Error(w, "User Already registered.", http.StatusBadRequest)
		return
	}

	// Validate OTP
	var storedOTP string
	err = db.DB.QueryRow("SELECT otp FROM otps_app WHERE email = $1 AND is_used = FALSE ORDER BY created_at DESC LIMIT 1", req.Email).Scan(&storedOTP)
	if err == sql.ErrNoRows || storedOTP != req.OTP {
		http.Error(w, "Invalid or expired OTP", http.StatusUnauthorized)
		return
	}
	

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert user into the database
	_, err = db.DB.Exec("INSERT INTO users_app (email, password) VALUES ($1, $2)", req.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Mark OTP as used
	_, err = db.DB.Exec("UPDATE otps_app SET is_used = TRUE WHERE email = $1", req.Email)
	if err != nil {
		http.Error(w, "Failed to update OTP status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User registered successfully")
}
