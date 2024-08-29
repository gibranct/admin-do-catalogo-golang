package video

type ImageMedia struct {
	ID       int64
	Checksum string
	Name     string
	Location string
}

func NewImageMediaWithId(
	id int64, checkSum, name, location string,
) *ImageMedia {
	return &ImageMedia{
		ID:       id,
		Checksum: checkSum,
		Name:     name,
		Location: location,
	}
}

func NewImageMediaWithoutId(
	checkSum, name, location string,
) *ImageMedia {
	return &ImageMedia{
		Checksum: checkSum,
		Name:     name,
		Location: location,
	}
}
