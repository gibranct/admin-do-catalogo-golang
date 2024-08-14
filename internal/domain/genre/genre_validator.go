package genre

import (
	"errors"
	"fmt"
	"strings"

	"github.com.br/gibranct/admin-do-catalogo/pkg/validator"
)

const nameMaxLength = 255
const nameMinLength = 3

type GenreValidator struct {
	genre    Genre
	vHandler validator.ValidationHandler
}

func (cv GenreValidator) Validate() {
	name := strings.Trim(cv.genre.Name, " ")
	if name == "" {
		cv.vHandler.Add(errors.New("'name' should not be empty"))
	}
	if len(name) < nameMinLength || len(name) > nameMaxLength {
		errorMsg := fmt.Sprintf("'name' must be between %d and %d characters", nameMinLength, nameMaxLength)
		cv.vHandler.Add(errors.New(errorMsg))
	}
}

func NewGenreValidator(genre Genre, vHandler validator.ValidationHandler) *GenreValidator {
	return &GenreValidator{
		genre:    genre,
		vHandler: vHandler,
	}
}
