// package main

// import (
// 	"log"
// 	"net/http"
// 	"os"
// 	"github.com/3007-Ani-Sha/distributed-storage-system/handlers"
// )

// func main() {
// 	// Create an uploads directory if it doesn't exist
// 	err := os.MkdirAll("./uploads", os.ModePerm)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Register routes
// 	http.HandleFunc("/upload", handlers.UploadHandler)
// 	http.HandleFunc("/download", handlers.DownloadHandler)

// 	// Start the server
// 	log.Println("Server started on :8080")
// 	err = http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		log.Fatal("ListenAndServe: ", err)
// 	}
// }

package main

import (
	"log"
	"github.com/joho/godotenv"
	"net/http"
	"github.com/rs/cors"
	"github.com/gorilla/mux"
	"github.com/3007-Ani-Sha/distributed-storage-system/db"
	"github.com/3007-Ani-Sha/distributed-storage-system/handlers"
)

func main() {

	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	// Initialize the database connection
	db.InitDB()

	// Create a new router
	r := mux.NewRouter()

	// Define routes for user management and file operations
	r.HandleFunc("/api/send-otp", handlers.SendOTP).Methods("POST")
	r.HandleFunc("/api/signup", handlers.SignUp).Methods("POST")
	r.HandleFunc("/api/login", handlers.Login).Methods("POST")

	// Routes for file upload and download
	r.HandleFunc("/api/upload", handlers.UploadHandler).Methods("POST")
	r.HandleFunc("/api/files", handlers.FetchFilesHandler).Methods("GET")
	// r.HandleFunc("/api/download/{filename:[a-zA-Z0-9_\\.]+}", handlers.DownloadHandler).Methods("GET")
	//Download with giving the email to use in extaction of file:
	r.HandleFunc("/api/download/{email}/{filename:[a-zA-Z0-9_\\.]+}", handlers.DownloadFileHandler).Methods("GET")
	r.HandleFunc("/api/delete/{fileID}", handlers.DeleteFileHandler).Methods("DELETE")

	// Initialize CORS middleware with allowed origins and methods
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3001"}, // Allow your frontend's origin
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // If you're using cookies or session-based authentication
	})

	handler := corsHandler.Handler(r)

	// Start the server
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
