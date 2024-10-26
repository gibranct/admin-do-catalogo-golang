package castmember_usecase_test

import (
	"testing"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain/castmember"
	castmember_usecase "github.com.br/gibranct/admin_do_catalogo/internal/usecases/castmember"
	"github.com.br/gibranct/admin_do_catalogo/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	gatewayMock := new(mocks.CastMemberGatewayMock)
	useCase := &castmember_usecase.DefaultCreateCastMemberUseCase{
		Gateway: gatewayMock,
	}
	command := castmember_usecase.CreateCastMemberCommand{
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
	useCase := &castmember_usecase.DefaultCreateCastMemberUseCase{
		Gateway: gatewayMock,
	}
	command := castmember_usecase.CreateCastMemberCommand{
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
