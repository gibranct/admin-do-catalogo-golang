package models

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com.br/gibranct/admin-do-catalogo/internal/domain/category"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/lib/pq"
)

type databaseContainer struct {
	*postgres.PostgresContainer
	connectionString string
}

func setup(ctx context.Context) (*databaseContainer, error) {
	container, err := postgres.Run(
		ctx,
		"postgres:15.7-alpine",
		postgres.WithInitScripts("../../../migrations/000001_create_categories_table.up.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}

	connString, err := container.ConnectionString(ctx)

	if err != nil {
		log.Fatalf("Could not get connection string: %s", err)
	}

	return &databaseContainer{
		PostgresContainer: container,
		connectionString:  connString + "sslmode=disable",
	}, nil
}

func TestCreateCategory(t *testing.T) {
	ctx := context.Background()
	container, _ := setup(ctx)
	db, err := sql.Open("postgres", container.connectionString)
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()
	defer container.Terminate(ctx)
	cg := NewCategoryGateway(db)
	category := category.NewCategory("drinks", "drinks desc", true)
	err = cg.Create(category)
	if err != nil {
		log.Fatalf("Could not save category: %s", err)
	}
	assert.Nil(t, err)
	assert.NotNil(t, category.ID)
	assert.NotNil(t, category.CreatedAt)
	assert.NotNil(t, category.UpdatedAt)
	assert.Equal(t, time.Time{}, category.DeletedAt)
}
