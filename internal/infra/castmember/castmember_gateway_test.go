package infra_castmember_test

import (
	"errors"
	"log"
	"testing"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain"
	"github.com.br/gibranct/admin_do_catalogo/internal/domain/castmember"
	infra_castmember "github.com.br/gibranct/admin_do_catalogo/internal/infra/castmember"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func TestCreateCastMember(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()
	cg := infra_castmember.NewCastMemberGateway(db)
	c := castmember.NewCastMember("John Doe", castmember.ACTOR)
	query := "INSERT INTO cast_members"
	mock.ExpectQuery(query).WithArgs(
		c.Name, c.Type.String(), c.CreatedAt, c.UpdatedAt,
	).WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	err = cg.Create(c)
	if err != nil {
		log.Fatalf("Could not save cast member: %s", err)
	}
	assert.Nil(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestCreateCastMemberWhenFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()
	cg := infra_castmember.NewCastMemberGateway(db)
	c := castmember.NewCastMember("John Doe", castmember.DIRECTOR)
	query := "INSERT INTO cast_members"
	expectedError := errors.New("failed to create cast member")
	mock.ExpectQuery(query).WithArgs(
		c.Name, c.Type.String(), c.CreatedAt, c.UpdatedAt,
	).WillReturnError(expectedError)

	err = cg.Create(c)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestFindById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()
	cg := infra_castmember.NewCastMemberGateway(db)
	c := castmember.NewCastMember("John Doe", castmember.DIRECTOR)
	c.ID = 45
	rows := sqlmock.NewRows([]string{"1", "2", "3", "4", "5"})
	rows.AddRow(
		c.ID,
		c.Name,
		c.Type.String(),
		c.CreatedAt,
		c.UpdatedAt,
	)
	mock.ExpectQuery("SELECT").WithArgs(c.ID).WillReturnRows(rows)

	castMemberFound, err := cg.FindById(c.ID)
	assert.Nil(t, err)
	assert.Equal(t, c.ID, castMemberFound.ID)
	assert.Equal(t, c.Name, castMemberFound.Name)
	assert.Equal(t, c.Type.String(), castMemberFound.Type.String())
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestFindByIdWhenFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()
	cg := infra_castmember.NewCastMemberGateway(db)
	castMember := castmember.NewCastMember("John Doe", castmember.DIRECTOR)
	castMember.ID = 45
	expectedError := errors.New("faild to find cast member")
	mock.ExpectQuery("SELECT").WithArgs(castMember.ID).WillReturnError(expectedError)

	castMemberFound, err := cg.FindById(castMember.ID)
	assert.Nil(t, castMemberFound)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()
	cg := infra_castmember.NewCastMemberGateway(db)
	c := castmember.NewCastMember("John Doe", castmember.ACTOR)
	mock.ExpectExec("UPDATE cast_members").
		WithArgs(c.Name, c.Type.String(), c.UpdatedAt, c.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = cg.Update(*c)

	assert.Nil(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUpdateWhenFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()
	cg := infra_castmember.NewCastMemberGateway(db)
	c := castmember.NewCastMember("John Doe", castmember.ACTOR)
	expectedError := errors.New("failed to update category")
	mock.ExpectExec("UPDATE cast_members").
		WithArgs(c.Name, c.Type.String(), c.UpdatedAt, c.ID).
		WillReturnError(expectedError)

	err = cg.Update(*c)

	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestFindAllWithFilters(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()
	totalRecords := 3
	castMember1 := castmember.NewCastMember("John Doe 1", castmember.ACTOR)
	castMember2 := castmember.NewCastMember("John Doe 2", castmember.ACTOR)
	castMember3 := castmember.NewCastMember("John Doe 3", castmember.ACTOR)
	test := struct {
		query         domain.SearchQuery
		expectedQuery domain.SearchQuery
		isLast        bool
	}{
		query: domain.SearchQuery{
			Page:      1,
			PerPage:   1,
			Term:      "",
			Sort:      "",
			Direction: "",
		},
		expectedQuery: domain.SearchQuery{
			Page:      1,
			PerPage:   1,
			Term:      "",
			Sort:      "name",
			Direction: "ASC",
		},
		isLast: false,
	}
	cg := infra_castmember.NewCastMemberGateway(db)
	rows := sqlmock.NewRows([]string{"1", "2", "3", "4", "5", "6"})
	rows.AddRow(
		totalRecords,
		castMember1.ID,
		castMember1.Name,
		castMember1.Type.String(),
		castMember1.CreatedAt,
		castMember1.UpdatedAt,
	)
	rows.AddRow(
		totalRecords,
		castMember2.ID,
		castMember2.Name,
		castMember2.Type.String(),
		castMember2.CreatedAt,
		castMember2.UpdatedAt,
	)
	rows.AddRow(
		totalRecords,
		castMember3.ID,
		castMember3.Name,
		castMember3.Type.String(),
		castMember3.CreatedAt,
		castMember3.UpdatedAt,
	)
	mock.ExpectQuery("SELECT").WithArgs(
		"%"+test.expectedQuery.Term+"%", test.expectedQuery.Limit(), test.expectedQuery.Offset(),
	).
		WillReturnRows(rows)

	page, err := cg.FindAll(test.query)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(page.Items))
	assert.Equal(t, 3, page.Total)
	assert.Equal(t, test.isLast, page.IsLast)
	assert.Equal(t, test.expectedQuery.Page, page.CurrentPage)
	assert.Equal(t, test.expectedQuery.PerPage, page.PerPage)
}

func TestFindAllWhenIsLastPage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()
	totalRecords := 3
	castMember1 := castmember.NewCastMember("John Doe 1", castmember.ACTOR)
	castMember2 := castmember.NewCastMember("John Doe 2", castmember.ACTOR)
	castMember3 := castmember.NewCastMember("John Doe 3", castmember.ACTOR)
	test := struct {
		query         domain.SearchQuery
		expectedQuery domain.SearchQuery
		isLast        bool
	}{
		query: domain.SearchQuery{
			Page:      1,
			PerPage:   3,
			Term:      "",
			Sort:      "",
			Direction: "",
		},
		expectedQuery: domain.SearchQuery{
			Page:      1,
			PerPage:   3,
			Term:      "",
			Sort:      "name",
			Direction: "ASC",
		},
		isLast: true,
	}
	cg := infra_castmember.NewCastMemberGateway(db)
	rows := sqlmock.NewRows([]string{"1", "2", "3", "4", "5", "6"})
	rows.AddRow(
		totalRecords,
		castMember1.ID,
		castMember1.Name,
		castMember1.Type.String(),
		castMember1.CreatedAt,
		castMember1.UpdatedAt,
	)
	rows.AddRow(
		totalRecords,
		castMember2.ID,
		castMember2.Name,
		castMember2.Type.String(),
		castMember2.CreatedAt,
		castMember2.UpdatedAt,
	)
	rows.AddRow(
		totalRecords,
		castMember3.ID,
		castMember3.Name,
		castMember3.Type.String(),
		castMember3.CreatedAt,
		castMember3.UpdatedAt,
	)
	mock.ExpectQuery("SELECT").WithArgs(
		"%"+test.expectedQuery.Term+"%", test.expectedQuery.Limit(), test.expectedQuery.Offset(),
	).
		WillReturnRows(rows)

	page, err := cg.FindAll(test.query)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(page.Items))
	assert.Equal(t, 3, page.Total)
	assert.Equal(t, test.isLast, page.IsLast)
	assert.Equal(t, test.expectedQuery.Page, page.CurrentPage)
	assert.Equal(t, test.expectedQuery.PerPage, page.PerPage)
}

func TestExistsByIds(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()
	cg := infra_castmember.NewCastMemberGateway(db)
	castMember1 := castmember.NewCastMember("John Doe 1", castmember.ACTOR)
	castMember2 := castmember.NewCastMember("John Doe 2", castmember.ACTOR)
	castMember3 := castmember.NewCastMember("John Doe 3", castmember.ACTOR)
	ids := []int64{castMember1.ID, castMember2.ID, castMember3.ID}

	rows := sqlmock.NewRows([]string{"1"})
	rows.AddRow(ids[0])
	rows.AddRow(ids[1])
	rows.AddRow(ids[2])
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	foundIds, err := cg.ExistsByIds(ids)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(foundIds))
	assert.Equal(t, castMember1.ID, foundIds[0])
	assert.Equal(t, castMember1.ID, foundIds[1])
	assert.Equal(t, castMember1.ID, foundIds[2])
}
