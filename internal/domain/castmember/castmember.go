package castmember

import (
	"errors"
	"time"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain"
	"github.com.br/gibranct/admin_do_catalogo/pkg/validator"
)

type CastMember struct {
	ID        int64
	Name      string
	Type      CastMemberType
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CastMemberGateway interface {
	Create(castMember *CastMember) error
	FindById(castMemberId int64) (*CastMember, error)
	DeleteById(castMemberId int64) error
	Update(castMember CastMember) error
	FindAll(query domain.SearchQuery) (*domain.Pagination[CastMember], error)
	ExistsByIds(castMemberIds []int64) ([]int64, error)
}

func NewCastMember(
	name string,
	castMemberType CastMemberType,
) *CastMember {
	now := time.Now().UTC()
	return &CastMember{
		Name:      name,
		Type:      castMemberType,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (cm *CastMember) Update(name string, aType CastMemberType) {
	cm.Name = name
	cm.Type = aType
	cm.UpdatedAt = time.Now().UTC()
}

func (cm *CastMember) Validate(handler validator.ValidationHandler) {
	NewCastMemberValidator(*cm, handler).Validate()
}

func (cm *CastMember) ChangeType(typeStr string) error {
	switch typeStr {
	case ACTOR.String():
		cm.Type = ACTOR
	case DIRECTOR.String():
		cm.Type = DIRECTOR
	default:
		return errors.New("unknown type")
	}
	return nil
}
