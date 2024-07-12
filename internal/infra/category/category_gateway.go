package models

import (
	"database/sql"

	"github.com.br/gibranct/admin-do-catalogo/internal/domain"
	"github.com.br/gibranct/admin-do-catalogo/internal/domain/category"
)

type CategoryGateway struct {
	Db *sql.DB
}

func NewCategoryGateway(db *sql.DB) *CategoryGateway {
	return &CategoryGateway{Db: db}
}

func (cg *CategoryGateway) Create(c *category.Category) error {
	query := `
		INSERT INTO categories (name, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	args := []any{c.Name, c.Description, c.IsActive, c.CreatedAt, c.UpdatedAt}

	return cg.Db.QueryRow(query, args...).Scan(&c.ID)
}

func (cg *CategoryGateway) DeleteById(categoryId int64) error {
	query := `
	 DELETE from categories where id = $1
	`
	_, err := cg.Db.Exec(query, categoryId)
	return err
}

func (gate *CategoryGateway) FindById(categoryId int64) (category.Category, error) {
	return *&category.Category{}, nil
}

func (gate *CategoryGateway) Update(c category.Category) (category.Category, error) {
	return *&category.Category{}, nil
}

func (gate *CategoryGateway) FindAll(query domain.SearchQuery) domain.Pagination[category.Category] {
	return domain.Pagination[category.Category]{}
}

func (gate *CategoryGateway) ExistsByIds(categoryIds []int64) []int64 {
	return []int64{1}
}
