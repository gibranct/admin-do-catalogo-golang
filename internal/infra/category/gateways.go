package gateway

import "database/sql"

type Gateways struct {
	Category CategoryGateway
}

func NewGateways(db *sql.DB) Gateways {
	return Gateways{
		Category: *NewCategoryGateway(db),
	}
}
