package video

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreteNewImageMediaWithId(t *testing.T) {
	id := int64(55)
	checkSum := "asdasdasd"
	name := "image_1.png"
	location := "/temp"
	imageMedia := NewImageMediaWithId(
		id, checkSum, name, location,
	)

	assert.NotNil(t, imageMedia)
	assert.Equal(t, id, imageMedia.ID)
	assert.Equal(t, checkSum, imageMedia.Checksum)
	assert.Equal(t, name, imageMedia.Name)
	assert.Equal(t, location, imageMedia.Location)
}

func TestCreteNewImageMediaWithoutId(t *testing.T) {
	checkSum := "asdasdasd"
	name := "image_1.png"
	location := "/temp"
	imageMedia := NewImageMediaWithoutId(
		checkSum, name, location,
	)

	assert.NotNil(t, imageMedia)
	assert.Equal(t, checkSum, imageMedia.Checksum)
	assert.Equal(t, name, imageMedia.Name)
	assert.Equal(t, location, imageMedia.Location)
}
