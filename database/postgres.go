package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := "host=localhost port=5432 user=postgres password=mysecretpassword dbname=snippet_manager sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	initQuery := `
    CREATE EXTENSION IF NOT EXISTS "pgcrypto";

    -- Create the users table
    CREATE TABLE IF NOT EXISTS users (
        user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        username TEXT NOT NULL UNIQUE,
        email TEXT NOT NULL UNIQUE,
        password_hash TEXT NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
    );

    -- Create the snippets table with a foreign key referencing users
    CREATE TABLE IF NOT EXISTS snippets (
        snippet_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        title TEXT NOT NULL,
        language TEXT NOT NULL,
        content TEXT NOT NULL,
        user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
    );
  `
	_, err = DB.Exec(initQuery)
	if err != nil {
		log.Fatal("Failed to create tables: ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database: ", err)
	}
	log.Println("Connected to database!")
}

func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Println("Failed to close the database connection: ", err)
		}
	}
}
