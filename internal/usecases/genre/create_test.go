package genre_usecase

import (
	"errors"
	"testing"

	"github.com.br/gibranct/admin-do-catalogo/internal/domain/genre"
	"github.com.br/gibranct/admin-do-catalogo/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenreCreationUseCase(t *testing.T) {
	gatewayMock := new(mocks.GenreGatewayMock)
	categoryGatewayMock := new(mocks.CategoryGatewayMock)
	useCase := DefaultCreateGenreUseCase{
		Gateway:         gatewayMock,
		CategoryGateway: categoryGatewayMock,
	}
	categoryIds := []int64{45, 59, 78}
	command := CreateGenreCommand{
		Name:        "Drinks",
		CategoryIds: &categoryIds,
	}
	newGenre := genre.NewGenre(command.Name)
	newGenre.AddCategoriesIds(*command.CategoryIds)

	gatewayMock.On("Create", mock.Anything).Return(nil)
	categoryGatewayMock.On("ExistsByIds", *command.CategoryIds).Return(*command.CategoryIds, nil)

	noti, cate := useCase.Execute(command)

	assert.NotNil(t, cate)
	assert.Nil(t, noti)
	gatewayMock.AssertExpectations(t)
	gatewayMock.AssertNumberOfCalls(t, "Create", 1)
	categoryGatewayMock.AssertNumberOfCalls(t, "ExistsByIds", 1)
}

func TestGenreCreationUseCaseWhenCategoryIsNotFound(t *testing.T) {
	gatewayMock := new(mocks.GenreGatewayMock)
	categoryGatewayMock := new(mocks.CategoryGatewayMock)
	useCase := DefaultCreateGenreUseCase{
		Gateway:         gatewayMock,
		CategoryGateway: categoryGatewayMock,
	}
	categoryIds := []int64{45, 59, 78}
	command := CreateGenreCommand{
		Name:        "Drinks",
		CategoryIds: &categoryIds,
	}
	categoryGatewayMock.On("ExistsByIds", *command.CategoryIds).Return([]int64{45}, nil)

	noti, cate := useCase.Execute(command)

	assert.Nil(t, cate)
	assert.NotNil(t, noti)
	assert.Len(t, noti.GetErrors(), 1)
	assert.Equal(t, noti.GetErrors()[0], errors.New("missing category ids: 59,78"))
	gatewayMock.AssertExpectations(t)
	gatewayMock.AssertNumberOfCalls(t, "Create", 0)
	categoryGatewayMock.AssertNumberOfCalls(t, "ExistsByIds", 1)
}

func TestValidateCategoriesIds(t *testing.T) {
	categoryGatewayMock := new(mocks.CategoryGatewayMock)
	categoryIds := []int64{45, 59, 78}
	useCase := DefaultCreateGenreUseCase{
		CategoryGateway: categoryGatewayMock,
	}

	categoryGatewayMock.On("ExistsByIds", categoryIds).Return(categoryIds, nil)

	err := useCase.validateCategories(categoryIds)

	assert.Nil(t, err)
}

func TestValidateCategoriesIdsDoesFindAllIds(t *testing.T) {
	categoryGatewayMock := new(mocks.CategoryGatewayMock)
	categoryIds := []int64{45, 59, 78}
	useCase := DefaultCreateGenreUseCase{
		CategoryGateway: categoryGatewayMock,
	}

	categoryGatewayMock.On("ExistsByIds", categoryIds).Return([]int64{45, 59}, nil)

	err := useCase.validateCategories(categoryIds)

	assert.NotNil(t, err)
	assert.Equal(t, "missing category ids: 78", err.Error())
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
