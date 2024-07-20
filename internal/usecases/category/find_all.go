package category_usecase

import (
	"github.com.br/gibranct/admin-do-catalogo/internal/domain"
	"github.com.br/gibranct/admin-do-catalogo/internal/domain/category"
)

type ListCategoriesOutput struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ListCategoriesUseCase interface {
	Execute(query domain.SearchQuery) (*domain.Pagination[ListCategoriesOutput], error)
}

type DefaultListCategoriesUseCase struct {
	Gateway category.CategoryGateway
}

func (useCase *DefaultListCategoriesUseCase) Execute(query domain.SearchQuery) (*domain.Pagination[ListCategoriesOutput], error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}

	page, err := useCase.Gateway.FindAll(query)

	if err != nil {
		return nil, err
	}

	var outputs []*ListCategoriesOutput

	for _, item := range page.Items {
		output := &ListCategoriesOutput{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
		}

		outputs = append(outputs, output)
	}

	return &domain.Pagination[ListCategoriesOutput]{
		Items:       outputs,
		CurrentPage: page.CurrentPage,
		PerPage:     page.PerPage,
		Total:       page.Total,
	}, nil

}
