package infra_video

import (
	"database/sql"
	"fmt"

	"github.com.br/gibranct/admin-do-catalogo/internal/domain/video"
)

type VideoGateway struct {
	Db *sql.DB
}

func NewVideoGateway(db *sql.DB) *VideoGateway {
	return &VideoGateway{Db: db}
}

func (vg VideoGateway) Create(aVideo video.Video) (*video.Video, error) {
	tx, err := vg.Db.Begin()

	if err != nil {
		return nil, fmt.Errorf("unable to create transaction: %s", err.Error())
	}

	defer tx.Rollback()

	var videoResourceId *int64
	var trailerResourceId *int64
	var bannerResourceId *int64
	var thumbnailResourceId *int64
	var thumbnailHalfResourceId *int64

	videoResourceId, err = saveVideoMedia(tx, aVideo.Video)
	if err != nil {
		return nil, err
	}

	trailerResourceId, err = saveVideoMedia(tx, aVideo.Trailer)
	if err != nil {
		return nil, err
	}

	bannerResourceId, err = saveImageMedia(tx, aVideo.Banner)
	if err != nil {
		return nil, err
	}

	thumbnailResourceId, err = saveImageMedia(tx, aVideo.ThumbNail)
	if err != nil {
		return nil, err
	}

	thumbnailHalfResourceId, err = saveImageMedia(tx, aVideo.ThumbNailHalf)
	if err != nil {
		return nil, err
	}

	createVideoQuery := `
		INSERT INTO VIDEOS (
		title, description, year_launched, opened, published, rating,
		duration, created_at, updated_at, video_id, trailer_id, banner_id, thumbnail_id, thumbnail_half_id)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id
	`

	var lastInsertId int64

	err = tx.QueryRow(createVideoQuery,
		aVideo.Title,
		aVideo.Description,
		aVideo.LaunchedAt,
		aVideo.Opened,
		aVideo.Published,
		aVideo.Rating.String(),
		aVideo.Duration,
		aVideo.CreatedAt,
		aVideo.UpdatedAt,
		videoResourceId,
		trailerResourceId,
		bannerResourceId,
		thumbnailResourceId,
		thumbnailHalfResourceId,
	).Scan(&lastInsertId)

	if err != nil {
		return nil, err
	}

	for _, cId := range aVideo.CategoryIds {
		err = saveCategory(tx, lastInsertId, cId)
		if err != nil {
			return nil, err
		}
	}

	for _, cId := range aVideo.CastMemberIds {
		err = saveCastMember(tx, lastInsertId, cId)
		if err != nil {
			return nil, err
		}
	}

	for _, gId := range aVideo.GenreIds {
		err = saveGenre(tx, lastInsertId, gId)
		if err != nil {
			return nil, err
		}
	}

	aVideo.ID = lastInsertId
	if videoResourceId != nil {
		aVideo.Video.ID = *videoResourceId
	}
	if trailerResourceId != nil {
		aVideo.Trailer.ID = *trailerResourceId
	}
	if bannerResourceId != nil {
		aVideo.Banner.ID = *bannerResourceId
	}
	if thumbnailResourceId != nil {
		aVideo.ThumbNail.ID = *thumbnailResourceId
	}
	if thumbnailHalfResourceId != nil {
		aVideo.ThumbNailHalf.ID = *thumbnailHalfResourceId
	}

	return &aVideo, nil
}

func (vg VideoGateway) DeleteById(aVideo int64) error {
	return nil
}

func (vg VideoGateway) FindById(videoId int64) (*video.Video, error) {
	return nil, nil
}

func saveVideoMedia(tx *sql.Tx, video *video.AudioVideoMedia) (*int64, error) {
	if video == nil {
		return nil, nil
	}

	createVideoMediaQuery := `
	INSERT INTO videos_video_media (name, checksum, file_path, encoded_path, media_status)
	VALUES ($1, $2, $3, $4, $5) RETURNING id
`

	var lastInsertId int64

	err := tx.QueryRow(createVideoMediaQuery,
		video.Name,
		video.Checksum,
		video.RawLocation,
		video.EncodedLocation,
		video.Status.String(),
	).Scan(&lastInsertId)

	if err != nil {
		return nil, err
	}

	return &lastInsertId, nil
}

func saveImageMedia(tx *sql.Tx, image *video.ImageMedia) (*int64, error) {
	if image == nil {
		return nil, nil
	}

	createImageMediaQuery := `
		INSERT INTO videos_image_media (name, checksum, file_path)
		VALUES ($1, $2, $3) RETURNING id
	`

	var lastInsertId int64

	err := tx.QueryRow(createImageMediaQuery, image.Name, image.Checksum, image.Location).Scan(&lastInsertId)

	if err != nil {
		return nil, err
	}

	return &lastInsertId, nil
}

func saveCategory(tx *sql.Tx, videoId, categoryId int64) error {
	query := `
		INSERT INTO videos_categories (video_id, category_id) VALUES ($1, $2)
	`
	_, err := tx.Exec(query, videoId, categoryId)

	return err
}

func saveGenre(tx *sql.Tx, videoId, genreId int64) error {
	query := `
		INSERT INTO videos_genres (video_id, genre_id) VALUES ($1, $2)
	`
	_, err := tx.Exec(query, videoId, genreId)

	return err
}

func saveCastMember(tx *sql.Tx, videoId, castId int64) error {
	query := `
		INSERT INTO videos_cast_members (video_id, cast_member_id) VALUES ($1, $2)
	`
	_, err := tx.Exec(query, videoId, castId)

	return err
}
