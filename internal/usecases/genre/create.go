package genre_usecase

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain/category"
	"github.com.br/gibranct/admin_do_catalogo/internal/domain/genre"
	"github.com.br/gibranct/admin_do_catalogo/pkg/notification"
)

type CreateGenreOutput struct {
	ID int64
}

type CreateGenreCommand struct {
	Name        string
	CategoryIds *[]int64
}

type CreateGenreUseCase interface {
	Execute(c CreateGenreCommand) (*notification.Notification, *CreateGenreOutput)
}

type DefaultCreateGenreUseCase struct {
	Gateway         genre.GenreGateway
	CategoryGateway category.CategoryGateway
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

	err := useCase.ValidateCategories(*command.CategoryIds)

	if err != nil {
		n.Add(err)
		return n, nil
	}
	genre.AddCategoriesIds(*command.CategoryIds)

	err = useCase.Gateway.Create(genre)

	if err != nil {
		n.Add(err)
		return n, nil
	}

	return nil, &CreateGenreOutput{
		ID: genre.ID,
	}
}

func (useCase DefaultCreateGenreUseCase) ValidateCategories(categoriesIds []int64) error {
	if len(categoriesIds) == 0 {
		return nil
	}

	ids, err := useCase.CategoryGateway.ExistsByIds(categoriesIds)
	if err != nil {
		return err
	}

	if slices.Equal(ids, categoriesIds) {
		return nil
	}

	var missingIds []string

	for _, id := range categoriesIds {
		if !slices.Contains(ids, id) {
			missingIds = append(missingIds, strconv.FormatInt(id, 10))
		}
	}

	if len(missingIds) != 0 {
		return fmt.Errorf("missing category ids: %s", strings.Join(missingIds, ","))
	}

	return nil
}
