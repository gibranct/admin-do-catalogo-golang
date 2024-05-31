package category

import (
	"github.com.br/gibranct/admin-do-catalogo/internal/domain/category"
	"github.com.br/gibranct/admin-do-catalogo/pkg/notification"
)

type CreateCategoryOutput struct {
	ID int64
}

type CreateCategoryCommand struct {
	Name        string
	Description string
	IsActive    bool
}

type CreateCategoryUseCase interface {
	Execute(c CreateCategoryCommand) (notification.Notification, CreateCategoryOutput)
}

type DefaultCreateCategoryUseCase struct {
	gateway category.CategoryGateway
}

func (useCase *DefaultCreateCategoryUseCase) Execute(
	command CreateCategoryCommand,
) (*notification.Notification, *CreateCategoryOutput) {

	category := category.NewCategory(
		command.Name,
		command.Description,
		command.IsActive,
	)

	n := notification.CreateNotification()

	category.Validate(n)

	if n.HasErrors() {
		return n, nil
	}

	nCategory, err := useCase.gateway.Create(*category)

	if err != nil {
		n.Add(err)
		return n, nil
	}

	return nil, &CreateCategoryOutput{
		ID: nCategory.ID,
	}
}
