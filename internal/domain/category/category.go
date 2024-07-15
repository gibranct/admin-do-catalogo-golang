package category

import (
	"time"

	"github.com.br/gibranct/admin-do-catalogo/internal/domain"
	"github.com.br/gibranct/admin-do-catalogo/pkg/validator"
)

type Category struct {
	ID          int64
	Name        string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type CategoryGateway interface {
	Create(category *Category) error
	DeleteById(categoryId int64) error
	FindById(categoryId int64) (*Category, error)
	Update(category Category) error
	FindAll(query domain.SearchQuery) (*domain.Pagination[Category], error)
	ExistsByIds(categoryIds []int64) ([]int64, error)
}

func NewCategory(
	name string,
	description string,
) *Category {
	now := time.Now().UTC()
	return &Category{
		Name:        name,
		Description: description,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   nil,
	}
}

func (c *Category) Deactivate() *Category {
	now := time.Now().UTC()
	c.DeletedAt = &now
	c.IsActive = false
	c.UpdatedAt = now
	return c
}

func (c *Category) Activate() *Category {
	c.DeletedAt = nil
	c.IsActive = true
	c.UpdatedAt = time.Now().UTC()
	return c
}

func (c *Category) Update(
	name, description string,
) *Category {
	c.Name = name
	c.Description = description
	c.UpdatedAt = time.Now().UTC()
	return c
}

func (c *Category) Validate(handler validator.ValidationHandler) {
	NewCategoryValidator(*c, handler).Validate()
}
