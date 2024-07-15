package category_usecase

import (
	"errors"
	"testing"

	"github.com.br/gibranct/admin-do-catalogo/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDeleteCategoryUseCase(t *testing.T) {
	gatewayMock := new(mocks.CategoryGatewayMock)
	useCase := DefaultDeleteCategoryUseCase{
		gateway: gatewayMock,
	}
	var expectedId int64 = 4545
	gatewayMock.On("DeleteById", expectedId).Return(nil)

	err := useCase.Execute(expectedId)

	assert.Nil(t, err)
	gatewayMock.AssertExpectations(t)
	gatewayMock.AssertNumberOfCalls(t, "DeleteById", 1)
}

func TestDeleteCategoryUseCaseFails(t *testing.T) {
	gatewayMock := new(mocks.CategoryGatewayMock)
	useCase := DefaultDeleteCategoryUseCase{
		gateway: gatewayMock,
	}
	var expectedId int64 = 4545
	expectedErr := errors.New("failed to delete category")
	gatewayMock.On("DeleteById", expectedId).Return(expectedErr)

	err := useCase.Execute(expectedId)

	assert.NotNil(t, err)
	assert.Equal(t, err, expectedErr)
	gatewayMock.AssertExpectations(t)
	gatewayMock.AssertNumberOfCalls(t, "DeleteById", 1)
}
