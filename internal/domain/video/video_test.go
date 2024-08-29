package video

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateBannerMedia(t *testing.T) {
	title := "title"
	description := "desc"
	launchedAt := 2025
	duration := 54.4
	opened := true
	published := true
	ids := []int64{12, 57}

	video := NewVideo(
		title, description, launchedAt, duration, opened, published, L, ids, ids, ids,
	)
	assert.Nil(t, video.Banner)
	imageMedia := NewImageMediaWithId(
		int64(565),
		"checksum",
		"name",
		"/asdasd/x.file",
	)

	updatedTime := video.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	video.UpdateBannerMedia(imageMedia)

	assert.NotNil(t, video.Banner)
	assert.Equal(t, video.Banner, imageMedia)
	assert.True(t, updatedTime.Before(video.UpdatedAt))
}

func TestUpdateThumbnailMedia(t *testing.T) {
	title := "title"
	description := "desc"
	launchedAt := 2025
	duration := 54.4
	opened := true
	published := true
	ids := []int64{12, 57}

	video := NewVideo(
		title, description, launchedAt, duration, opened, published, L, ids, ids, ids,
	)
	assert.Nil(t, video.ThumbNail)
	imageMedia := NewImageMediaWithId(
		int64(565),
		"checksum",
		"name",
		"/asdasd/x.file",
	)

	updatedTime := video.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	video.UpdateThumbnailMedia(imageMedia)

	assert.NotNil(t, video.ThumbNail)
	assert.Equal(t, video.ThumbNail, imageMedia)
	assert.True(t, updatedTime.Before(video.UpdatedAt))
}

func TestUpdateThumbnailHalfMedia(t *testing.T) {
	title := "title"
	description := "desc"
	launchedAt := 2025
	duration := 54.4
	opened := true
	published := true
	ids := []int64{12, 57}

	video := NewVideo(
		title, description, launchedAt, duration, opened, published, L, ids, ids, ids,
	)
	assert.Nil(t, video.ThumbNailHalf)
	imageMedia := NewImageMediaWithId(
		int64(565),
		"checksum",
		"name",
		"/asdasd/x.file",
	)

	updatedTime := video.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	video.UpdateThumbnailHalfMedia(imageMedia)

	assert.NotNil(t, video.ThumbNailHalf)
	assert.Equal(t, video.ThumbNailHalf, imageMedia)
	assert.True(t, updatedTime.Before(video.UpdatedAt))
}

func TestUpdateTrailerMedia(t *testing.T) {
	title := "title"
	description := "desc"
	launchedAt := 2025
	duration := 54.4
	opened := true
	published := true
	ids := []int64{12, 57}

	video := NewVideo(
		title, description, launchedAt, duration, opened, published, L, ids, ids, ids,
	)
	assert.Nil(t, video.Trailer)
	status := PENDING
	audioVideoMedia := NewAudioVideoMediaWith(
		int64(565),
		&status,
		"checksum",
		"name",
		"/asdasd/x.file",
		"encoded-location",
	)

	updatedTime := video.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	video.UpdateTrailerMedia(audioVideoMedia)

	assert.NotNil(t, video.Trailer)
	assert.Equal(t, video.Trailer, audioVideoMedia)
	assert.True(t, updatedTime.Before(video.UpdatedAt))
}

func TestUpdateVideoMedia(t *testing.T) {
	title := "title"
	description := "desc"
	launchedAt := 2025
	duration := 54.4
	opened := true
	published := true
	ids := []int64{12, 57}

	video := NewVideo(
		title, description, launchedAt, duration, opened, published, L, ids, ids, ids,
	)
	assert.Nil(t, video.Video)
	status := PENDING
	audioVideoMedia := NewAudioVideoMediaWith(
		int64(565),
		&status,
		"checksum",
		"name",
		"/asdasd/x.file",
		"encoded-location",
	)

	updatedTime := video.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	video.UpdateVideoMedia(audioVideoMedia)

	assert.NotNil(t, video.Video)
	assert.Equal(t, video.Video, audioVideoMedia)
	assert.True(t, updatedTime.Before(video.UpdatedAt))
}

func TestProcessingVideo(t *testing.T) {
	title := "title"
	description := "desc"
	launchedAt := 2025
	duration := 54.4
	opened := true
	published := true
	ids := []int64{12, 57}

	video := NewVideo(
		title, description, launchedAt, duration, opened, published, L, ids, ids, ids,
	)
	updatedTime := video.UpdatedAt
	aType := VIDEO
	status := PENDING
	audioVideoMedia := NewAudioVideoMediaWith(
		int64(565),
		&status,
		"checksum",
		"name",
		"/asdasd/x.file",
		"encoded-location",
	)

	video.Video = audioVideoMedia
	assert.Equal(t, status, *video.Video.Status)

	video.Processing(aType)

	assert.Equal(t, PROCESSING, *video.Video.Status)
	assert.NotNil(t, video.Video)
	assert.True(t, updatedTime.Before(video.UpdatedAt))
}

func TestCompletedTrailer(t *testing.T) {
	title := "title"
	description := "desc"
	launchedAt := 2025
	duration := 54.4
	opened := true
	published := true
	ids := []int64{12, 57}
	expectedEncodedPath := "/new-path/encoded"

	video := NewVideo(
		title, description, launchedAt, duration, opened, published, L, ids, ids, ids,
	)
	updatedTime := video.UpdatedAt
	aType := TRAILER
	status := PENDING
	audioVideoMedia := NewAudioVideoMediaWith(
		int64(565),
		&status,
		"checksum",
		"name",
		"/asdasd/x.file",
		"encoded-location",
	)

	video.Trailer = audioVideoMedia
	assert.Equal(t, status, *video.Trailer.Status)

	video.Completed(aType, expectedEncodedPath)

	assert.Equal(t, COMPLETED, *video.Trailer.Status)
	assert.Equal(t, expectedEncodedPath, video.Trailer.EncodedLocation)
	assert.True(t, updatedTime.Before(video.UpdatedAt))
}

func TestProcessingTrailer(t *testing.T) {
	title := "title"
	description := "desc"
	launchedAt := 2025
	duration := 54.4
	opened := true
	published := true
	ids := []int64{12, 57}

	video := NewVideo(
		title, description, launchedAt, duration, opened, published, L, ids, ids, ids,
	)
	updatedTime := video.UpdatedAt
	aType := TRAILER
	status := PENDING
	audioVideoMedia := NewAudioVideoMediaWith(
		int64(565),
		&status,
		"checksum",
		"name",
		"/asdasd/x.file",
		"encoded-location",
	)

	video.Trailer = audioVideoMedia
	assert.Equal(t, status, *video.Trailer.Status)

	video.Processing(aType)

	assert.Equal(t, PROCESSING, *video.Trailer.Status)
	assert.NotNil(t, video.Trailer)
	assert.True(t, updatedTime.Before(video.UpdatedAt))
}

func TestCompletedVideo(t *testing.T) {
	title := "title"
	description := "desc"
	launchedAt := 2025
	duration := 54.4
	opened := true
	published := true
	ids := []int64{12, 57}
	expectedEncodedPath := "/new-path/encoded"

	video := NewVideo(
		title, description, launchedAt, duration, opened, published, L, ids, ids, ids,
	)
	updatedTime := video.UpdatedAt
	aType := VIDEO
	status := PENDING
	audioVideoMedia := NewAudioVideoMediaWith(
		int64(565),
		&status,
		"checksum",
		"name",
		"/asdasd/x.file",
		"encoded-location",
	)

	video.Video = audioVideoMedia
	assert.Equal(t, status, *video.Video.Status)

	video.Completed(aType, expectedEncodedPath)

	assert.Equal(t, COMPLETED, *video.Video.Status)
	assert.Equal(t, expectedEncodedPath, video.Video.EncodedLocation)
	assert.True(t, updatedTime.Before(video.UpdatedAt))
}
