package mocks

import (
	"github.com.br/gibranct/admin_do_catalogo/internal/domain/genre"
	"github.com/stretchr/testify/mock"
)

type GenreGatewayMock struct {
	mock.Mock
}

func (m *GenreGatewayMock) Create(g *genre.Genre) error {
	args := m.Called(g)
	return args.Error(0)
}

func (m *GenreGatewayMock) FindAll() ([]*genre.Genre, error) {
	args := m.Called()
	return args.Get(0).([]*genre.Genre), args.Error(1)

}
func (m *GenreGatewayMock) ExistsByIds(genreIds []int64) ([]int64, error) {
	args := m.Called(genreIds)
	return args.Get(0).([]int64), args.Error(1)
}

func (m *GenreGatewayMock) DeleteById(genreId int64) error {
	args := m.Called(genreId)
	return args.Error(0)
}
