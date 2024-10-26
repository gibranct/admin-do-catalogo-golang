package video

import (
	"errors"

	"github.com.br/gibranct/admin_do_catalogo/pkg/validator"
)

type VideoValidator struct {
	video    Video
	vHandler validator.ValidationHandler
}

const (
	TITLE_MAX_LENGTH       = 255
	DESCRIPTION_MAX_LENGTH = 4_000
)

func (vv VideoValidator) Validate() {
	title := vv.video.Title
	description := vv.video.Description
	if title == "" {
		vv.vHandler.Add(errors.New("'title' should not be null or empty"))
	}
	if description == "" {
		vv.vHandler.Add(errors.New("'description' should not be null or empty"))
	}
	if len(title) > TITLE_MAX_LENGTH {
		vv.vHandler.Add(errors.New("'name' must be between 1 and 255 characters"))
	}
	if len(description) > DESCRIPTION_MAX_LENGTH {
		vv.vHandler.Add(errors.New("'description' must be between 1 and 4000 characters"))
	}
}

func NewVideoValidator(v Video, handler validator.ValidationHandler) *VideoValidator {
	return &VideoValidator{
		video:    v,
		vHandler: handler,
	}
}
