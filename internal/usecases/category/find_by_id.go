package category_usecase

import (
	"fmt"
	"time"

	"github.com.br/gibranct/admin-do-catalogo/internal/domain"
	"github.com.br/gibranct/admin-do-catalogo/internal/domain/category"
)

type CategoryOutput struct {
	ID          int64
	Name        string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type GetCategoryByIdUseCase interface {
	Execute(categoryId int64) (CategoryOutput, error)
}

type DefaultGetCategoryByIdUseCase struct {
	gateway category.CategoryGateway
}

func (fo *DefaultGetCategoryByIdUseCase) Execute(categoryId int64) (*CategoryOutput, error) {

	cat, err := fo.gateway.FindById(categoryId)

	if err != nil {
		return nil, domain.NotFoundException{
			Message: fmt.Sprintf("Could not find category for id: %d", categoryId),
		}
	}

	return &CategoryOutput{
		ID:          cat.ID,
		Name:        cat.Name,
		Description: cat.Description,
		IsActive:    cat.IsActive,
		CreatedAt:   cat.CreatedAt,
		UpdatedAt:   cat.UpdatedAt,
		DeletedAt:   cat.DeletedAt,
	}, nil
}
