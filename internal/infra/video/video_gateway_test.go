package infra_video_test

import (
	"errors"
	"log"
	"testing"

	"github.com.br/gibranct/admin-do-catalogo/internal/domain/video"
	infra_video "github.com.br/gibranct/admin-do-catalogo/internal/infra/video"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateVideo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()

	vg := infra_video.NewVideoGateway(db)
	video := dummyVideo()
	videoId := int64(85)
	categoryId := video.CategoryIds[0]
	genreId := video.GenreIds[0]
	castMemberId := video.CastMemberIds[0]

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO VIDEOS").WithArgs(
		video.Title,
		video.Description,
		video.LaunchedAt,
		video.Opened,
		video.Published,
		video.Rating.String(),
		video.Duration,
		video.CreatedAt,
		video.UpdatedAt,
		nil,
		nil,
		nil,
		nil,
		nil,
	).WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(videoId))

	mock.ExpectExec("INSERT INTO videos_categories").WithArgs(
		videoId, categoryId,
	).WillReturnResult(sqlmock.NewResult(categoryId, 1))

	mock.ExpectExec("INSERT INTO videos_cast_members").WithArgs(
		videoId, castMemberId,
	).WillReturnResult(sqlmock.NewResult(castMemberId, 1))

	mock.ExpectExec("INSERT INTO videos_genres").WithArgs(
		videoId, genreId,
	).WillReturnResult(sqlmock.NewResult(genreId, 1))

	mock.ExpectCommit()

	savedVideo, err := vg.Create(video)

	assert.Nil(t, err)
	assert.Equal(t, videoId, savedVideo.ID)
}

func TestCreateVideoWhenItFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()

	vg := infra_video.NewVideoGateway(db)
	video := dummyVideo()
	videoId := int64(85)
	createVideoError := errors.New("failed to save category")

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO VIDEOS").WithArgs(
		video.Title,
		video.Description,
		video.LaunchedAt,
		video.Opened,
		video.Published,
		video.Rating.String(),
		video.Duration,
		video.CreatedAt,
		video.UpdatedAt,
		nil,
		nil,
		nil,
		nil,
		nil,
	).WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(videoId))
	mock.ExpectExec("INSERT INTO videos_categories").WithArgs(
		videoId, video.CategoryIds[0],
	).WillReturnError(createVideoError)

	mock.ExpectRollback()

	savedVideo, err := vg.Create(video)

	assert.Nil(t, savedVideo)
	assert.NotNil(t, err)
	assert.Equal(t, createVideoError.Error(), err.Error())
}

func dummyVideo() video.Video {
	return *video.NewVideo(
		"dummy title",
		"dummy desc",
		2024,
		120.0,
		true,
		true,
		video.L,
		[]int64{78},
		[]int64{39},
		[]int64{55},
	)
}
