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
	DeletedAt   time.Time
}

type CategoryGateway interface {
	Create(category Category) (Category, error)
	DeleteById(categoryId string) error
	FindById(categoryId string) (Category, error)
	Update(category Category) (Category, error)
	FindAll(query domain.SearchQuery) domain.Pagination[Category]
	ExistsByIds(categoryIds []string) []string
}

func NewCategory(
	name string,
	description string,
	isActive bool,
) *Category {
	now := time.Now()
	var deletedAt time.Time
	if !isActive {
		deletedAt = now
	}
	return &Category{
		Name:        name,
		Description: description,
		IsActive:    isActive,
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   deletedAt,
	}
}

func (c *Category) Deactivate() *Category {
	if (c.DeletedAt == time.Time{}) {
		c.DeletedAt = time.Now()
	}
	c.DeletedAt = time.Now()
	c.IsActive = false
	c.UpdatedAt = time.Now()
	return c
}

func (c *Category) Activate() *Category {
	c.DeletedAt = time.Time{}
	c.IsActive = true
	c.UpdatedAt = time.Now()
	return c
}

func (c *Category) Update(
	name, description string,
	isActive bool,
) *Category {
	if isActive {
		c.Activate()
	} else {
		c.Deactivate()
	}
	c.Name = name
	c.Description = description
	c.UpdatedAt = time.Now()
	return c
}

func (c *Category) Validate(handler validator.ValidationHandler) {
	NewCategoryValidator(*c, handler).Validate()
}
