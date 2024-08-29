package video

type AudioVideoMedia struct {
	ID              int64
	Checksum        string
	Name            string
	RawLocation     string
	EncodedLocation string
	Status          *MediaStatus
}

func NewAudioVideoMediaWith(
	id int64,
	status *MediaStatus,
	checksum,
	name,
	rawLocation,
	encodedLocation string,
) *AudioVideoMedia {
	return &AudioVideoMedia{
		ID:              id,
		Checksum:        checksum,
		Name:            name,
		RawLocation:     rawLocation,
		EncodedLocation: encodedLocation,
		Status:          status,
	}
}

func (avm *AudioVideoMedia) processing() *AudioVideoMedia {
	status := PROCESSING
	return NewAudioVideoMediaWith(
		avm.ID,
		&status,
		avm.Checksum,
		avm.Name,
		avm.RawLocation,
		avm.EncodedLocation,
	)
}

func (avm *AudioVideoMedia) completed(encodedPath string) *AudioVideoMedia {
	status := COMPLETED
	return NewAudioVideoMediaWith(
		avm.ID,
		&status,
		avm.Checksum,
		avm.Name,
		avm.RawLocation,
		encodedPath,
	)
}

func (avm *AudioVideoMedia) IsPendingEncode() bool {
	return PENDING == *avm.Status
}
