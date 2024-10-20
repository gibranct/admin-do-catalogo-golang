package category_usecase

import (
	"testing"

	"github.com.br/gibranct/admin_do_catalogo/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryCreationUseCase(t *testing.T) {
	gatewayMock := new(mocks.CategoryGatewayMock)
	useCase := DefaultCreateCategoryUseCase{
		Gateway: gatewayMock,
	}
	command := CreateCategoryCommand{
		Name:        "Drinks",
		Description: "All cool drinks",
	}
	gatewayMock.On("Create", mock.Anything).Return(nil)

	noti, cate := useCase.Execute(command)

	assert.NotNil(t, cate)
	assert.Nil(t, noti)
	gatewayMock.AssertExpectations(t)
	gatewayMock.AssertNumberOfCalls(t, "Create", 1)
}

func TestCategoryCreationWithEmptyName(t *testing.T) {
	gatewayMock := new(mocks.CategoryGatewayMock)
	useCase := DefaultCreateCategoryUseCase{
		Gateway: gatewayMock,
	}
	command := CreateCategoryCommand{
		Name:        "",
		Description: "All cool drinks",
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
