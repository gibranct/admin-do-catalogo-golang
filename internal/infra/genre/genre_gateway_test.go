package infra_genre

import (
	"errors"
	"log"
	"testing"

	"github.com.br/gibranct/admin-do-catalogo/internal/domain/genre"
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
	cg := NewGenreGateway(db)
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
	gg := NewGenreGateway(db)
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
	genre2 := genre.NewGenre("tv show")
	genre3 := genre.NewGenre("documentary")

	gg := NewGenreGateway(db)
	rows := sqlmock.NewRows([]string{"1", "2", "3", "4", "5", "6"})
	rows.AddRow(
		genre1.ID,
		genre1.Name,
		genre1.IsActive,
		genre1.CreatedAt,
		genre1.UpdatedAt,
		genre1.DeletedAt,
	)
	rows.AddRow(
		genre2.ID,
		genre2.Name,
		genre2.IsActive,
		genre2.CreatedAt,
		genre2.UpdatedAt,
		genre2.DeletedAt,
	)
	rows.AddRow(
		genre3.ID,
		genre3.Name,
		genre3.IsActive,
		genre3.CreatedAt,
		genre3.UpdatedAt,
		genre3.DeletedAt,
	)
	mock.ExpectQuery("SELECT").WithArgs().WillReturnRows(rows)

	genres, err := gg.FindAll()

	assert.Nil(t, err)
	assert.Equal(t, 3, len(genres))
}

func TestExistsByIds(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()
	gg := NewGenreGateway(db)
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
