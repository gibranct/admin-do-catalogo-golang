package application_category

import (
	"errors"
	"testing"
	"time"

	"github.com.br/gibranct/admin-do-catalogo/internal/domain"
	"github.com.br/gibranct/admin-do-catalogo/internal/domain/category"
	"github.com.br/gibranct/admin-do-catalogo/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFindCategoryByIdUseCase(t *testing.T) {
	categoryId := int64(54545)
	categoryGatewayMock := &mocks.CategoryGatewayMock{}

	sut := DefaultGetCategoryByIdUseCase{
		gateway: categoryGatewayMock,
	}
	cate := category.Category{
		ID:          categoryId,
		Name:        "A",
		Description: "B",
		IsActive:    true,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		DeletedAt:   nil,
	}
	categoryGatewayMock.On("FindById", categoryId).Return(cate, nil)

	foundCate, err := sut.Execute(categoryId)

	assert.NotNil(t, foundCate)
	assert.Nil(t, err)
	assert.Equal(t, cate.ID, foundCate.ID)
	assert.Equal(t, cate.Name, foundCate.Name)
	assert.Equal(t, cate.Description, foundCate.Description)
	assert.Equal(t, cate.IsActive, foundCate.IsActive)
	assert.Equal(t, cate.CreatedAt, foundCate.CreatedAt)
	assert.Equal(t, cate.UpdatedAt, foundCate.UpdatedAt)
	assert.Equal(t, cate.DeletedAt, foundCate.DeletedAt)
	categoryGatewayMock.AssertExpectations(t)
	categoryGatewayMock.AssertNumberOfCalls(t, "FindById", 1)
}

func TestFailToFindCategoryById(t *testing.T) {
	categoryId := int64(54)
	categoryGatewayMock := &mocks.CategoryGatewayMock{}

	sut := DefaultGetCategoryByIdUseCase{
		gateway: categoryGatewayMock,
	}
	expectedErr := domain.NotFoundException{
		Message: "Could not find category for id: 54",
	}
	emptyCate := category.Category{}
	categoryGatewayMock.On("FindById", categoryId).Return(emptyCate, errors.New(""))

	foundCate, err := sut.Execute(categoryId)

	assert.Nil(t, foundCate)
	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
	categoryGatewayMock.AssertExpectations(t)
	categoryGatewayMock.AssertNumberOfCalls(t, "FindById", 1)
}
