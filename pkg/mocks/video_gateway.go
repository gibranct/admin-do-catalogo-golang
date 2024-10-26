package mocks

import (
	"github.com.br/gibranct/admin_do_catalogo/internal/domain/video"
	"github.com/stretchr/testify/mock"
)

type VideoGatewayMock struct {
	mock.Mock
}

func (vg *VideoGatewayMock) Create(aVideo video.Video) (*video.Video, error) {
	args := vg.Called(aVideo)
	return args.Get(0).(*video.Video), args.Error(1)
}

func (vg *VideoGatewayMock) DeleteById(aVideo int64) error {
	args := vg.Called(aVideo)
	return args.Error(0)
}

func (vg *VideoGatewayMock) FindById(videoId int64) (*video.Video, error) {
	args := vg.Called(videoId)
	return args.Get(0).(*video.Video), args.Error(1)
}
