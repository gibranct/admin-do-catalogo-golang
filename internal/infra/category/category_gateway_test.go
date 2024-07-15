package gateway

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com.br/gibranct/admin-do-catalogo/internal/domain"
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
	category := category.NewCategory("drinks", "drinks desc")
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
	category := category.NewCategory("drinks", "drinks desc")
	cg.Create(category)

	categoryFound, err := cg.FindById(category.ID)
	assert.Nil(t, err)
	assert.Equal(t, category.ID, categoryFound.ID)
	assert.Equal(t, category.Name, categoryFound.Name)
	assert.Equal(t, category.Description, categoryFound.Description)
	assert.Equal(t, category.IsActive, categoryFound.IsActive)
}

func TestUpdate(t *testing.T) {
	db, err := sql.Open("postgres", dbContainer.connectionString)
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	db.Exec("DELETE FROM categories")
	defer db.Close()
	cg := NewCategoryGateway(db)
	category := category.NewCategory("drinks", "drinks desc")
	cg.Create(category)

	updatedCategory := category.Update("new name", "new description")

	err = cg.Update(*updatedCategory)

	categoryFound, _ := cg.FindById(category.ID)

	assert.Nil(t, err)
	assert.Equal(t, updatedCategory.ID, categoryFound.ID)
	assert.Equal(t, category.Name, categoryFound.Name)
	assert.Equal(t, category.Description, categoryFound.Description)
	assert.Equal(t, category.IsActive, categoryFound.IsActive)
}

func TestFindAllWithFilters(t *testing.T) {
	db, err := sql.Open("postgres", dbContainer.connectionString)
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	db.Exec("DELETE FROM categories")
	defer db.Close()
	cg := NewCategoryGateway(db)
	category1 := category.NewCategory("movie", "drinks desc")
	category2 := category.NewCategory("tv show", "drinks desc")
	category3 := category.NewCategory("documentary", "drinks desc")
	cg.Create(category1)
	cg.Create(category2)
	cg.Create(category3)

	query := domain.SearchQuery{
		Page:      1,
		PerPage:   1,
		Term:      "doc",
		Sort:      "name",
		Direction: "ASC",
	}

	page, err := cg.FindAll(query)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(page.Items))
	assert.Equal(t, 1, page.Total)
	assert.Equal(t, 1, page.CurrentPage)
	assert.Equal(t, 1, page.PerPage)
}

func TestFindAllWithoutFilters(t *testing.T) {
	db, err := sql.Open("postgres", dbContainer.connectionString)
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	db.Exec("DELETE FROM categories")
	defer db.Close()
	cg := NewCategoryGateway(db)
	category1 := category.NewCategory("movie", "drinks desc")
	category2 := category.NewCategory("tv show", "drinks desc")
	category3 := category.NewCategory("documentary", "drinks desc")
	cg.Create(category1)
	cg.Create(category2)
	cg.Create(category3)

	query := domain.SearchQuery{
		Page:      1,
		PerPage:   10,
		Term:      "",
		Sort:      "name",
		Direction: "ASC",
	}

	page, err := cg.FindAll(query)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(page.Items))
	assert.Equal(t, 3, page.Total)
	assert.Equal(t, 1, page.CurrentPage)
	assert.Equal(t, 10, page.PerPage)
	assert.Equal(t, category3.Name, page.Items[0].Name)
	assert.Equal(t, category1.Name, page.Items[1].Name)
	assert.Equal(t, category2.Name, page.Items[2].Name)
}

func TestExistsByIds(t *testing.T) {
	db, err := sql.Open("postgres", dbContainer.connectionString)
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	db.Exec("DELETE FROM categories")
	defer db.Close()
	cg := NewCategoryGateway(db)
	category1 := category.NewCategory("movie", "drinks desc")
	category2 := category.NewCategory("tv show", "drinks desc")
	category3 := category.NewCategory("documentary", "drinks desc")
	cg.Create(category1)
	cg.Create(category2)
	cg.Create(category3)
	ids := []int64{category1.ID, category2.ID, category3.ID}

	foundIds, err := cg.ExistsByIds(ids)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(foundIds))
	assert.Equal(t, category1.ID, foundIds[0])
	assert.Equal(t, category2.ID, foundIds[1])
	assert.Equal(t, category3.ID, foundIds[2])
}
