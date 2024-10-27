package category_usecase_test

import (
	"testing"

	category_usecase "github.com.br/gibranct/admin_do_catalogo/internal/usecases/category"
	"github.com.br/gibranct/admin_do_catalogo/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryCreationUseCase(t *testing.T) {
	gatewayMock := new(mocks.CategoryGatewayMock)
	useCase := category_usecase.DefaultCreateCategoryUseCase{
		Gateway: gatewayMock,
	}
	command := category_usecase.CreateCategoryCommand{
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
	useCase := category_usecase.DefaultCreateCategoryUseCase{
		Gateway: gatewayMock,
	}
	command := category_usecase.CreateCategoryCommand{
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
