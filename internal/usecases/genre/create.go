package genre_usecase

import (
	"github.com.br/gibranct/admin-do-catalogo/internal/domain/genre"
	"github.com.br/gibranct/admin-do-catalogo/pkg/notification"
)

type CreateGenreOutput struct {
	ID int64
}

type CreateGenreCommand struct {
	Name string
}

type CreateGenreUseCase interface {
	Execute(c CreateGenreCommand) (*notification.Notification, *CreateGenreOutput)
}

type DefaultCreateGenreUseCase struct {
	Gateway genre.GenreGateway
}

func (useCase DefaultCreateGenreUseCase) Execute(
	command CreateGenreCommand,
) (*notification.Notification, *CreateGenreOutput) {

	genre := genre.NewGenre(command.Name)

	n := notification.CreateNotification()

	genre.Validate(n)

	if n.HasErrors() {
		return n, nil
	}

	err := useCase.Gateway.Create(genre)

	if err != nil {
		n.Add(err)
		return n, nil
	}

	return nil, &CreateGenreOutput{
		ID: genre.ID,
	}
}
