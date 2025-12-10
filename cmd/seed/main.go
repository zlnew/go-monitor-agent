package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found")
	}

	defaultDSN := os.Getenv("DATABASE_URL")
	dsn := flag.String("dsn", defaultDSN, "database url")
	flag.Parse()

	if *dsn == "" {
		log.Fatal("DSN required via flag -dsn or DATABASE_URL env")
	}

	db, err := sql.Open("pgx", *dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Cannot ping DB:", err)
	}

	seedAdmin(db)
	seedServer(db)
}

func seedAdmin(db *sql.DB) {
	email := "admin@horizonx.local"
	rawPassword := "password"

	if envEmail := os.Getenv("DB_ADMIN_EMAIL"); envEmail != "" {
		email = envEmail
	}

	if envPass := os.Getenv("DB_ADMIN_PASSWORD"); envPass != "" {
		rawPassword = envPass
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)

	query := `
		INSERT INTO users (name, email, password, role_id) 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) DO UPDATE SET password = excluded.password;
	`

	_, err := db.Exec(query, "Super Admin", email, string(hashed), 1)
	if err != nil {
		log.Fatalf("Failed to seed admin: %v", err)
	}

	fmt.Printf("✅ User Seeded!\n   User: %s\n   Pass: %s\n", email, rawPassword)
}

func seedServer(db *sql.DB) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM servers").Scan(&count)
	if err != nil {
		log.Printf("⚠️  Failed to check existing servers: %v", err)
		return
	}

	if count > 0 {
		fmt.Println("ℹ️  Server table populated. Skipping default server creation.")
		return
	}

	bytes := make([]byte, 32)
	rand.Read(bytes)
	token := "hzx_" + hex.EncodeToString(bytes)

	query := `
		INSERT INTO servers (name, ip_address, api_token, is_online, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id
	`

	var id int64
	err = db.QueryRow(query, "Primary Node (Supernova)", "127.0.0.1", token, false).Scan(&id)
	if err != nil {
		log.Fatalf("Failed to seed server: %v", err)
	}

	fmt.Println("---------------------------------------------------")
	fmt.Println("✅ Supernova Server Created!")
	fmt.Printf("   ID:    %d\n", id)
	fmt.Printf("   Name:  Primary Node (Supernova)\n")
	fmt.Printf("   TOKEN: %s\n", token)
	fmt.Println("---------------------------------------------------")
	fmt.Println("👉 COPY THE TOKEN ABOVE TO YOUR .env AGENT NOW! 👈")
	fmt.Println("---------------------------------------------------")
}
