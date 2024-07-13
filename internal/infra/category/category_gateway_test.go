package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
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

var dbContainer *databaseContainer

func TestMain(t *testing.M) {
	ctx := context.Background()
	defer func() {
		if r := recover(); r != nil {
			dbContainer.PostgresContainer.Terminate(ctx)
			fmt.Println("Panic")
		}
	}()
	setup(ctx)
	code := t.Run()
	dbContainer.PostgresContainer.Terminate(ctx)
	os.Exit(code)
}

func setup(ctx context.Context) error {
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

	connString, err := container.ConnectionString(ctx, "sslmode=disable")

	if err != nil {
		log.Fatalf("Could not get connection string: %s", err)
	}

	dbContainer = &databaseContainer{
		connectionString:  connString,
		PostgresContainer: container,
	}

	return nil
}

func TestCreateCategory(t *testing.T) {
	db, err := sql.Open("postgres", dbContainer.connectionString)
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	db.Exec("DELETE FROM categories")
	defer db.Close()
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
	assert.Nil(t, category.DeletedAt)
}

func TestFindById(t *testing.T) {
	db, err := sql.Open("postgres", dbContainer.connectionString)
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	db.Exec("DELETE FROM categories")
	defer db.Close()
	cg := NewCategoryGateway(db)
	category := category.NewCategory("drinks", "drinks desc", true)
	cg.Create(category)

	categoryFound, err := cg.FindById(category.ID)
	assert.Nil(t, err)
	assert.Equal(t, category.ID, categoryFound.ID)
	assert.Equal(t, category.Name, categoryFound.Name)
	assert.Equal(t, category.Description, categoryFound.Description)
	assert.Equal(t, category.IsActive, categoryFound.IsActive)
}
