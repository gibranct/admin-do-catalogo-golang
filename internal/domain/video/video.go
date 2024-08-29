package video

import (
	"time"

	"github.com.br/gibranct/admin-do-catalogo/pkg/validator"
)

type Video struct {
	ID            int64
	Title         string
	Description   string
	LaunchedAt    int
	Duration      float64
	Rating        Rating
	Opened        bool
	Published     bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Banner        *ImageMedia
	ThumbNail     *ImageMedia
	ThumbNailHalf *ImageMedia
	Video         *AudioVideoMedia
	Trailer       *AudioVideoMedia
	GenreIds      []int64
	CategoryIds   []int64
	CastMemberIds []int64
}

func (v *Video) Validate(handler validator.ValidationHandler) {
	NewVideoValidator(*v, handler).Validate()
}

func NewVideo(
	title string,
	description string,
	launchedAt int,
	duration float64,
	opened bool,
	published bool,
	rating Rating,
	categoryIds []int64,
	genreIds []int64,
	castMemberIds []int64,
) *Video {
	now := time.Now().UTC()
	return &Video{
		Title:         title,
		Description:   description,
		LaunchedAt:    launchedAt,
		Duration:      duration,
		Opened:        opened,
		Published:     published,
		Rating:        rating,
		CategoryIds:   categoryIds,
		GenreIds:      genreIds,
		CastMemberIds: castMemberIds,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (v *Video) UpdateBannerMedia(banner *ImageMedia) *Video {
	v.Banner = banner
	v.UpdatedAt = time.Now().UTC()
	return v
}

func (v *Video) UpdateThumbnailMedia(thumbnail *ImageMedia) *Video {
	v.ThumbNail = thumbnail
	v.UpdatedAt = time.Now().UTC()
	return v
}

func (v *Video) UpdateThumbnailHalfMedia(thumbnailHalf *ImageMedia) *Video {
	v.ThumbNailHalf = thumbnailHalf
	v.UpdatedAt = time.Now().UTC()
	return v
}

func (v *Video) UpdateTrailerMedia(trailer *AudioVideoMedia) *Video {
	v.Trailer = trailer
	v.UpdatedAt = time.Now().UTC()
	return v
}

func (v *Video) UpdateVideoMedia(video *AudioVideoMedia) *Video {
	v.Video = video
	v.UpdatedAt = time.Now().UTC()
	return v
}

func (v *Video) Processing(aType VideoMediaType) *Video {
	if VIDEO == aType && v.Video != nil {
		v.UpdateVideoMedia(v.Video.processing())
	} else if TRAILER == aType && v.Trailer != nil {
		v.UpdateTrailerMedia(v.Trailer.processing())
	}
	return v
}

func (v *Video) Completed(aType VideoMediaType, encodedPath string) *Video {
	if VIDEO == aType {
		v.UpdateVideoMedia(v.Video.completed(encodedPath))
	} else if TRAILER == aType {
		v.UpdateTrailerMedia(v.Trailer.completed(encodedPath))
	}
	return v
}
