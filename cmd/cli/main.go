package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/funchooooza-ossh/protego/internal/composition"
	"github.com/funchooooza-ossh/protego/internal/config"
	"github.com/funchooooza-ossh/protego/internal/domain"
	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "register" {
		fmt.Println("usage: register --login=<login> --password=<password>")
		os.Exit(1)
	}

	// CLI flags
	registerCmd := flag.NewFlagSet("register", flag.ExitOnError)
	login := registerCmd.String("login", "", "Login")
	password := registerCmd.String("password", "", "Password")
	registerCmd.Parse(os.Args[2:])

	if *login == "" || *password == "" {
		fmt.Println("Both --login and --password are required")
		os.Exit(1)
	}

	// Load config and init
	cfg := config.Load()

	// Connect infrastructure (DB, Redis, etc.)
	conns := composition.ProvideConnections(cfg)

	// Apply DB migrations
	if err := composition.ApplyMigrations(cfg.DatabaseDsn(), "./migrations"); err != nil {
		fmt.Println("Migration error:", err)
		os.Exit(1)
	}

	// Prepare hashed password
	hashed, err := bcrypt.GenerateFromPassword([]byte(*password), cfg.PassCost)
	if err != nil {
		fmt.Println("Invalid password:", err)
		os.Exit(1)
	}

	repos := composition.ProvideAdapters(conns, cfg)
	ctx := context.Background()

	// Prevent duplicate superuser
	existing, err := repos.User.GetByEmail(ctx, *login)
	if err == nil && existing != nil {
		fmt.Println("Superuser already exists.")
		os.Exit(1)
	}

	// Ensure supermanager role exists
	role, err := repos.Role.GetByCode(ctx, "supermanager")
	if errors.Is(err, e.ErrNotFound) || role == nil {
		role = &domain.Role{
			ID:   uuid.NewString(),
			Code: "supermanager",
		}
		if err := repos.Role.Create(ctx, role); err != nil {
			fmt.Println("Role creation error:", err)
			os.Exit(1)
		}
	} else if err != nil {
		fmt.Println("Role fetch error:", err)
		os.Exit(1)
	}

	// Create superuser
	user := &domain.User{
		ID:       uuid.NewString(),
		Email:    *login,
		Password: string(hashed),
		Blocked:  false,
		Role:     *role,
	}
	if err := repos.User.Create(ctx, user); err != nil {
		fmt.Println("User creation error:", err)
		os.Exit(1)
	}

	fmt.Printf("Superuser created: %s\n", *login)
}
