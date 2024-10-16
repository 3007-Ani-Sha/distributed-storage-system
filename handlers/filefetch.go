package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "github.com/gorilla/mux"
	"github.com/dgrijalva/jwt-go"
	"github.com/3007-Ani-Sha/distributed-storage-system/db"
)

// type File struct {
//     id        uint   `json:"id" gorm:"primaryKey"`
//     filename      string `json:"filename"`
//     path      string `json:"path"`
//     email string `json:"user_email"`
// }

type File struct {
    File_id       int    `json:"fid"`
    Filename string `json:"filename"`
    Path     string `json:"path"`
    Email    string `json:"email"`
}


// Fetch all files uploaded by the user
func FetchFilesHandler(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("token")
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    tokenStr := cookie.Value
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil || !token.Valid {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

	rows, err := db.DB.Query("SELECT id, filename, path, email FROM user_files WHERE email = $1", claims.Email)
	if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error fetching files: ", err)
		return
    }
	defer rows.Close()

	// Putting the file info into the arrayto pass in the response:
    var files []File

	for rows.Next() {
        var file File
        if err := rows.Scan(&file.File_id, &file.Filename, &file.Path,
			&file.Email); err != nil {
            w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error Scanning files")
			return
        }
        // fmt.Println("File:", file)
        files = append(files, file)
    }

	if err:= rows.Err(); err!=nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
    // Explicitly check if files array is empty
    if len(files) == 0 {
        // Return an explicit empty array
        w.Write([]byte("[]"))  // This returns an empty JSON array
        return
    }
	json.NewEncoder(w).Encode(files)
}

// Download a file by its ID:
func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
    
	//Verifying the token:
	cookie, err := r.Cookie("token")
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    tokenStr := cookie.Value
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil || !token.Valid {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
	
	// Extracting the requested information: 
    // Ensure the file belongs to the logged-in user
	vars := mux.Vars(r)
    filename := vars["filename"]
	email_asked := vars["email"]

	if email_asked != claims.Email {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "User for which data is requested and you do not match!")
		return
	}

	// Path to the stored file in the database:
    var file File
	err = db.DB.QueryRow("SELECT id, filename, path, email FROM user_files WHERE email = $1 and filename = $2", 
			claims.Email, filename).Scan(&file.File_id, &file.Filename, &file.Path,
			&file.Email)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintln(w, "File not found")
        return
    }

    // Serve the file for download
    http.ServeFile(w, r, file.Path)
}

// Delete a file by its ID:
func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {

	// Extract the logged in user data to ensure the file belongs to the logged-in user
    cookie, err := r.Cookie("token")
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    tokenStr := cookie.Value
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil || !token.Valid {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

	// Extract fileID from the request:
    vars := mux.Vars(r)
    fileID := vars["fileID"]

    var file File
    err = db.DB.QueryRow("SELECT id, filename, path, email FROM user_files WHERE id = $1", 
			fileID).Scan(&file.File_id, &file.Filename, &file.Path,
			&file.Email)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintln(w, "File not found")
        return
    }

	// Ensure the file belongs to the logged-in user:
	if file.Email != claims.Email {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "This file does not belong to you!")
		return
	}

    // Delete the file record from the database
	_, err = db.DB.Exec("DELETE FROM user_files WHERE id = $1", fileID)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintln(w, "Error deleting file from database")
        return
    }
    // Remove the file from disk
    err = os.Remove(file.Path)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintln(w, "Error deleting file")
        return
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "File deleted successfully")
}
