package genre_usecase

import (
	"github.com.br/gibranct/admin-do-catalogo/internal/domain/genre"
)

type ListGenresOutput struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Active      bool    `json:"active"`
	CategoryIds []int64 `json:"categoryIds"`
}

type ListGenresUseCase interface {
	Execute() ([]*ListGenresOutput, error)
}

type DefaultListGenresUseCase struct {
	Gateway genre.GenreGateway
}

func (useCase DefaultListGenresUseCase) Execute() ([]*ListGenresOutput, error) {
	genres, err := useCase.Gateway.FindAll()

	if err != nil {
		return nil, err
	}

	var outputs []*ListGenresOutput

	for _, item := range genres {
		output := &ListGenresOutput{
			ID:          item.ID,
			Name:        item.Name,
			Active:      item.IsActive,
			CategoryIds: item.CategoryIds,
		}

		outputs = append(outputs, output)
	}

	return outputs, nil

}
