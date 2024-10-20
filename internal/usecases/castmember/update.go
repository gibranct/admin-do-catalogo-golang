package castmember_usecase

import (
	"errors"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain/castmember"
	"github.com.br/gibranct/admin_do_catalogo/pkg/notification"
)

type UpdateCastMemberCommand struct {
	ID   int64
	Name string
	Type string
}

type UpdateCategoryUseCase interface {
	Execute(c UpdateCastMemberCommand) *notification.Notification
}

type DefaultUpdateCastMemberUseCase struct {
	Gateway castmember.CastMemberGateway
}

func (useCase DefaultUpdateCastMemberUseCase) Execute(
	command UpdateCastMemberCommand,
) *notification.Notification {

	castMember, err := useCase.Gateway.FindById(command.ID)

	n := notification.CreateNotification()

	if err != nil {
		n.Add(errors.New("cast member not found"))
		return n
	}

	newType, err := castmember.TypeFromString(command.Type)

	if err != nil {
		n.Add(err)
		return n
	}

	castMember.Update(command.Name, newType)

	castMember.Validate(n)

	if n.HasErrors() {
		return n
	}

	err = useCase.Gateway.Update(*castMember)

	if err != nil {
		n.Add(err)
		return n
	}

	return nil
}
