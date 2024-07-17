package category_usecase

import (
	"testing"

	"github.com.br/gibranct/admin-do-catalogo/internal/domain"
	"github.com.br/gibranct/admin-do-catalogo/internal/domain/category"
	"github.com.br/gibranct/admin-do-catalogo/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFindAllCategories(t *testing.T) {
	gatewayMock := new(mocks.CategoryGatewayMock)
	sut := DefaultListCategoriesUseCase{
		Gateway: gatewayMock,
	}
	query := domain.SearchQuery{
		Page:      1,
		PerPage:   10,
		Term:      "aa",
		Sort:      "name",
		Direction: "ASC",
	}
	categories := []*category.Category{
		category.NewCategory("Cate1", "Desc 1"),
		category.NewCategory("Cate2", "Desc 2"),
		category.NewCategory("Cate3", "Desc 3"),
	}
	pageCategories := &domain.Pagination[category.Category]{
		CurrentPage: 1,
		PerPage:     10,
		Total:       1,
		Items:       categories,
	}

	gatewayMock.On("FindAll", query).Return(pageCategories, nil)

	page, err := sut.Execute(query)

	assert.Nil(t, err)
	assert.Len(t, page.Items, 3)
	assert.Equal(t, pageCategories.CurrentPage, page.CurrentPage)
	assert.Equal(t, pageCategories.PerPage, page.PerPage)
	assert.Equal(t, pageCategories.Total, page.Total)
	for idx, item := range page.Items {
		assert.Equal(t, categories[idx].Name, item.Name)
		assert.Equal(t, categories[idx].Description, item.Description)
	}
}
