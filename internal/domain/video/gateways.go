package video

type Resource struct {
	Content     []byte
	Checksum    string
	ContentType string
	Name        string
}

type VideoResource struct {
	Type     VideoMediaType
	Resource Resource
}

type MediaResourceGateway interface {
	StoreAudioVideo(videoId int64, resource VideoResource) (*AudioVideoMedia, error)
	StoreImage(videoId int64, resource VideoResource) (*ImageMedia, error)
	GetResource(videoId int64, aType VideoMediaType) (*Resource, error)
	ClearResources(videoId int64) error
}

type VideoGateway interface {
	Create(aVideo Video) (*Video, error)
	DeleteById(aVideo int64) error
	FindById(videoId int64) (*Video, error)
}
