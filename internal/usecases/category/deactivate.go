package category_usecase

import "github.com.br/gibranct/admin-do-catalogo/internal/domain/category"

type DeactivateCategoryUseCase interface {
	Execute(categoryId int64) error
}

type DefaultDeactivateCategoryUseCase struct {
	Gateway category.CategoryGateway
}

func (d *DefaultDeactivateCategoryUseCase) Execute(categoryId int64) error {
	cate, err := d.Gateway.FindById(categoryId)
	if err != nil {
		return err
	}
	cate.Deactivate()

	return d.Gateway.Update(*cate)
}
