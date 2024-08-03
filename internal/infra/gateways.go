package gateway

import (
	"database/sql"

	infra_castmember "github.com.br/gibranct/admin-do-catalogo/internal/infra/castmember"
	infra_category "github.com.br/gibranct/admin-do-catalogo/internal/infra/category"
)

type Gateways struct {
	Category   infra_category.CategoryGateway
	CastMember infra_castmember.CastMemberGateway
}

func NewGateways(db *sql.DB) Gateways {
	return Gateways{
		Category:   *infra_category.NewCategoryGateway(db),
		CastMember: *infra_castmember.NewCastMemberGateway(db),
	}
}
