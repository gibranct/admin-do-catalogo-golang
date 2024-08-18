package genre

import (
	"time"

	"github.com.br/gibranct/admin-do-catalogo/pkg/validator"
)

type Genre struct {
	ID          int64
	Name        string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	CategoryIds []int64
}

type GenreGateway interface {
	Create(genre *Genre) error
	FindAll() ([]*Genre, error)
	ExistsByIds(genreIds []int64) ([]int64, error)
}

func NewGenre(
	name string,
) *Genre {
	now := time.Now().UTC()
	return &Genre{
		Name:        name,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   nil,
		CategoryIds: make([]int64, 0, 5),
	}
}

func (c *Genre) Deactivate() *Genre {
	now := time.Now().UTC()
	c.DeletedAt = &now
	c.IsActive = false
	c.UpdatedAt = now
	return c
}

func (c *Genre) Activate() *Genre {
	c.DeletedAt = nil
	c.IsActive = true
	c.UpdatedAt = time.Now().UTC()
	return c
}

func (c *Genre) Update(name string) *Genre {
	c.Name = name
	c.UpdatedAt = time.Now().UTC()
	return c
}

func (c *Genre) Validate(handler validator.ValidationHandler) {
	NewGenreValidator(*c, handler).Validate()
}

func (c *Genre) AddCategoryId(categoryID int64) *Genre {
	c.CategoryIds = append(c.CategoryIds, categoryID)
	c.UpdatedAt = time.Now().UTC()
	return c
}

func (c *Genre) AddCategoriesIds(categoriesIDs []int64) *Genre {
	c.CategoryIds = append(c.CategoryIds, categoriesIDs...)
	c.UpdatedAt = time.Now().UTC()
	return c
}

func (c *Genre) RemoveCategoryId(categoryID int64) *Genre {
	var idx int
	for i := 0; i < len(c.CategoryIds); i++ {
		if c.CategoryIds[i] == categoryID {
			idx = i
		}
	}
	c.CategoryIds = append(c.CategoryIds[:idx], c.CategoryIds[idx+1:]...)
	c.UpdatedAt = time.Now().UTC()
	return c
}
