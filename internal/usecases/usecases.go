package usecase

import (
	"database/sql"

	castmember "github.com.br/gibranct/admin-do-catalogo/internal/infra/castmember"
	gateway "github.com.br/gibranct/admin-do-catalogo/internal/infra/category"
	infra_genre "github.com.br/gibranct/admin-do-catalogo/internal/infra/genre"
	infra_video "github.com.br/gibranct/admin-do-catalogo/internal/infra/video"
	castmemberUsecase "github.com.br/gibranct/admin-do-catalogo/internal/usecases/castmember"
	categoryUsecase "github.com.br/gibranct/admin-do-catalogo/internal/usecases/category"
	genre_usecase "github.com.br/gibranct/admin-do-catalogo/internal/usecases/genre"
	video_usecase "github.com.br/gibranct/admin-do-catalogo/internal/usecases/video"
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

type GenreUseCase struct {
	Create     genre_usecase.CreateGenreUseCase
	FindAll    genre_usecase.ListGenresUseCase
	DeleteById genre_usecase.DeleteGenreUseCase
}

type VideoUseCase struct {
	Create video_usecase.CreateVideoUseCase
}

type UseCases struct {
	Category   CategoryUseCase
	CastMember CastMemberUseCase
	Genre      GenreUseCase
	Video      VideoUseCase
}

func NewUseCases(db *sql.DB) UseCases {
	cGateway := gateway.NewCategoryGateway(db)
	cmGateway := castmember.NewCastMemberGateway(db)
	gGateway := infra_genre.NewGenreGateway(db)
	vg := infra_video.NewVideoGateway(db)
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
		Genre: GenreUseCase{
			Create: genre_usecase.DefaultCreateGenreUseCase{
				Gateway:         gGateway,
				CategoryGateway: cGateway,
			},
			FindAll: genre_usecase.DefaultListGenresUseCase{
				Gateway: gGateway,
			},
			DeleteById: genre_usecase.DefaultDeleteGenreUseCase{
				Gateway: gGateway,
			},
		},
		Video: VideoUseCase{
			Create: video_usecase.DefaultCreateVideoUseCase{
				Gateway:           vg,
				CategoryGateway:   cGateway,
				GenreGateway:      gGateway,
				CastMemberGateway: cmGateway,
			},
		},
	}
}
