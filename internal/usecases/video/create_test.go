package video_usecase_test

import (
	"errors"
	"strings"
	"testing"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain/video"
	video_usecase "github.com.br/gibranct/admin_do_catalogo/internal/usecases/video"
	"github.com.br/gibranct/admin_do_catalogo/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func dummyCreateVideoCommand() video_usecase.CreateVideoCommand {
	return video_usecase.CreateVideoCommand{
		Title:       "dummy title",
		Description: "dummy desc",
		LaunchedAt:  2024,
		Duration:    120.0,
		Opened:      true,
		Published:   true,
		Rating:      "Livre",
		CategoryIds: []int64{78, 45},
		GenreIds:    []int64{39},
		MemberIds:   []int64{55},
	}
}

func TestCreateVideo(t *testing.T) {
	videoGateway := new(mocks.VideoGatewayMock)
	categoryGateway := new(mocks.CategoryGatewayMock)
	genreGateway := new(mocks.GenreGatewayMock)
	castGateway := new(mocks.CastMemberGatewayMock)
	sut := video_usecase.NewDefaultCreateVideoUseCase(
		videoGateway, categoryGateway, genreGateway, castGateway,
	)
	video := video.Video{
		ID: 999,
	}

	command := dummyCreateVideoCommand()

	categoryGateway.On("ExistsByIds", command.CategoryIds).Return(
		command.CategoryIds, nil,
	)
	genreGateway.On("ExistsByIds", command.GenreIds).Return(
		command.GenreIds, nil,
	)
	castGateway.On("ExistsByIds", command.MemberIds).Return(
		command.MemberIds, nil,
	)
	videoGateway.On("Create", mock.Anything).Return(&video, nil)

	noti, output := sut.Execute(command)

	assert.NotNil(t, output)
	assert.Equal(t, video.ID, output.ID)
	assert.Nil(t, noti)

	categoryGateway.AssertExpectations(t)
	categoryGateway.AssertNumberOfCalls(t, "ExistsByIds", 1)
	genreGateway.AssertExpectations(t)
	genreGateway.AssertNumberOfCalls(t, "ExistsByIds", 1)
	castGateway.AssertExpectations(t)
	castGateway.AssertNumberOfCalls(t, "ExistsByIds", 1)
	videoGateway.AssertExpectations(t)
	videoGateway.AssertNumberOfCalls(t, "Create", 1)
}

func TestCreateVideoWithInvalidRating(t *testing.T) {
	videoGateway := new(mocks.VideoGatewayMock)
	categoryGateway := new(mocks.CategoryGatewayMock)
	genreGateway := new(mocks.GenreGatewayMock)
	castGateway := new(mocks.CastMemberGatewayMock)
	sut := video_usecase.NewDefaultCreateVideoUseCase(
		videoGateway, categoryGateway, genreGateway, castGateway,
	)
	command := dummyCreateVideoCommand()
	command.Rating = "dummy"

	noti, output := sut.Execute(command)

	assert.Nil(t, output)
	assert.NotNil(t, noti)
	assert.Len(t, noti.GetErrors(), 1)
	assert.Equal(t, "unknown type", noti.GetErrors()[0].Error())
}

func TestCreateVideoWithInvalidTitleAndDescription(t *testing.T) {
	videoGateway := new(mocks.VideoGatewayMock)
	categoryGateway := new(mocks.CategoryGatewayMock)
	genreGateway := new(mocks.GenreGatewayMock)
	castGateway := new(mocks.CastMemberGatewayMock)
	sut := video_usecase.NewDefaultCreateVideoUseCase(
		videoGateway, categoryGateway, genreGateway, castGateway,
	)

	tests := []struct {
		title string
		desc  string
		err   error
	}{
		{
			title: "",
			desc:  "dummy value",
			err:   errors.New("'title' should not be null or empty"),
		},
		{
			title: strings.Repeat("x", 256),
			desc:  "dummy desc",
			err:   errors.New("'name' must be between 1 and 255 characters"),
		},
		{
			title: "dummy value",
			desc:  "",
			err:   errors.New("'description' should not be null or empty"),
		},
		{
			title: "dummy value",
			desc:  strings.Repeat("x", 4001),
			err:   errors.New("'description' must be between 1 and 4000 characters"),
		},
	}

	for _, test := range tests {
		command := dummyCreateVideoCommand()
		command.Title = test.title
		command.Description = test.desc

		noti, output := sut.Execute(command)

		assert.Nil(t, output)
		assert.NotNil(t, noti)
		assert.Len(t, noti.GetErrors(), 1)
		assert.Equal(t, test.err.Error(), noti.GetErrors()[0].Error())
	}
}
