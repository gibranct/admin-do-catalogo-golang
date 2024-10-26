package castmember_usecase

import (
	"github.com.br/gibranct/admin_do_catalogo/internal/domain/castmember"
	"github.com.br/gibranct/admin_do_catalogo/pkg/notification"
)

type CreateCastMemberOutput struct {
	ID int64
}

type CreateCastMemberCommand struct {
	Name string
	Type castmember.CastMemberType
}

type CreateCastMemberUseCase interface {
	Execute(c CreateCastMemberCommand) (*notification.Notification, *CreateCastMemberOutput)
}

type DefaultCreateCastMemberUseCase struct {
	Gateway castmember.CastMemberGateway
}

func (useCase DefaultCreateCastMemberUseCase) Execute(
	command CreateCastMemberCommand,
) (*notification.Notification, *CreateCastMemberOutput) {
	castMember := castmember.NewCastMember(
		command.Name,
		command.Type,
	)

	n := notification.CreateNotification()

	castMember.Validate(n)

	if n.HasErrors() {
		return n, nil
	}

	err := useCase.Gateway.Create(castMember)

	if err != nil {
		n.Add(err)
		return n, nil
	}

	return nil, &CreateCastMemberOutput{
		ID: castMember.ID,
	}
}
