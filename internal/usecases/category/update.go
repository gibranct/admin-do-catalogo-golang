package category_usecase

import (
	"errors"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain/category"
	"github.com.br/gibranct/admin_do_catalogo/pkg/notification"
)

type UpdateCategoryCommand struct {
	ID          int64
	Name        string
	Description string
}

type UpdateCategoryUseCase interface {
	Execute(c UpdateCategoryCommand) *notification.Notification
}

type DefaultUpdateCategoryUseCase struct {
	Gateway category.CategoryGateway
}

func (useCase DefaultUpdateCategoryUseCase) Execute(
	command UpdateCategoryCommand,
) *notification.Notification {

	category, err := useCase.Gateway.FindById(command.ID)

	n := notification.CreateNotification()

	if err != nil {
		n.Add(errors.New("category not found"))
		return n
	}

	category.Update(command.Name, command.Description)

	category.Validate(n)

	if n.HasErrors() {
		return n
	}

	err = useCase.Gateway.Update(*category)

	if err != nil {
		n.Add(err)
		return n
	}

	return nil
}
