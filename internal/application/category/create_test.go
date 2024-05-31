package category

import (
	"testing"

	"github.com.br/gibranct/admin-do-catalogo/internal/domain/category"
	"github.com.br/gibranct/admin-do-catalogo/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryCreationUseCase(t *testing.T) {
	gatewayMock := new(mocks.CategoryGatewayMock)
	useCase := DefaultCreateCategoryUseCase{
		gateway: gatewayMock,
	}
	command := CreateCategoryCommand{
		Name:        "Drinks",
		Description: "All cool drinks",
		IsActive:    true,
	}
	var expectedId int64 = 4545
	newCategory := category.Category{
		ID: expectedId,
	}
	gatewayMock.On("Create", mock.Anything).Return(newCategory, nil)

	noti, cate := useCase.Execute(command)

	assert.NotNil(t, cate)
	assert.Nil(t, noti)
	assert.Equal(t, cate.ID, expectedId)
	gatewayMock.AssertExpectations(t)
	gatewayMock.AssertNumberOfCalls(t, "Create", 1)
}

func TestCategoryCreationWithEmptyName(t *testing.T) {
	gatewayMock := new(mocks.CategoryGatewayMock)
	useCase := DefaultCreateCategoryUseCase{
		gateway: gatewayMock,
	}
	command := CreateCategoryCommand{
		Name:        "",
		Description: "All cool drinks",
		IsActive:    true,
	}
	expectedMsg := "'name' should not be empty"
	newCategory := category.Category{
		ID: 552,
	}
	gatewayMock.On("Create", mock.Anything).Return(newCategory, nil)

	noti, cate := useCase.Execute(command)

	assert.Nil(t, cate)
	assert.NotNil(t, noti)
	assert.True(t, noti.HasErrors())
	assert.Len(t, noti.GetErrors(), 2)
	assert.Equal(t, noti.GetErrors()[0].Error(), expectedMsg)
	gatewayMock.AssertNumberOfCalls(t, "Create", 0)
}
