package infra_genre

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com.br/gibranct/admin-do-catalogo/internal/domain/genre"
)

type GenreGateway struct {
	Db *sql.DB
}

func NewGenreGateway(db *sql.DB) *GenreGateway {
	return &GenreGateway{Db: db}
}

func (cg *GenreGateway) Create(c *genre.Genre) error {

	tx, err := cg.Db.Begin()

	if err != nil {
		return errors.New("unable to create transaction")
	}

	defer tx.Rollback()

	query1 := `
		INSERT INTO genres (name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	query2 := `INSERT INTO genres_categories (genre_id, category_id) VALUES (%d, %d)`

	args := []any{c.Name, c.IsActive, c.CreatedAt, c.UpdatedAt}

	err = tx.QueryRow(query1, args...).Scan(&c.ID)

	if err != nil {
		return err
	}

	for _, cId := range c.CategoryIds {
		_, err = tx.Exec(fmt.Sprintf(query2, c.ID, cId))
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (cg *GenreGateway) FindAll() ([]*genre.Genre, error) {
	sql := `
		SELECT id, name, is_active, created_at, updated_at, deleted_at, gc.category_id FROM genres as g
		LEFT JOIN genres_categories as gc ON g.id = gc.genre_id
	`

	genreRows, err := cg.Db.Query(sql)

	if err != nil {
		return nil, err
	}

	defer genreRows.Close()
	genres := []*genre.Genre{}
	genresMap := make(map[int64]*genre.Genre)

	for genreRows.Next() {
		var g genre.Genre
		var categoryId int64
		err := genreRows.Scan(
			&g.ID,
			&g.Name,
			&g.IsActive,
			&g.CreatedAt,
			&g.UpdatedAt,
			&g.DeletedAt,
			&categoryId,
		)

		if err != nil {
			return nil, err
		}

		if genresMap[g.ID] == nil {
			genresMap[g.ID] = &g
		}
		genresMap[g.ID].AddCategoryId(categoryId)

	}

	if err = genreRows.Err(); err != nil {
		return nil, err
	}

	for _, v := range genresMap {
		genres = append(genres, v)
	}

	return genres, nil
}

func (cg *GenreGateway) ExistsByIds(genreIds []int64) ([]int64, error) {
	var stringIds []string
	for _, id := range genreIds {
		stringIds = append(stringIds, strconv.Itoa(int(id)))
	}

	query := fmt.Sprintf(`
		SELECT id from genres WHERE id IN (%s)
		ORDER BY id ASC
	`, strings.Join(stringIds, ","))

	rows, err := cg.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64

	for rows.Next() {
		var id int64
		err = rows.Scan(&id)

		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (cg *GenreGateway) DeleteById(genreId int64) error {
	query := "DELETE FROM genres g where g.id=$1"

	_, err := cg.Db.Exec(query, genreId)

	return err
}
