package genre_usecase

import (
	"time"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain/genre"
)

type ListGenresOutput struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Active      bool      `json:"active"`
	CategoryIds []int64   `json:"categoryIds"`
	CreatedAt   time.Time `json:"createdAt"`
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
			CreatedAt:   item.CreatedAt,
		}

		outputs = append(outputs, output)
	}

	return outputs, nil

}
