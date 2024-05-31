package category

import (
	"errors"
	"fmt"
	"strings"

	"github.com.br/gibranct/admin-do-catalogo/pkg/validator"
)

const nameMaxLength = 255
const nameMinLength = 3

type CategoryValidator struct {
	category Category
	vHandler validator.ValidationHandler
}

func (cv CategoryValidator) Validate() {
	name := strings.Trim(cv.category.Name, " ")
	if name == "" {
		cv.vHandler.Add(errors.New("'name' should not be empty"))
	}
	if len(name) < nameMinLength || len(name) > nameMaxLength {
		errorMsg := fmt.Sprintf("'name' must be between %d and %d characters", nameMinLength, nameMaxLength)
		cv.vHandler.Add(errors.New(errorMsg))
	}
}

func NewCategoryValidator(category Category, vHandler validator.ValidationHandler) *CategoryValidator {
	return &CategoryValidator{
		category: category,
		vHandler: vHandler,
	}
}
