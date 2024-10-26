package category_usecase

import (
	"errors"
	"testing"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain/category"
	"github.com.br/gibranct/admin_do_catalogo/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryUpdateUseCase(t *testing.T) {
	gatewayMock := new(mocks.CategoryGatewayMock)
	useCase := DefaultUpdateCategoryUseCase{
		Gateway: gatewayMock,
	}
	command := UpdateCategoryCommand{
		ID:          56,
		Name:        "Drinks",
		Description: "All cool drinks",
	}
	category := category.NewCategory(command.Name, command.Description)
	category.ID = command.ID
	gatewayMock.On("FindById", command.ID).Return(category, nil)
	gatewayMock.On("Update", mock.Anything).Return(nil)

	noti := useCase.Execute(command)

	assert.Nil(t, noti)
	gatewayMock.AssertExpectations(t)
	gatewayMock.AssertNumberOfCalls(t, "Update", 1)
	gatewayMock.AssertNumberOfCalls(t, "FindById", 1)
}

func TestCategoryUpdateUseCaseWhenCategoryIsNotFound(t *testing.T) {
	gatewayMock := new(mocks.CategoryGatewayMock)
	useCase := DefaultUpdateCategoryUseCase{
		Gateway: gatewayMock,
	}
	command := UpdateCategoryCommand{
		ID:          56,
		Name:        "Drinks",
		Description: "All cool drinks",
	}
	category := category.NewCategory(command.Name, command.Description)
	category.ID = command.ID
	expectedMsg := "category not found"
	gatewayMock.On("FindById", command.ID).Return(category, errors.New(""))

	noti := useCase.Execute(command)

	assert.NotNil(t, noti)
	assert.True(t, noti.HasErrors())
	assert.Equal(t, 1, len(noti.GetErrors()))
	assert.Equal(t, expectedMsg, noti.GetErrors()[0].Error())
	gatewayMock.AssertExpectations(t)
	gatewayMock.AssertNumberOfCalls(t, "FindById", 1)
}

func TestCategoryUpdateWithEmptyName(t *testing.T) {
	gatewayMock := new(mocks.CategoryGatewayMock)
	useCase := DefaultUpdateCategoryUseCase{
		Gateway: gatewayMock,
	}
	command := UpdateCategoryCommand{
		ID:          56,
		Name:        "",
		Description: "All cool drinks",
	}
	category := category.NewCategory(command.Name, command.Description)
	category.ID = command.ID
	gatewayMock.On("FindById", command.ID).Return(category, nil)
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
