package genre_usecase

import (
	"github.com.br/gibranct/admin_do_catalogo/internal/domain/genre"
)

type DeleteGenreCommand struct {
	GenreId int64
}

type DeleteGenreUseCase interface {
	Execute(c DeleteGenreCommand) error
}

type DefaultDeleteGenreUseCase struct {
	Gateway genre.GenreGateway
}

func (useCase DefaultDeleteGenreUseCase) Execute(command DeleteGenreCommand) error {
	return useCase.Gateway.DeleteById(command.GenreId)
}
