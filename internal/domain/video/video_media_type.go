package video

import (
	"errors"
	"fmt"
)

type VideoMediaType uint8

const (
	VIDEO VideoMediaType = iota
	TRAILER
	BANNER
	THUMBNAIL
	THUMBNAIL_HALF
)

func (v VideoMediaType) String() string {
	switch v {
	case VIDEO:
		return "Video"
	case TRAILER:
		return "Trailer"
	case BANNER:
		return "Banner"
	case THUMBNAIL:
		return "Thumbnail"
	case THUMBNAIL_HALF:
		return "Thumbnail_half"

	}
	return "unknown"
}

func GetVideoType(value string) (VideoMediaType, error) {
	switch value {
	case "Video":
		return VIDEO, nil
	case "Trailer":
		return TRAILER, nil
	case "Banner":
		return BANNER, nil
	case "Thumbnail":
		return THUMBNAIL, nil
	case "Thumbnail_half":
		return THUMBNAIL_HALF, nil
	}

	return 0, errors.New(fmt.Sprintf("unknown video type: %s", value))
}
