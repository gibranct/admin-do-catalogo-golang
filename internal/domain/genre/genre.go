package genre

import (
	"time"

	"github.com.br/gibranct/admin-do-catalogo/pkg/validator"
)

type Genre struct {
	ID        int64
	Name      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
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
		Name:      name,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
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
