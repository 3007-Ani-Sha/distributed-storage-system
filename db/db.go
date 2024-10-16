package db

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {

    // connStr := os.Getenv("DATABASE_URL")
	connStr :="postgresql://postgres.vwbonkhkfkkzybjtzotl:dss@dell@1107@aws-0-ap-south-1.pooler.supabase.com:6543/postgres"
    if connStr == "" {
        log.Fatal("DATABASE_URL environment variable not set")
    }
	
    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to connect to Supabase: ", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatal("Failed to ping the database: ", err)
    }

    fmt.Println("Successfully connected to the database!")
}
