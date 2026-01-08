package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"horizonx/internal/adapters/postgres"
	"horizonx/internal/application/role"
	"horizonx/internal/application/user"
	"horizonx/internal/domain"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("WARNING! .env file not found")
	}

	dbURL := os.Getenv("DATABASE_URL")

	dbPool, err := postgres.InitDB(dbURL)
	if err != nil {
		log.Fatalf("failed to init DB: %v", err)
	}
	defer dbPool.Close()

	roleRepo := postgres.NewRoleRepository(dbPool)
	roleSvc := role.NewService(roleRepo)

	userRepo := postgres.NewUserRepository(dbPool)
	userSvc := user.NewService(userRepo)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	seedRolePermissions(ctx, roleSvc)
	seedAdmin(ctx, userSvc)
}

func seedRolePermissions(ctx context.Context, role domain.RoleService) {
	if err := role.SyncPermissions(ctx); err != nil {
		log.Printf("[x] failed to seed roles and permissions: %v", err)
		return
	}

	log.Println("[v] Roles and Permissions seeded")
}

func seedAdmin(ctx context.Context, user domain.UserService) {
	email := os.Getenv("DB_ADMIN_EMAIL")
	if email == "" {
		email = "admin@horizonx.local"
	}

	password := os.Getenv("DB_ADMIN_PASSWORD")
	if password == "" {
		password = "password"
	}

	req := domain.UserSaveRequest{
		Name:     "Admin",
		Email:    email,
		Password: password,
		RoleID:   1,
	}

	if err := user.Create(ctx, req); err != nil {
		log.Printf("[x] failed to seed admin: %v", err)
		return
	}

	log.Printf("[v] Admin seeded | User: %s | Pass: %s", email, password)
}
