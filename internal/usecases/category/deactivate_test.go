package category_usecase

import (
	"errors"
	"testing"
	"time"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain/category"
	"github.com.br/gibranct/admin_do_catalogo/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeactivateCategoryUseCase(t *testing.T) {
	gatewayMock := new(mocks.CategoryGatewayMock)
	useCase := DefaultDeactivateCategoryUseCase{
		Gateway: gatewayMock,
	}
	var expectedId int64 = 4545
	cate := &category.Category{
		ID:          expectedId,
		Name:        "fake",
		Description: "fake desc",
		IsActive:    true,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		DeletedAt:   nil,
	}
	gatewayMock.On("FindById", expectedId).Return(cate, nil)
	gatewayMock.On("Update", mock.Anything).Return(nil)

	err := useCase.Execute(expectedId)

	assert.Nil(t, err)
	assert.False(t, cate.IsActive)
	gatewayMock.AssertExpectations(t)
	gatewayMock.AssertNumberOfCalls(t, "FindById", 1)
	gatewayMock.AssertNumberOfCalls(t, "Update", 1)
}

func TestDeactivateCategoryUseCaseFails(t *testing.T) {
	gatewayMock := new(mocks.CategoryGatewayMock)
	useCase := DefaultDeactivateCategoryUseCase{
		Gateway: gatewayMock,
	}
	var expectedId int64 = 4545
	expectedErr := errors.New("failed to find category")
	gatewayMock.On("FindById", expectedId).Return(&category.Category{}, expectedErr)

	err := useCase.Execute(expectedId)

	assert.NotNil(t, err)
	assert.Equal(t, err, expectedErr)
	gatewayMock.AssertExpectations(t)
	gatewayMock.AssertNumberOfCalls(t, "FindById", 1)
}
