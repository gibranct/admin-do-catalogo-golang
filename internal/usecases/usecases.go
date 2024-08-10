package usecase

import (
	"database/sql"

	castmember "github.com.br/gibranct/admin-do-catalogo/internal/infra/castmember"
	gateway "github.com.br/gibranct/admin-do-catalogo/internal/infra/category"
	castmemberUsecase "github.com.br/gibranct/admin-do-catalogo/internal/usecases/castmember"
	categoryUsecase "github.com.br/gibranct/admin-do-catalogo/internal/usecases/category"
)

type CategoryUseCase struct {
	Create     categoryUsecase.CreateCategoryUseCase
	FindOne    categoryUsecase.GetCategoryByIdUseCase
	Activate   categoryUsecase.ActivateCategoryUseCase
	Deactivate categoryUsecase.DeactivateCategoryUseCase
	FindAll    categoryUsecase.ListCategoriesUseCase
	Update     categoryUsecase.UpdateCategoryUseCase
}

type CastMemberUseCase struct {
	Create  castmemberUsecase.CreateCastMemberUseCase
	Update  castmemberUsecase.UpdateCategoryUseCase
	FindAll castmemberUsecase.ListCastMembersUseCase
}

type UseCases struct {
	Category   CategoryUseCase
	CastMember CastMemberUseCase
}

func NewUseCases(db *sql.DB) UseCases {
	cGateway := gateway.NewCategoryGateway(db)
	cmGateway := castmember.NewCastMemberGateway(db)
	return UseCases{
		Category: CategoryUseCase{
			Create: categoryUsecase.DefaultCreateCategoryUseCase{
				Gateway: cGateway,
			},
			FindOne: &categoryUsecase.DefaultGetCategoryByIdUseCase{
				Gateway: cGateway,
			},
			Activate: &categoryUsecase.DefaultActivateCategoryUseCase{
				Gateway: cGateway,
			},
			Deactivate: &categoryUsecase.DefaultDeactivateCategoryUseCase{
				Gateway: cGateway,
			},
			FindAll: &categoryUsecase.DefaultListCategoriesUseCase{
				Gateway: cGateway,
			},
			Update: &categoryUsecase.DefaultUpdateCategoryUseCase{
				Gateway: cGateway,
			},
		},
		CastMember: CastMemberUseCase{
			Create: &castmemberUsecase.DefaultCreateCastMemberUseCase{
				Gateway: cmGateway,
			},
			Update: &castmemberUsecase.DefaultUpdateCastMemberUseCase{
				Gateway: cmGateway,
			},
			FindAll: &castmemberUsecase.DefaultListCastMembersUseCase{
				Gateway: cmGateway,
			},
		},
	}
}
