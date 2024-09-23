package test

import (
	"context"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var scripts = []string{
	"../../migrations/000001_create_categories_table.up.sql",
	"../../migrations/000002_create_cast_members_table.up.sql",
	"../../migrations/000003_create_genres_table.sql.up.sql",
	"../../migrations/000004_create_videos_table.sql.up.sql",
}

func InitDatabase(ctx context.Context) (string, *postgres.PostgresContainer, error) {
	container, err := postgres.Run(
		ctx,
		"postgres:15.7-alpine",
		postgres.WithInitScripts(scripts...),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}

	connString, err := container.ConnectionString(ctx, "sslmode=disable")

	if err != nil {
		log.Fatalf("Could not get connection string: %s", err)
	}

	return connString, container, nil
}
