package gateway

import (
	"database/sql"

	infra_castmember "github.com.br/gibranct/admin-do-catalogo/internal/infra/castmember"
	infra_category "github.com.br/gibranct/admin-do-catalogo/internal/infra/category"
	infra_genre "github.com.br/gibranct/admin-do-catalogo/internal/infra/genre"
)

type Gateways struct {
	Category   infra_category.CategoryGateway
	CastMember infra_castmember.CastMemberGateway
	Genre      infra_genre.GenreGateway
}

func NewGateways(db *sql.DB) Gateways {
	return Gateways{
		Category:   *infra_category.NewCategoryGateway(db),
		CastMember: *infra_castmember.NewCastMemberGateway(db),
		Genre:      *infra_genre.NewGenreGateway(db),
	}
}
