package infra_category

import (
	"errors"
	"log"
	"testing"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain"
	"github.com.br/gibranct/admin_do_catalogo/internal/domain/category"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func TestCreateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()
	cg := NewCategoryGateway(db)
	c := category.NewCategory("drinks", "drinks desc")
	query := "INSERT INTO categories"
	mock.ExpectQuery(query).WithArgs(
		c.Name, c.Description, c.IsActive, c.CreatedAt, c.UpdatedAt,
	).WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	err = cg.Create(c)
	if err != nil {
		log.Fatalf("Could not save category: %s", err)
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
	cg := NewCategoryGateway(db)
	c := category.NewCategory("drinks", "drinks desc")
	query := "INSERT INTO categories"
	expectedError := errors.New("failed to create category")
	mock.ExpectQuery(query).WithArgs(
		c.Name, c.Description, c.IsActive, c.CreatedAt, c.UpdatedAt,
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
	cg := NewCategoryGateway(db)
	category := category.NewCategory("drinks", "drinks desc")
	category.ID = 45
	rows := sqlmock.NewRows([]string{"1", "2", "3", "4", "5", "6", "7"})
	rows.AddRow(
		category.ID,
		category.Name,
		category.Description,
		category.IsActive,
		category.CreatedAt,
		category.UpdatedAt,
		category.DeletedAt,
	)
	mock.ExpectQuery("SELECT").WithArgs(category.ID).WillReturnRows(rows)

	categoryFound, err := cg.FindById(category.ID)
	assert.Nil(t, err)
	assert.Equal(t, category.ID, categoryFound.ID)
	assert.Equal(t, category.Name, categoryFound.Name)
	assert.Equal(t, category.Description, categoryFound.Description)
	assert.Equal(t, category.IsActive, categoryFound.IsActive)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestFindByIdWhenFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s", err)
	}
	defer db.Close()
	cg := NewCategoryGateway(db)
	category := category.NewCategory("drinks", "drinks desc")
	category.ID = 45
	expectedError := errors.New("faild to find category")
	mock.ExpectQuery("SELECT").WithArgs(category.ID).WillReturnError(expectedError)

	categoryFound, err := cg.FindById(category.ID)
	assert.Nil(t, categoryFound)
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
	cg := NewCategoryGateway(db)
	c := category.NewCategory("drinks", "drinks desc")
	mock.ExpectExec("UPDATE categories").
		WithArgs(c.Name, c.Description, c.IsActive, c.UpdatedAt, c.DeletedAt, c.ID).
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
	cg := NewCategoryGateway(db)
	c := category.NewCategory("drinks", "drinks desc")
	expectedError := errors.New("failed to update category")
	mock.ExpectExec("UPDATE categories").
		WithArgs(c.Name, c.Description, c.IsActive, c.UpdatedAt, c.DeletedAt, c.ID).
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
	category1 := category.NewCategory("movie", "drinks desc")
	category2 := category.NewCategory("tv show", "drinks desc")
	category3 := category.NewCategory("documentary", "drinks desc")
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
	cg := NewCategoryGateway(db)
	rows := sqlmock.NewRows([]string{"1", "2", "3", "4", "5", "6", "7", "8"})
	rows.AddRow(
		totalRecords,
		category1.ID,
		category1.Name,
		category1.Description,
		category1.IsActive,
		category1.CreatedAt,
		category1.UpdatedAt,
		category1.DeletedAt,
	)
	rows.AddRow(
		totalRecords,
		category2.ID,
		category2.Name,
		category2.Description,
		category2.IsActive,
		category2.CreatedAt,
		category2.UpdatedAt,
		category2.DeletedAt,
	)
	rows.AddRow(
		totalRecords,
		category3.ID,
		category3.Name,
		category3.Description,
		category3.IsActive,
		category3.CreatedAt,
		category3.UpdatedAt,
		category3.DeletedAt,
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
	category1 := category.NewCategory("movie", "drinks desc")
	category2 := category.NewCategory("tv show", "drinks desc")
	category3 := category.NewCategory("documentary", "drinks desc")
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
	cg := NewCategoryGateway(db)
	rows := sqlmock.NewRows([]string{"1", "2", "3", "4", "5", "6", "7", "8"})
	rows.AddRow(
		totalRecords,
		category1.ID,
		category1.Name,
		category1.Description,
		category1.IsActive,
		category1.CreatedAt,
		category1.UpdatedAt,
		category1.DeletedAt,
	)
	rows.AddRow(
		totalRecords,
		category2.ID,
		category2.Name,
		category2.Description,
		category2.IsActive,
		category2.CreatedAt,
		category2.UpdatedAt,
		category2.DeletedAt,
	)
	rows.AddRow(
		totalRecords,
		category3.ID,
		category3.Name,
		category3.Description,
		category3.IsActive,
		category3.CreatedAt,
		category3.UpdatedAt,
		category3.DeletedAt,
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
	cg := NewCategoryGateway(db)
	category1 := category.NewCategory("movie", "drinks desc")
	category2 := category.NewCategory("tv show", "drinks desc")
	category3 := category.NewCategory("documentary", "drinks desc")
	ids := []int64{category1.ID, category2.ID, category3.ID}

	rows := sqlmock.NewRows([]string{"1"})
	rows.AddRow(ids[0])
	rows.AddRow(ids[1])
	rows.AddRow(ids[2])
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	foundIds, err := cg.ExistsByIds(ids)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(foundIds))
	assert.Equal(t, category1.ID, foundIds[0])
	assert.Equal(t, category2.ID, foundIds[1])
	assert.Equal(t, category3.ID, foundIds[2])
}
