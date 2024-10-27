package castmember_usecase_test

import (
	"errors"
	"testing"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain/castmember"
	castmember_usecase "github.com.br/gibranct/admin_do_catalogo/internal/usecases/castmember"
	"github.com.br/gibranct/admin_do_catalogo/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryUpdateUseCase(t *testing.T) {
	gatewayMock := new(mocks.CastMemberGatewayMock)
	useCase := castmember_usecase.DefaultUpdateCastMemberUseCase{
		Gateway: gatewayMock,
	}
	command := castmember_usecase.UpdateCastMemberCommand{
		ID:   56,
		Name: "John Doe",
		Type: "actor",
	}
	newType, _ := castmember.TypeFromString(command.Type)
	castMember := castmember.NewCastMember(command.Name, newType)
	castMember.ID = command.ID
	gatewayMock.On("FindById", command.ID).Return(castMember, nil)
	gatewayMock.On("Update", mock.Anything).Return(nil)

	noti := useCase.Execute(command)

	assert.Nil(t, noti)
	gatewayMock.AssertExpectations(t)
	gatewayMock.AssertNumberOfCalls(t, "Update", 1)
	gatewayMock.AssertNumberOfCalls(t, "FindById", 1)
}

func TestUpdateCastMemberUseCaseWhenCastMemberIsNotFound(t *testing.T) {
	gatewayMock := new(mocks.CastMemberGatewayMock)
	useCase := castmember_usecase.DefaultUpdateCastMemberUseCase{
		Gateway: gatewayMock,
	}
	command := castmember_usecase.UpdateCastMemberCommand{
		ID:   56,
		Name: "John Doe",
		Type: "actor",
	}

	newType, _ := castmember.TypeFromString(command.Type)
	castMember := castmember.NewCastMember(command.Name, newType)
	castMember.ID = command.ID
	expectedMsg := "cast member not found"
	gatewayMock.On("FindById", command.ID).Return(castMember, errors.New(""))

	noti := useCase.Execute(command)

	assert.NotNil(t, noti)
	assert.True(t, noti.HasErrors())
	assert.Equal(t, 1, len(noti.GetErrors()))
	assert.Equal(t, expectedMsg, noti.GetErrors()[0].Error())
	gatewayMock.AssertExpectations(t)
	gatewayMock.AssertNumberOfCalls(t, "FindById", 1)
}

func TestCategoryUpdateWithEmptyName(t *testing.T) {
	gatewayMock := new(mocks.CastMemberGatewayMock)
	useCase := castmember_usecase.DefaultUpdateCastMemberUseCase{
		Gateway: gatewayMock,
	}
	command := castmember_usecase.UpdateCastMemberCommand{
		ID:   56,
		Name: "",
		Type: "actor",
	}

	newType, _ := castmember.TypeFromString(command.Type)
	castMember := castmember.NewCastMember(command.Name, newType)
	castMember.ID = command.ID
	gatewayMock.On("FindById", command.ID).Return(castMember, nil)
	expectedMsg := "'name' should not be empty"

	noti := useCase.Execute(command)

	assert.NotNil(t, noti)
	assert.True(t, noti.HasErrors())
	assert.Len(t, noti.GetErrors(), 2)
	assert.Equal(t, noti.GetErrors()[0].Error(), expectedMsg)
	gatewayMock.AssertExpectations(t)
	gatewayMock.AssertNumberOfCalls(t, "Update", 0)
	gatewayMock.AssertNumberOfCalls(t, "FindById", 1)
}
