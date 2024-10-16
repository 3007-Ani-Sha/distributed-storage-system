package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	// "fmt"
	"github.com/3007-Ani-Sha/distributed-storage-system/db"
	"github.com/dgrijalva/jwt-go"
    "time"
)

type LoginRequest struct {
	Email string `json:"email"`
	Password   string `json:"password"`
}

var jwtKey = []byte("e2bb9e8e6d7b60f4864b2b3e9d1e120e8b2c9f1b") // Random key that only the backend knows to authenticate the sign for any user, i.e. my signature on the user data.

type Claims struct {
    Email string `json:"email"`
    jwt.StandardClaims
}

func generateJWT(email string) (string, time.Time, error) {
    expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
    claims := &Claims{
        Email: email,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return "", expirationTime, err
    }
    return tokenString, expirationTime, nil
}

// Login handles user login
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate Password:
	var storedpass string
	err = db.DB.QueryRow("SELECT password FROM users_app WHERE email = $1", req.Email).Scan(&storedpass)
	if err == sql.ErrNoRows {
		http.Error(w, "Not Registered User", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedpass), []byte(req.Password))

	if err!=nil {
		http.Error(w, "Invalid Credentials!", http.StatusUnauthorized)
		return
	}

	tokenString, expirationTime, err := generateJWT(req.Email)

	if err!=nil {
		http.Error(w, "Error in creating the JWT", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
        Name:    "token",
        Value:   tokenString,
        Expires: expirationTime,
    })

    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
		"message": "Login Successfull!",
	})

	// fmt.Fprint(w, "Login successful")
}
