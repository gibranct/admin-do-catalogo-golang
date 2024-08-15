package genre_usecase

import (
	"testing"

	"github.com.br/gibranct/admin-do-catalogo/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenreCreationUseCase(t *testing.T) {
	gatewayMock := new(mocks.GenreGatewayMock)
	useCase := DefaultCreateGenreUseCase{
		Gateway: gatewayMock,
	}
	command := CreateGenreCommand{
		Name: "Drinks",
	}
	gatewayMock.On("Create", mock.Anything).Return(nil)

	noti, cate := useCase.Execute(command)

	assert.NotNil(t, cate)
	assert.Nil(t, noti)
	gatewayMock.AssertExpectations(t)
	gatewayMock.AssertNumberOfCalls(t, "Create", 1)
}

func TestGenreCreationWithEmptyName(t *testing.T) {
	gatewayMock := new(mocks.GenreGatewayMock)
	useCase := DefaultCreateGenreUseCase{
		Gateway: gatewayMock,
	}
	command := CreateGenreCommand{
		Name: "",
	}
	expectedMsg := "'name' should not be empty"
	gatewayMock.On("Create", mock.Anything).Return(nil)

	noti, cate := useCase.Execute(command)

	assert.Nil(t, cate)
	assert.NotNil(t, noti)
	assert.True(t, noti.HasErrors())
	assert.Len(t, noti.GetErrors(), 2)
	assert.Equal(t, noti.GetErrors()[0].Error(), expectedMsg)
	gatewayMock.AssertNumberOfCalls(t, "Create", 0)
}
