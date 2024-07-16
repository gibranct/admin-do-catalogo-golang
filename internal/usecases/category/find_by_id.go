package category_usecase

import (
	"github.com.br/gibranct/admin-do-catalogo/internal/domain/category"
)

type CategoryOutput struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"active"`
}

type GetCategoryByIdUseCase interface {
	Execute(categoryId int64) (*CategoryOutput, error)
}

type DefaultGetCategoryByIdUseCase struct {
	Gateway category.CategoryGateway
}

func (fo *DefaultGetCategoryByIdUseCase) Execute(categoryId int64) (*CategoryOutput, error) {
	cat, err := fo.Gateway.FindById(categoryId)

	if err != nil {
		return nil, err
	}

	return &CategoryOutput{
		ID:          cat.ID,
		Name:        cat.Name,
		Description: cat.Description,
		IsActive:    cat.IsActive,
	}, nil
}
