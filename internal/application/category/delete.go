package application_category

import "github.com.br/gibranct/admin-do-catalogo/internal/domain/category"

type DeleteCategoryUseCase interface {
	Execute(categoryId int64) error
}

type DefaultDeleteCategoryUseCase struct {
	gateway category.CategoryGateway
}

func (d *DefaultDeleteCategoryUseCase) Execute(categoryId int64) error {
	return d.gateway.DeleteById(categoryId)
}
