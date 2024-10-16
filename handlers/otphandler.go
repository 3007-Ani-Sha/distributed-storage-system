// package handlers

// import (
// 	"crypto/rand"
// 	"encoding/json"
// 	"fmt"
// 	"math/big"
// 	"net/http"
// 	"log"

// 	"github.com/3007-Ani-Sha/distributed-storage-system/db"
// 	"github.com/3007-Ani-Sha/distributed-storage-system/mail"
// )

// // Generate random OTP
// func generateOTP() (string, error) {
// 	const otpChars = "0123456789"
// 	otpLength := 6
// 	otp := ""

// 	for i := 0; i < otpLength; i++ {
// 		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(otpChars))))
// 		if err != nil {
// 			return "", err
// 		}
// 		otp += string(otpChars[n.Int64()])
// 	}

// 	return otp, nil
// }
// func SendOTP(w http.ResponseWriter, r *http.Request) {
// 	// Set CORS headers for all requests
// 	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3001")
// 	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

// 	// Handle preflight request
// 	if r.Method == http.MethodOptions {
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	}

// 	// Continue with normal POST request handling
// 	// Parse request body
// 	var data struct {
// 		Email string `json:"email"`
// 	}

// 	err := json.NewDecoder(r.Body).Decode(&data)
// 	if err != nil || data.Email == "" {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		return
// 	}

// 	// Generate OTP
// 	otp, err := generateOTP()
// 	if err != nil {
// 		log.Println("Error generating OTP:", err)
// 		http.Error(w, "Failed to generate OTP", http.StatusInternalServerError)
// 		return
// 	}

// 	// Store OTP in the database
// 	_, err = db.DB.Exec("INSERT INTO otps_app (email, otp) VALUES ($1, $2)", data.Email, otp)
// 	if err != nil {
// 		http.Error(w, "Failed to store OTP", http.StatusInternalServerError)
// 		return
// 	}

// 	// Send OTP via email
// 	err = mail.SendOTPEmail(data.Email, otp)
// 	if err != nil {
// 		http.Error(w, "Failed to send OTP email", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprint(w, "OTP sent successfully\n")
// }

package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"log"

	"github.com/3007-Ani-Sha/distributed-storage-system/db"
	"github.com/3007-Ani-Sha/distributed-storage-system/mail"
)

// Generate random OTP
func generateOTP() (string, error) {
	const otpChars = "0123456789"
	otpLength := 6
	otp := ""

	for i := 0; i < otpLength; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(otpChars))))
		if err != nil {
			return "", err
		}
		otp += string(otpChars[n.Int64()])
	}

	return otp, nil
}
func SendOTP(w http.ResponseWriter, r *http.Request) {
	// Handle the POST request (sending the OTP)
    if r.Method == http.MethodPost {
        var data struct {
            Email string `json:"email"`
        }

        err := json.NewDecoder(r.Body).Decode(&data)
        if err != nil || data.Email == "" {
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }

        // OTP generation and sending logic
        otp, err := generateOTP()
        if err != nil {
            log.Println("Error generating OTP:", err)
            http.Error(w, "Failed to generate OTP", http.StatusInternalServerError)
            return
        }

        _, err = db.DB.Exec("INSERT INTO otps_app (email, otp) VALUES ($1, $2)", data.Email, otp)
        if err != nil {
            http.Error(w, "Failed to store OTP", http.StatusInternalServerError)
            return
        }

        err = mail.SendOTPEmail(data.Email, otp)
        if err != nil {
            http.Error(w, "Failed to send OTP email", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, "OTP sent successfully\n")
    } else {
        // If the method is not POST or OPTIONS, return 405
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}

