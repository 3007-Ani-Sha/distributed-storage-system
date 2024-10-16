package handlers

import (
	"fmt"
	"io"
	// "log"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"github.com/3007-Ani-Sha/distributed-storage-system/db"
	"os"
)

// UploadHandler handles file uploads
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	//CHeck for the correct JWT:
	cookie, err := r.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
        return
	}

	tokenStr := cookie.Value
	claims:=&Claims{}
	token, err:= jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

	if err != nil || !token.Valid {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

	// Limit upload size
	r.ParseMultipartForm(10 << 20) // 10 MB limit

	// Get the uploaded file
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error Retrieving File", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Create a user directory:
	userDir := "./uploads/" + claims.Email
	err = os.MkdirAll(userDir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Check if for the same user file is not uploaded.
	filePath := userDir + "/" + handler.Filename
	if _, err := os.Stat(filePath); err == nil{
		// If the file exists, notify the user to change the file name
        w.WriteHeader(http.StatusConflict)
        fmt.Fprintln(w, "File already exists. Please rename the file or delete the existing one.")
        return
	}

	// Create a file on the server
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error Creating File", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the server
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error Saving File", http.StatusInternalServerError)
		return
	}

	//Saving the information in the Database incase of the retieval:
	_, err = db.DB.Exec("INSERT INTO user_files (filename, path, email) VALUES ($1, $2, $3)", handler.Filename, filePath, claims.Email)
	if err != nil {
		http.Error(w, "Failed to store Data for the uploaded file.", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s", handler.Filename)
}

// DownloadHandler handles file downloads
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	// Get the file name from query parameters
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	// Open the file
	file, err := os.Open("./uploads/" + filename)
	if err != nil {
		http.Error(w, "File Not Found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set the headers and write the file to the response
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "application/octet-stream")

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error Downloading File", http.StatusInternalServerError)
		return
	}
}
