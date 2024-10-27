package infra_genre_test

import (
	"errors"
	"log"
	"testing"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain/genre"
	infra_genre "github.com.br/gibranct/admin_do_catalogo/internal/infra/genre"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func TestCreateGenre(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()
	cg := infra_genre.NewGenreGateway(db)
	g := genre.NewGenre("drinks")
	categoryId := int64(55)
	g.AddCategoryId(categoryId)
	query1 := "INSERT INTO genres"

	mock.ExpectBegin()
	mock.ExpectQuery(query1).WithArgs(
		g.Name, g.IsActive, g.CreatedAt, g.UpdatedAt,
	).WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	query2 := "INSERT INTO genres_categories"
	mock.ExpectExec(query2).WithoutArgs().WillReturnResult(sqlmock.NewResult(categoryId, 1))
	mock.ExpectCommit()

	err = cg.Create(g)
	if err != nil {
		log.Fatalf("Could not save genre: %s", err)
	}
	assert.Nil(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestCreateCategoryWhenFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()
	gg := infra_genre.NewGenreGateway(db)
	g := genre.NewGenre("drinks")

	query := "INSERT INTO genres"
	expectedError := errors.New("failed to create genre")

	mock.ExpectBegin()
	mock.ExpectQuery(query).WithArgs(
		g.Name, g.IsActive, g.CreatedAt, g.UpdatedAt,
	).WillReturnError(expectedError)
	mock.ExpectRollback()

	err = gg.Create(g)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestFindAllGenres(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()
	genre1 := genre.NewGenre("movie")
	genre1.ID = int64(1)
	categoryId1 := int64(23)
	genre2 := genre.NewGenre("tv show")
	genre2.ID = int64(2)
	categoryId2 := int64(3)
	genre3 := genre.NewGenre("documentary")
	categoryId3 := int64(2)
	genre3.ID = int64(3)

	gg := infra_genre.NewGenreGateway(db)
	rows := sqlmock.NewRows([]string{"1", "2", "3", "4", "5", "6", "7"})
	rows.AddRow(
		genre1.ID,
		genre1.Name,
		genre1.IsActive,
		genre1.CreatedAt,
		genre1.UpdatedAt,
		genre1.DeletedAt,
		categoryId1,
	)
	rows.AddRow(
		genre2.ID,
		genre2.Name,
		genre2.IsActive,
		genre2.CreatedAt,
		genre2.UpdatedAt,
		genre2.DeletedAt,
		categoryId2,
	)
	rows.AddRow(
		genre3.ID,
		genre3.Name,
		genre3.IsActive,
		genre3.CreatedAt,
		genre3.UpdatedAt,
		genre3.DeletedAt,
		categoryId3,
	)
	mock.ExpectQuery("SELECT").WithArgs().WillReturnRows(rows)

	genres, err := gg.FindAll()

	assert.Nil(t, err)
	assert.Len(t, genres, 3)
	assert.Len(t, genres[0].CategoryIds, 1)
	assert.Len(t, genres[1].CategoryIds, 1)
	assert.Len(t, genres[2].CategoryIds, 1)
}

func TestExistsByIds(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()
	gg := infra_genre.NewGenreGateway(db)
	genre1 := genre.NewGenre("movie")
	genre2 := genre.NewGenre("tv show")
	genre3 := genre.NewGenre("documentary")
	ids := []int64{genre1.ID, genre2.ID, genre3.ID}

	rows := sqlmock.NewRows([]string{"1"})
	rows.AddRow(ids[0])
	rows.AddRow(ids[1])
	rows.AddRow(ids[2])
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	foundIds, err := gg.ExistsByIds(ids)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(foundIds))
	assert.Equal(t, genre1.ID, foundIds[0])
	assert.Equal(t, genre2.ID, foundIds[1])
	assert.Equal(t, genre3.ID, foundIds[2])
}

func TestDeleteGenre(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()
	gg := infra_genre.NewGenreGateway(db)
	genreId := int64(56)

	mock.ExpectExec("DELETE").WithArgs(genreId).WillReturnResult(sqlmock.NewResult(int64(1), 1))

	err = gg.DeleteById(genreId)

	assert.Nil(t, err)
}

func TestDeleteGenreWhenFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()
	gg := infra_genre.NewGenreGateway(db)
	genreId := int64(56)
	erro := errors.New("failed to delete genre")

	mock.ExpectExec("DELETE").WithArgs(genreId).WillReturnError(erro)

	err = gg.DeleteById(genreId)

	assert.NotNil(t, err)
	assert.Equal(t, erro, err)
}
