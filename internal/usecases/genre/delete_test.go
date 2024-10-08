package genre_usecase

import (
	"errors"
	"testing"

	"github.com.br/gibranct/admin-do-catalogo/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDeleteGenreById(t *testing.T) {
	genreGatewayMock := new(mocks.GenreGatewayMock)
	sut := DefaultDeleteGenreUseCase{
		Gateway: genreGatewayMock,
	}
	genreId := int64(45)

	genreGatewayMock.On("DeleteById", genreId).Return(nil)

	err := sut.Gateway.DeleteById(genreId)

	assert.Nil(t, err)
}

func TestDeleteGenreByIdWhenFails(t *testing.T) {
	genreGatewayMock := new(mocks.GenreGatewayMock)
	sut := DefaultDeleteGenreUseCase{
		Gateway: genreGatewayMock,
	}
	genreId := int64(45)
	err := errors.New("failed to delete genre")

	genreGatewayMock.On("DeleteById", genreId).Return(err)

	result := sut.Gateway.DeleteById(genreId)

	assert.NotNil(t, err)
	assert.Equal(t, err, result)
}
