package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() *sql.DB {
	defaultDB := connectToServer()
	createDatabaseIfNotExists(defaultDB)
	defaultDB.Close()

	appDB := connectToAppDB()
	createTableIfNotExists(appDB)
	DB = appDB
	return appDB
}

func connectToServer() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if port == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("Missing one or more required environment variables (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)")
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, user)
	serverDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Cannot connect to PostgreSQL server: %v", err)
	}
	if err = serverDB.Ping(); err != nil {
		log.Fatalf("PostgreSQL server not reachable: %v", err)
	}
	fmt.Println("Connected to PostgreSQL server (default db).")
	return serverDB
}

func createDatabaseIfNotExists(db *sql.DB) {
	dbname := os.Getenv("DB_NAME")

	query := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", dbname)

	var exists int
	err := db.QueryRow(query).Scan(&exists)

	if err == sql.ErrNoRows {
		_, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
		if err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
		fmt.Printf("Database '%s' created successfully.\n", dbname)
	} else if err != nil {
		log.Fatalf("Database existence check failed: %v", err)
	} else {
		fmt.Printf("Database '%s' already exists.\n", dbname)
	}
}

func connectToAppDB() *sql.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s application_name=GoServer sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Could not connect to app database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Could not ping app database: %v", err)
	}

	fmt.Printf("Connected to '%s' successfully.\n", dbname)
	configurePool(db)
	return db
}

func createTableIfNotExists(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS kvstore (
		key INT PRIMARY KEY,
		value TEXT NOT NULL
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	fmt.Println("Table 'kvstore' ready.")
}

func configurePool(db *sql.DB) {
	maxOpen := intFromEnv("DB_MAX_OPEN_CONNS", 50)
	maxIdle := intFromEnv("DB_MAX_IDLE_CONNS", 25)
	lifetimeMinutes := intFromEnv("DB_CONN_MAX_LIFETIME_MIN", 5)

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(time.Duration(lifetimeMinutes) * time.Minute)
}

func intFromEnv(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return fallback
}
