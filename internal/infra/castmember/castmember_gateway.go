package infra_castmember

import (
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com.br/gibranct/admin-do-catalogo/internal/domain"
	"github.com.br/gibranct/admin-do-catalogo/internal/domain/castmember"
)

type CastMemberGateway struct {
	Db *sql.DB
}

func NewCastMemberGateway(db *sql.DB) *CastMemberGateway {
	return &CastMemberGateway{Db: db}
}

func (cg *CastMemberGateway) Create(c *castmember.CastMember) error {
	query := `
		INSERT INTO cast_members (name, type, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	args := []any{c.Name, c.Type.String(), c.CreatedAt, c.UpdatedAt}

	return cg.Db.QueryRow(query, args...).Scan(&c.ID)
}

func (cg *CastMemberGateway) FindById(castMemberId int64) (*castmember.CastMember, error) {
	query := `
	 SELECT id, name, type, created_at, updated_at, FROM
	 cast_members 
	 where id = $1
	`

	castMember := &castmember.CastMember{}
	var castMemberType string

	err := cg.Db.QueryRow(query, castMemberId).Scan(
		&castMember.ID,
		&castMember.Name,
		&castMemberType,
		&castMember.CreatedAt,
		&castMember.UpdatedAt,
	)
	castMember.ChangeType(castMemberType)

	if err != nil {
		return nil, err
	}

	return castMember, nil
}

func (cg *CastMemberGateway) Update(c castmember.CastMember) error {
	query := `
	 UPDATE cast_members set name=$1, type=$2, updated_at=$3
	 where id = $4
	`

	args := []any{c.Name, c.Type.String(), c.UpdatedAt, c.ID}

	_, err := cg.Db.Exec(query, args...)

	return err
}

func (cg *CastMemberGateway) FindAll(query domain.SearchQuery) (*domain.Pagination[castmember.CastMember], error) {
	sql := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, name, type, created_at, updated_at
		FROM cast_members
		WHERE name ILIKE $1
		ORDER BY %s %s, id 
		LIMIT $2 OFFSET $3`,
		query.SortColumn(), query.SortDirection())

	args := []any{"%" + query.Term + "%", query.Limit(), query.Offset()}

	rows, err := cg.Db.Query(sql, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	castMembers := []*castmember.CastMember{}
	totalRecords := 0
	var castMemberType string

	for rows.Next() {
		var c castmember.CastMember
		err := rows.Scan(
			&totalRecords,
			&c.ID,
			&c.Name,
			&castMemberType,
			&c.CreatedAt,
			&c.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		c.ChangeType(castMemberType)

		castMembers = append(castMembers, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	lastPage := math.Ceil(float64(totalRecords) / float64(query.PerPage))
	return &domain.Pagination[castmember.CastMember]{
		Items:       castMembers,
		PerPage:     query.PerPage,
		CurrentPage: query.Page,
		Total:       totalRecords,
		IsLast:      lastPage == float64(query.Page),
	}, nil
}

func (cg *CastMemberGateway) ExistsByIds(castMemberIds []int64) ([]int64, error) {
	var stringIds []string
	for _, id := range castMemberIds {
		stringIds = append(stringIds, strconv.Itoa(int(id)))
	}

	query := fmt.Sprintf(`
		SELECT id from cast_members WHERE id IN (%s)
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
