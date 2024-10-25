package mocks

import (
	"github.com.br/gibranct/admin_do_catalogo/internal/domain"
	"github.com.br/gibranct/admin_do_catalogo/internal/domain/category"
	"github.com/stretchr/testify/mock"
)

type CategoryGatewayMock struct {
	mock.Mock
}

func (m *CategoryGatewayMock) Create(c *category.Category) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *CategoryGatewayMock) FindById(categoryId int64) (*category.Category, error) {
	args := m.Called(categoryId)
	return args.Get(0).(*category.Category), args.Error(1)
}

func (m *CategoryGatewayMock) Update(c category.Category) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *CategoryGatewayMock) FindAll(query domain.SearchQuery) (*domain.Pagination[category.Category], error) {
	args := m.Called(query)
	return args.Get(0).(*domain.Pagination[category.Category]), args.Error(1)

}
func (m *CategoryGatewayMock) ExistsByIds(categoryIds []int64) ([]int64, error) {
	args := m.Called(categoryIds)
	return args.Get(0).([]int64), args.Error(1)
}
