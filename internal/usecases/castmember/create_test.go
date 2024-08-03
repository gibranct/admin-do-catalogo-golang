package castmember_usecase

import (
	"testing"

	"github.com.br/gibranct/admin-do-catalogo/internal/domain/castmember"
	"github.com.br/gibranct/admin-do-catalogo/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	gatewayMock := new(mocks.CastMemberGatewayMock)
	useCase := &DefaultCreateCastMemberUseCase{
		Gateway: gatewayMock,
	}
	command := CreateCastMemberCommand{
		Name: "Keannu Reeves",
		Type: castmember.ACTOR,
	}
	gatewayMock.On("Create", mock.Anything).Return(nil)

	noti, output := useCase.Execute(command)

	assert.NotNil(t, output)
	assert.Nil(t, noti)
	gatewayMock.AssertExpectations(t)
	gatewayMock.AssertNumberOfCalls(t, "Create", 1)
}

func TestCreateWithEmptyName(t *testing.T) {
	gatewayMock := new(mocks.CastMemberGatewayMock)
	useCase := &DefaultCreateCastMemberUseCase{
		Gateway: gatewayMock,
	}
	command := CreateCastMemberCommand{
		Name: "",
		Type: castmember.ACTOR,
	}
	gatewayMock.On("Create", mock.Anything).Return(nil)
	expectedMsg := "'name' should not be empty"

	noti, output := useCase.Execute(command)

	assert.Nil(t, output)
	assert.NotNil(t, noti)
	assert.True(t, noti.HasErrors())
	assert.Len(t, noti.GetErrors(), 2)
	assert.Equal(t, noti.GetErrors()[0].Error(), expectedMsg)
	gatewayMock.AssertNumberOfCalls(t, "Create", 0)
}
