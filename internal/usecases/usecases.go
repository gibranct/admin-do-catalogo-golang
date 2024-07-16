package usecase

import (
	"database/sql"

	gateway "github.com.br/gibranct/admin-do-catalogo/internal/infra/category"
	categoryUsecase "github.com.br/gibranct/admin-do-catalogo/internal/usecases/category"
)

type CategoryUseCase struct {
	Create  categoryUsecase.CreateCategoryUseCase
	FindOne categoryUsecase.GetCategoryByIdUseCase
	Delete  categoryUsecase.DeleteCategoryUseCase
}
type UseCases struct {
	Category CategoryUseCase
}

func NewUseCases(db *sql.DB) UseCases {
	cGateway := gateway.NewCategoryGateway(db)
	return UseCases{
		Category: CategoryUseCase{
			Create: categoryUsecase.DefaultCreateCategoryUseCase{
				Gateway: cGateway,
			},
			FindOne: &categoryUsecase.DefaultGetCategoryByIdUseCase{
				Gateway: cGateway,
			},
		},
	}
}
