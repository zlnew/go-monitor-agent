package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	storagePg "horizonx-server/internal/storage/postgres"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	defaultDSN := os.Getenv("DATABASE_URL")

	cmd := flag.String("op", "", "operation: up, down, version, force")
	steps := flag.Int("steps", 0, "number of steps for up/down (0 = all)")
	dsn := flag.String("dsn", defaultDSN, "database url (postgres://user:pass@host:port/db)")
	flag.Parse()

	if *cmd == "" || *dsn == "" {
		fmt.Println("Usage: go run cmd/migrate/main.go -op=[up|down] -dsn=[postgres://...]")
		os.Exit(1)
	}

	db, err := sql.Open("pgx", *dsn)
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("could not create driver: %v", err)
	}

	src, err := iofs.New(storagePg.MigrationsFS, "migrations")
	if err != nil {
		log.Fatalf("could not create source driver: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		log.Fatalf("could not create migrate instance: %v", err)
	}

	log.Printf("Running migration op: %s...", *cmd)
	switch *cmd {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
	case "down":
		if *steps > 0 {
			err = m.Steps(-(*steps))
		} else {
			err = m.Down()
		}
	case "version":
		v, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			log.Fatal(err)
		}
		fmt.Printf("Version: %d, Dirty: %v\n", v, dirty)
		return
	case "force":
		if *steps == 0 {
			log.Fatal("please specify version to force")
		}
		err = m.Force(*steps)
	default:
		log.Fatal("unknown command")
	}

	if err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No changes detected.")
		} else {
			log.Fatalf("Migration failed: %v", err)
		}
	} else {
		fmt.Println("Migration success!")
	}
}
