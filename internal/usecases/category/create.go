package category_usecase

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
}

type CreateCategoryUseCase interface {
	Execute(c CreateCategoryCommand) (*notification.Notification, *CreateCategoryOutput)
}

type DefaultCreateCategoryUseCase struct {
	Gateway category.CategoryGateway
}

func (useCase DefaultCreateCategoryUseCase) Execute(
	command CreateCategoryCommand,
) (*notification.Notification, *CreateCategoryOutput) {

	category := category.NewCategory(
		command.Name,
		command.Description,
	)

	n := notification.CreateNotification()

	category.Validate(n)

	if n.HasErrors() {
		return n, nil
	}

	err := useCase.Gateway.Create(category)

	if err != nil {
		n.Add(err)
		return n, nil
	}

	return nil, &CreateCategoryOutput{
		ID: category.ID,
	}
}
