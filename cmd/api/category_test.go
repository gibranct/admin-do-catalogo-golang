package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com.br/gibranct/admin-do-catalogo/internal/domain"
	usecase "github.com.br/gibranct/admin-do-catalogo/internal/usecases"
	cs "github.com.br/gibranct/admin-do-catalogo/internal/usecases/category"
	"github.com.br/gibranct/admin-do-catalogo/pkg/notification"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CreateCategoryUseCaseMock struct {
	mock.Mock
}

func (m *CreateCategoryUseCaseMock) Execute(c cs.CreateCategoryCommand) (*notification.Notification, *cs.CreateCategoryOutput) {
	args := m.Called(c)
	return args.Get(0).(*notification.Notification), args.Get(1).(*cs.CreateCategoryOutput)
}

type ActivateCategoryUseCaseMock struct {
	mock.Mock
}

func (m *ActivateCategoryUseCaseMock) Execute(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

type DeactivateCategoryUseCaseMock struct {
	mock.Mock
}

func (m *DeactivateCategoryUseCaseMock) Execute(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

type GetCategoryByIdUseCaseMock struct {
	mock.Mock
}

func (m *GetCategoryByIdUseCaseMock) Execute(id int64) (*cs.CategoryOutput, error) {
	args := m.Called(id)
	return args.Get(0).(*cs.CategoryOutput), args.Error(1)
}

type ListCategoriesUseCaseMock struct {
	mock.Mock
}

func (m *ListCategoriesUseCaseMock) Execute(query domain.SearchQuery) (*domain.Pagination[cs.ListCategoriesOutput], error) {
	args := m.Called(query)
	return args.Get(0).(*domain.Pagination[cs.ListCategoriesOutput]), args.Error(1)
}

func TestCreateCategory(t *testing.T) {
	createUseCaseMock := new(CreateCategoryUseCaseMock)
	app := &application{
		useCases: usecase.UseCases{
			Category: usecase.CategoryUseCase{
				Create: createUseCaseMock,
			},
		},
	}
	data, _ := json.Marshal(map[string]any{
		"name":        "Test",
		"description": "test description",
	})
	command := cs.CreateCategoryCommand{
		Name:        "Test",
		Description: "test description",
	}
	var body struct {
		ID int64 `json:"id"`
	}
	id := int64(544)

	req := httptest.NewRequest(http.MethodPost, "/v1/categories", bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	createUseCaseMock.On("Execute", command).Return(&notification.Notification{}, &cs.CreateCategoryOutput{ID: id})

	app.createCategoryHandler(w, req)

	err := json.NewDecoder(w.Body).Decode(&body)

	assert.Nil(t, err)
	assert.Equal(t, id, body.ID)
	assert.Equal(t, http.StatusCreated, w.Code)
	createUseCaseMock.AssertExpectations(t)
	createUseCaseMock.AssertNumberOfCalls(t, "Execute", 1)
}

func TestGetCategoryById(t *testing.T) {
	getByIdUseCaseMock := new(GetCategoryByIdUseCaseMock)
	app := &application{
		useCases: usecase.UseCases{
			Category: usecase.CategoryUseCase{
				FindOne: getByIdUseCaseMock,
			},
		},
	}
	id := int64(488)
	cate := cs.CategoryOutput{
		ID:          id,
		Name:        "Name",
		Description: "desc",
		IsActive:    true,
	}
	body := cs.CategoryOutput{}
	req := httptest.NewRequest(http.MethodGet, "/v1/categories/{id}", nil)
	w := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.FormatInt(id, 10))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	getByIdUseCaseMock.On("Execute", id).Return(&cate, nil)

	app.getCategoryByIdHandler(w, req)

	err := json.NewDecoder(w.Body).Decode(&body)

	assert.Nil(t, err)
	assert.Equal(t, id, body.ID)
	assert.Equal(t, cate.Name, body.Name)
	assert.Equal(t, cate.Description, body.Description)
	assert.Equal(t, cate.IsActive, body.IsActive)
	assert.Equal(t, http.StatusOK, w.Code)
	getByIdUseCaseMock.AssertExpectations(t)
	getByIdUseCaseMock.AssertNumberOfCalls(t, "Execute", 1)
}

func TestGetCategoryByIdWhenIdIsString(t *testing.T) {
	getByIdUseCaseMock := new(GetCategoryByIdUseCaseMock)
	app := &application{
		useCases: usecase.UseCases{
			Category: usecase.CategoryUseCase{
				FindOne: getByIdUseCaseMock,
			},
		},
	}
	id := int64(488)

	var body struct {
		Message string   `json:"message"`
		Errors  []string `json:"errors"`
	}
	req := httptest.NewRequest(http.MethodGet, "/v1/categories/{id}", nil)
	w := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "a")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	getByIdUseCaseMock.On("Execute", id).Return(nil, errors.New(""))

	app.getCategoryByIdHandler(w, req)

	err := json.NewDecoder(w.Body).Decode(&body)

	assert.Nil(t, err)
	assert.Equal(t, "invalid id", body.Message)
	assert.Equal(t, 0, len(body.Errors))
	assert.Equal(t, http.StatusBadRequest, w.Code)
	getByIdUseCaseMock.AssertNumberOfCalls(t, "Execute", 0)
}

func TestGetCategoryByIdWhenIdIsInvalid(t *testing.T) {
	getByIdUseCaseMock := new(GetCategoryByIdUseCaseMock)
	app := &application{
		useCases: usecase.UseCases{
			Category: usecase.CategoryUseCase{
				FindOne: getByIdUseCaseMock,
			},
		},
	}
	id := int64(-1)

	var body struct {
		Message string   `json:"message"`
		Errors  []string `json:"errors"`
	}
	req := httptest.NewRequest(http.MethodGet, "/v1/categories/{id}", nil)
	w := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.FormatInt(id, 10))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	getByIdUseCaseMock.On("Execute", id).Return(nil, errors.New(""))

	app.getCategoryByIdHandler(w, req)

	err := json.NewDecoder(w.Body).Decode(&body)

	assert.Nil(t, err)
	assert.Equal(t, "the requested resource could not be found", body.Message)
	assert.Equal(t, 0, len(body.Errors))
	assert.Equal(t, http.StatusNotFound, w.Code)
	getByIdUseCaseMock.AssertNumberOfCalls(t, "Execute", 0)
}

func TestActivateCategory(t *testing.T) {
	activateCategoryUseCaseMock := new(ActivateCategoryUseCaseMock)
	app := &application{
		useCases: usecase.UseCases{
			Category: usecase.CategoryUseCase{
				Activate: activateCategoryUseCaseMock,
			},
		},
	}
	id := int64(488)
	req := httptest.NewRequest(http.MethodPost, "/v1/categories/{id}/activate", nil)
	w := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.FormatInt(id, 10))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	activateCategoryUseCaseMock.On("Execute", id).Return(nil)

	app.activateCategoryHandler(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	activateCategoryUseCaseMock.AssertExpectations(t)
	activateCategoryUseCaseMock.AssertNumberOfCalls(t, "Execute", 1)
}

func TestDeactivateCategory(t *testing.T) {
	deleteCategoryUseCaseMock := new(DeactivateCategoryUseCaseMock)
	app := &application{
		useCases: usecase.UseCases{
			Category: usecase.CategoryUseCase{
				Deactivate: deleteCategoryUseCaseMock,
			},
		},
	}
	id := int64(488)
	req := httptest.NewRequest(http.MethodPost, "/v1/categories/{id}/deactivate", nil)
	w := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.FormatInt(id, 10))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	deleteCategoryUseCaseMock.On("Execute", id).Return(nil)

	app.deactivateCategoryHandler(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	deleteCategoryUseCaseMock.AssertExpectations(t)
	deleteCategoryUseCaseMock.AssertNumberOfCalls(t, "Execute", 1)
}

func TestListCategories(t *testing.T) {
	listCategoriesUseCaseMock := new(ListCategoriesUseCaseMock)
	app := &application{
		useCases: usecase.UseCases{
			Category: usecase.CategoryUseCase{
				FindAll: listCategoriesUseCaseMock,
			},
		},
	}
	search := ""
	page := 0
	perPage := 10
	sort := "name"
	direction := "ASC"
	req := httptest.NewRequest(http.MethodGet, "/v1/categories", nil)
	query := req.URL.Query()
	query.Add("search", search)
	query.Add("page", strconv.Itoa(page))
	query.Add("perPage", strconv.Itoa(perPage))
	query.Add("sort", sort)
	query.Add("dir", direction)
	req.URL.RawQuery = query.Encode()
	output := []*cs.ListCategoriesOutput{
		{
			ID:          1,
			Name:        "name1",
			Description: "desc1",
		},
		{
			ID:          2,
			Name:        "name2",
			Description: "desc2",
		},
	}
	pageCategories := domain.Pagination[cs.ListCategoriesOutput]{
		CurrentPage: page,
		PerPage:     perPage,
		Items:       output,
		Total:       perPage,
	}
	var body struct {
		CurrentPage int                       `json:"currentPage"`
		PerPage     int                       `json:"perPage"`
		Total       int                       `json:"total"`
		Items       []cs.ListCategoriesOutput `json:"items"`
	}

	w := httptest.NewRecorder()

	listCategoriesUseCaseMock.On("Execute", domain.SearchQuery{
		Term:      search,
		Page:      page,
		PerPage:   perPage,
		Sort:      sort,
		Direction: direction,
	}).Return(&pageCategories, nil)

	app.listCategoriesHandler(w, req)
	err := json.NewDecoder(w.Body).Decode(&body)

	assert.Nil(t, err)
	assert.Equal(t, page, body.CurrentPage)
	assert.Equal(t, perPage, body.PerPage)
	assert.Equal(t, perPage, body.Total)
	assert.Equal(t, http.StatusOK, w.Code)
	listCategoriesUseCaseMock.AssertExpectations(t)
	listCategoriesUseCaseMock.AssertNumberOfCalls(t, "Execute", 1)
	for idx, item := range body.Items {
		assert.Equal(t, output[idx].ID, item.ID)
		assert.Equal(t, output[idx].Name, item.Name)
		assert.Equal(t, output[idx].Description, item.Description)
	}
}

func runTestServer() *httptest.Server {
	var cfg config = config{
		port: 4000,
		env:  "test",
	}
	app := &application{
		logger:   slog.New(slog.NewTextHandler(os.Stdout, nil)),
		useCases: usecase.NewUseCases(nil),
		config:   cfg,
	}
	return httptest.NewServer(app.routes())
}

func TestCreateCategory2(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()

	t.Run("should 201 when creation is success", func(t *testing.T) {
		data, _ := json.Marshal(map[string]any{
			"name":        "Test",
			"description": "test description",
		})
		_, err := http.Post(
			fmt.Sprintf("%s/v1/categories", ts.URL),
			"application/json",
			bytes.NewBuffer(data),
		)
		assert.Nil(t, err)
	})
}
