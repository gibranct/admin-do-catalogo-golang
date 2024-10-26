package gateway

import (
	"database/sql"

	infra_castmember "github.com.br/gibranct/admin_do_catalogo/internal/infra/castmember"
	infra_category "github.com.br/gibranct/admin_do_catalogo/internal/infra/category"
	infra_genre "github.com.br/gibranct/admin_do_catalogo/internal/infra/genre"
	infra_video "github.com.br/gibranct/admin_do_catalogo/internal/infra/video"
)

type Gateways struct {
	Category   infra_category.CategoryGateway
	CastMember infra_castmember.CastMemberGateway
	Genre      infra_genre.GenreGateway
	Video      infra_video.VideoGateway
}

func NewGateways(db *sql.DB) Gateways {
	return Gateways{
		Category:   *infra_category.NewCategoryGateway(db),
		CastMember: *infra_castmember.NewCastMemberGateway(db),
		Genre:      *infra_genre.NewGenreGateway(db),
		Video:      *infra_video.NewVideoGateway(db),
	}
}
