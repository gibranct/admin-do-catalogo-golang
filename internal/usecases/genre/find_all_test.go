package genre_usecase_test

import (
	"testing"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain/genre"
	genre_usecase "github.com.br/gibranct/admin_do_catalogo/internal/usecases/genre"
	"github.com.br/gibranct/admin_do_catalogo/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFindAllGenres(t *testing.T) {
	gatewayMock := new(mocks.GenreGatewayMock)
	sut := genre_usecase.DefaultListGenresUseCase{
		Gateway: gatewayMock,
	}
	genres := []*genre.Genre{
		genre.NewGenre("Genre 1"),
		genre.NewGenre("Genre 2"),
		genre.NewGenre("Genre 3"),
	}

	gatewayMock.On("FindAll").Return(genres, nil)

	page, err := sut.Execute()

	assert.Nil(t, err)
	assert.Len(t, page, 3)
	for idx, item := range page {
		assert.Equal(t, genres[idx].Name, item.Name)
		assert.Equal(t, genres[idx].CategoryIds, item.CategoryIds)
	}
}
