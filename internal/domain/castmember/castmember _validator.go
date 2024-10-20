package castmember

import (
	"errors"
	"fmt"
	"strings"

	"github.com.br/gibranct/admin_do_catalogo/pkg/validator"
)

const nameMaxLength = 255
const nameMinLength = 3

type CastMemberValidator struct {
	castMember CastMember
	vHandler   validator.ValidationHandler
}

func (cm CastMemberValidator) Validate() {
	name := strings.Trim(cm.castMember.Name, " ")
	if name == "" {
		cm.vHandler.Add(errors.New("'name' should not be empty"))
	}
	if len(name) < nameMinLength || len(name) > nameMaxLength {
		errorMsg := fmt.Sprintf("'name' must be between %d and %d characters", nameMinLength, nameMaxLength)
		cm.vHandler.Add(errors.New(errorMsg))
	}
}

func NewCastMemberValidator(castMember CastMember, vHandler validator.ValidationHandler) *CastMemberValidator {
	return &CastMemberValidator{
		castMember: castMember,
		vHandler:   vHandler,
	}
}
