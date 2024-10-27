package video_usecase

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain/castmember"
	"github.com.br/gibranct/admin_do_catalogo/internal/domain/category"
	"github.com.br/gibranct/admin_do_catalogo/internal/domain/genre"
	"github.com.br/gibranct/admin_do_catalogo/internal/domain/video"
	"github.com.br/gibranct/admin_do_catalogo/pkg/notification"
)

type CreateVideoOutput struct {
	ID int64
}

type CreateVideoCommand struct {
	Title         string
	Description   string
	LaunchedAt    int
	Duration      float64
	Opened        bool
	Published     bool
	Rating        string
	CategoryIds   []int64
	GenreIds      []int64
	MemberIds     []int64
	Video         *video.Resource
	Trailer       *video.Resource
	Banner        *video.Resource
	Thumbnail     *video.Resource
	ThumbnailHalf *video.Resource
}

type CreateVideoUseCase interface {
	Execute(c CreateVideoCommand) (*notification.Notification, *CreateVideoOutput)
}

type DefaultCreateVideoUseCase struct {
	Gateway           video.VideoGateway
	CategoryGateway   category.CategoryGateway
	GenreGateway      genre.GenreGateway
	CastMemberGateway castmember.CastMemberGateway
}

func NewDefaultCreateVideoUseCase(
	vg video.VideoGateway,
	cg category.CategoryGateway,
	gg genre.GenreGateway,
	ccg castmember.CastMemberGateway,
) *DefaultCreateVideoUseCase {
	return &DefaultCreateVideoUseCase{
		Gateway:           vg,
		CategoryGateway:   cg,
		GenreGateway:      gg,
		CastMemberGateway: ccg,
	}
}

func (useCase DefaultCreateVideoUseCase) Execute(
	command CreateVideoCommand,
) (*notification.Notification, *CreateVideoOutput) {

	n := notification.CreateNotification()

	rating, err := video.StringToRating(command.Rating)

	if err != nil {
		n.Add(err)
		return n, nil
	}

	video := video.NewVideo(
		command.Title,
		command.Description,
		command.LaunchedAt,
		command.Duration,
		command.Opened,
		command.Published,
		rating,
		command.CategoryIds,
		command.GenreIds,
		command.MemberIds,
	)

	video.Validate(n)

	if n.HasErrors() {
		return n, nil
	}

	n.Append(useCase.ValidateCategories(command.CategoryIds))
	n.Append(useCase.validateGenres(command.GenreIds))
	n.Append(useCase.validateMembers(command.MemberIds))

	if n.HasErrors() {
		return n, nil
	}

	savedVideo, err := useCase.Gateway.Create(*video)

	if err != nil {
		n.Add(err)
		return n, nil
	}

	return nil, &CreateVideoOutput{
		ID: savedVideo.ID,
	}
}

func (useCase DefaultCreateVideoUseCase) ValidateCategories(ids []int64) *notification.Notification {
	return validateAggregate("categories", ids, useCase.CategoryGateway.ExistsByIds)
}

func (useCase DefaultCreateVideoUseCase) validateGenres(ids []int64) *notification.Notification {
	return validateAggregate("genres", ids, useCase.GenreGateway.ExistsByIds)
}

func (useCase DefaultCreateVideoUseCase) validateMembers(ids []int64) *notification.Notification {
	return validateAggregate("cast members", ids, useCase.CastMemberGateway.ExistsByIds)
}

func validateAggregate(aggregate string, ids []int64, fn func(ids []int64) ([]int64, error)) *notification.Notification {
	n := notification.CreateNotification()
	if len(ids) == 0 {
		return n
	}

	retrievedIds, err := fn(ids)

	if err != nil {
		n.Add(err)
		return n
	}

	var missingIds []string

	for _, id := range ids {
		if !slices.Contains(retrievedIds, id) {
			missingIds = append(missingIds, strconv.FormatInt(id, 10))
		}
	}

	if len(missingIds) != 0 {
		err = fmt.Errorf("missing %s ids: %s", aggregate, strings.Join(missingIds, ","))
		n.Add(err)
	}

	return n
}
