package category_usecase

import "github.com.br/gibranct/admin-do-catalogo/internal/domain/category"

type ActivateCategoryUseCase interface {
	Execute(categoryId int64) error
}

type DefaultActivateCategoryUseCase struct {
	Gateway category.CategoryGateway
}

func (d *DefaultActivateCategoryUseCase) Execute(categoryId int64) error {
	cate, err := d.Gateway.FindById(categoryId)
	if err != nil {
		return err
	}
	cate.Activate()

	return d.Gateway.Update(*cate)
}
