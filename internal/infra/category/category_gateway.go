package gateway

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

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

func (cg *CategoryGateway) FindById(categoryId int64) (*category.Category, error) {
	query := `
	 SELECT id, name, description, is_active, created_at, updated_at, deleted_at FROM
	 categories 
	 where id = $1
	`

	category := category.Category{}

	err := cg.Db.QueryRow(query, categoryId).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (cg *CategoryGateway) Update(c category.Category) error {
	query := `
	 UPDATE categories set name=$1, description=$2, is_active=$3, updated_at=$4, deleted_at=$5
	 where id = $6
	`

	args := []any{c.Name, c.Description, c.IsActive, c.UpdatedAt, c.DeletedAt, c.ID}

	_, err := cg.Db.Exec(query, args...)

	return err
}

func (cg *CategoryGateway) FindAll(query domain.SearchQuery) (*domain.Pagination[category.Category], error) {
	sql := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, name, description, is_active, created_at, updated_at, deleted_at
		FROM categories
		WHERE name ILIKE $1 OR description ILIKE $1
		ORDER BY %s, id %s 
		LIMIT $2 OFFSET $3`,
		query.SortColumn(), query.SortDirection())

	args := []any{"%" + query.Term + "%", query.Limit(), query.Offset()}

	rows, err := cg.Db.Query(sql, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	categories := []*category.Category{}
	totalRecords := 0

	for rows.Next() {
		var c category.Category
		err := rows.Scan(
			&totalRecords,
			&c.ID,
			&c.Name,
			&c.Description,
			&c.IsActive,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		categories = append(categories, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &domain.Pagination[category.Category]{
		Items:       categories,
		PerPage:     query.PerPage,
		CurrentPage: query.Page,
		Total:       totalRecords,
	}, nil

}

func (cg *CategoryGateway) ExistsByIds(categoryIds []int64) ([]int64, error) {
	var stringIds []string
	for _, id := range categoryIds {
		stringIds = append(stringIds, strconv.Itoa(int(id)))
	}

	query := fmt.Sprintf(`
		SELECT id from categories WHERE id IN (%s)
		ORDER BY id ASC
	`, strings.Join(stringIds, ","))

	rows, err := cg.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64

	for rows.Next() {
		var id int64
		err = rows.Scan(&id)

		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}
