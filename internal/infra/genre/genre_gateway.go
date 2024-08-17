package infra_genre

import (
	"database/sql"
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
	query := `
		INSERT INTO genres (name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	args := []any{c.Name, c.IsActive, c.CreatedAt, c.UpdatedAt}

	return cg.Db.QueryRow(query, args...).Scan(&c.ID)
}

func (cg *GenreGateway) FindAll() ([]*genre.Genre, error) {
	sql := `SELECT id, name, is_active, created_at, updated_at, deleted_at FROM genres`

	rows, err := cg.Db.Query(sql)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	genres := []*genre.Genre{}

	for rows.Next() {
		var g genre.Genre
		err := rows.Scan(
			&g.ID,
			&g.Name,
			&g.IsActive,
			&g.CreatedAt,
			&g.UpdatedAt,
			&g.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		genres = append(genres, &g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
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
