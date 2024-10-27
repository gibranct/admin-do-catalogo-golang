package genre_usecase_test

import (
	"errors"
	"testing"

	genre_usecase "github.com.br/gibranct/admin_do_catalogo/internal/usecases/genre"
	"github.com.br/gibranct/admin_do_catalogo/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDeleteGenreById(t *testing.T) {
	genreGatewayMock := new(mocks.GenreGatewayMock)
	sut := genre_usecase.DefaultDeleteGenreUseCase{
		Gateway: genreGatewayMock,
	}
	genreId := int64(45)

	genreGatewayMock.On("DeleteById", genreId).Return(nil)

	err := sut.Gateway.DeleteById(genreId)

	assert.Nil(t, err)
}

func TestDeleteGenreByIdWhenFails(t *testing.T) {
	genreGatewayMock := new(mocks.GenreGatewayMock)
	sut := genre_usecase.DefaultDeleteGenreUseCase{
		Gateway: genreGatewayMock,
	}
	genreId := int64(45)
	err := errors.New("failed to delete genre")

	genreGatewayMock.On("DeleteById", genreId).Return(err)

	result := sut.Gateway.DeleteById(genreId)

	assert.NotNil(t, err)
	assert.Equal(t, err, result)
}
