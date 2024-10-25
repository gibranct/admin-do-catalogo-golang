package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	usecase "github.com.br/gibranct/admin_do_catalogo/internal/usecases"
	category_usecase "github.com.br/gibranct/admin_do_catalogo/internal/usecases/category"
	"github.com.br/gibranct/admin_do_catalogo/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type databaseContainer struct {
	*postgres.PostgresContainer
	connectionString string
	db               *sql.DB
}

var dbContainer *databaseContainer

func TestMain(t *testing.M) {
	ctx := context.Background()
	defer func() {
		if r := recover(); r != nil {
			dbContainer.PostgresContainer.Terminate(ctx)
		}
	}()
	connString, container, _ := test.InitDatabase(ctx)
	dbContainer = &databaseContainer{
		connectionString:  connString,
		PostgresContainer: container,
	}
	code := t.Run()
	container.Terminate(ctx)
	os.Exit(code)
}

func runTestServer() (*httptest.Server, *application) {
	var cfg config = config{
		port: 4000,
		env:  "test",
		db: struct {
			dsn          string
			maxOpenConns int
			maxIdleConns int
			maxIdleTime  time.Duration
		}{
			dsn:          dbContainer.connectionString,
			maxOpenConns: 1,
			maxIdleConns: 1,
			maxIdleTime:  5 * time.Minute,
		},
	}
	db, err := OpenDB(cfg)
	if err != nil {
		panic("failed to start db connection: " + err.Error())
	}
	dbContainer.db = db
	app := &application{
		logger:   slog.New(slog.NewTextHandler(os.Stdout, nil)),
		useCases: usecase.NewUseCases(db),
		config:   cfg,
	}
	return httptest.NewServer(app.routes()), app
}

func cleanUp() {
	tx, err := dbContainer.db.Begin()
	if err != nil {
		log.Fatalf("failed to create transaction: %s", err)
	}
	tx.Exec("DELETE FROM videos_video_media")
	tx.Exec("DELETE FROM videos_image_media")
	tx.Exec("DELETE FROM videos_categories")
	tx.Exec("DELETE FROM videos_genres")
	tx.Exec("DELETE FROM videos_cast_members")
	tx.Exec("DELETE FROM videos")
	tx.Exec("DELETE FROM categories")
	tx.Exec("DELETE FROM cast_members")
	tx.Exec("DELETE FROM genres")
	err = tx.Commit()
	if err != nil {
		log.Fatalf("failed to commit: %s", err)
	}
}

func TestCreateCategory(t *testing.T) {
	t.Cleanup(cleanUp)
	ts, _ := runTestServer()
	defer ts.Close()

	t.Run("should return 201 when creation is success", func(t *testing.T) {
		data, _ := json.Marshal(map[string]any{
			"name":        "Test",
			"description": "test description",
		})
		resp, err := http.Post(
			fmt.Sprintf("%s/v1/categories", ts.URL),
			"application/json",
			bytes.NewBuffer(data),
		)
		expecBody := `{"id":1}`
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, expecBody, body)
	})

	t.Run("should return 400 when creation fails", func(t *testing.T) {
		data, _ := json.Marshal(map[string]any{
			"name":        "",
			"description": "",
		})
		resp, err := http.Post(
			fmt.Sprintf("%s/v1/categories", ts.URL),
			"application/json",
			bytes.NewBuffer(data),
		)
		expecBody := `{"errors":["'name' should not be empty","'name' must be between 3 and 255 characters"],"message":"Could not save category"}`
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, expecBody, body)
	})
}

func TestFindCategoryById(t *testing.T) {
	t.Cleanup(cleanUp)
	ts, app := runTestServer()
	defer ts.Close()

	t.Run("should return 200 when category exists", func(t *testing.T) {
		command := category_usecase.CreateCategoryCommand{
			Name:        "test 1",
			Description: "desc fake",
		}
		_, output := app.useCases.Category.Create.Execute(command)
		resp, err := http.Get(
			fmt.Sprintf("%s/v1/categories/%d", ts.URL, output.ID),
		)
		expecBody := fmt.Sprintf(
			`{"id":%d,"name":"%s","description":"%s","active":true}`,
			output.ID,
			command.Name,
			command.Description,
		)
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expecBody, body)
	})

	t.Run("should return 404 when category does not exists", func(t *testing.T) {
		resp, err := http.Get(
			fmt.Sprintf("%s/v1/categories/%d", ts.URL, 999),
		)
		expecBody := `{"errors":[],"message":"the requested resource could not be found"}`
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Equal(t, expecBody, body)
	})
}

func TestDeactivateCategory(t *testing.T) {
	t.Cleanup(cleanUp)
	ts, app := runTestServer()
	defer ts.Close()

	t.Run("should return 204 when category is deactivated", func(t *testing.T) {
		command := category_usecase.CreateCategoryCommand{
			Name:        "test 1",
			Description: "desc fake",
		}
		_, output := app.useCases.Category.Create.Execute(command)
		resp, err := http.Post(
			fmt.Sprintf("%s/v1/categories/%d/deactivate", ts.URL, output.ID),
			"application/json",
			nil,
		)
		assert.Nil(t, err)
		out, err := app.useCases.Category.FindOne.Execute(output.ID)

		assert.False(t, out.IsActive)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

func TestActivateCategory(t *testing.T) {
	t.Cleanup(cleanUp)
	ts, app := runTestServer()
	defer ts.Close()

	t.Run("should return 204 when category is activated", func(t *testing.T) {
		command := category_usecase.CreateCategoryCommand{
			Name:        "test 1",
			Description: "desc fake",
		}
		_, output := app.useCases.Category.Create.Execute(command)
		_, err := http.Post(
			fmt.Sprintf("%s/v1/categories/%d/deactivate", ts.URL, output.ID),
			"application/json",
			nil,
		)
		assert.Nil(t, err)
		out, err := app.useCases.Category.FindOne.Execute(output.ID)
		assert.Nil(t, err)
		assert.False(t, out.IsActive)
		resp, err := http.Post(
			fmt.Sprintf("%s/v1/categories/%d/activate", ts.URL, output.ID),
			"application/json",
			nil,
		)
		assert.Nil(t, err)
		outActive, err := app.useCases.Category.FindOne.Execute(out.ID)
		assert.Nil(t, err)
		assert.True(t, outActive.IsActive)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

func TestFindAllCategories(t *testing.T) {
	t.Cleanup(cleanUp)
	ts, app := runTestServer()
	defer ts.Close()
	command1 := category_usecase.CreateCategoryCommand{
		Name:        "test 1 b",
		Description: "desc fake a",
	}
	command2 := category_usecase.CreateCategoryCommand{
		Name:        "test 1 a",
		Description: "desc fake b",
	}
	_, cate1 := app.useCases.Category.Create.Execute(command1)
	_, cate2 := app.useCases.Category.Create.Execute(command2)

	t.Run("should return 200 when find all categories without filter sorted by name ASC", func(t *testing.T) {
		page := 1
		perPage := 2
		total := 2
		isLast := true
		sort := "name"
		dir := "ASC"
		resp, err := http.Get(
			fmt.Sprintf("%s/v1/categories?page=%d&perPage=%d&sort=%s&dir=%s", ts.URL, page, perPage, sort, dir),
		)
		expecBody := fmt.Sprintf(
			`{"currentPage":%d,"perPage":%d,"total":%d,"isLast":%t,"items":[{"id":%d,"name":"%s","description":"%s"},{"id":%d,"name":"%s","description":"%s"}]}`,
			page, perPage, total, isLast, cate2.ID, command2.Name, command2.Description, cate1.ID, command1.Name, command1.Description,
		)
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expecBody, body)
	})

	t.Run("should return 200 when find all categories without filter sorted by name DESC", func(t *testing.T) {
		page := 1
		perPage := 2
		total := 2
		isLast := true
		sort := "name"
		dir := "DESC"
		resp, err := http.Get(
			fmt.Sprintf("%s/v1/categories?page=%d&perPage=%d&sort=%s&dir=%s", ts.URL, page, perPage, sort, dir),
		)
		expecBody := fmt.Sprintf(
			`{"currentPage":%d,"perPage":%d,"total":%d,"isLast":%t,"items":[{"id":%d,"name":"%s","description":"%s"},{"id":%d,"name":"%s","description":"%s"}]}`,
			page, perPage, total, isLast, cate1.ID, command1.Name, command1.Description, cate2.ID, command2.Name, command2.Description,
		)
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expecBody, body)
	})

	t.Run("should return 200 when find all categories without filter sorted by description DESC", func(t *testing.T) {
		page := 1
		perPage := 2
		total := 2
		isLast := true
		sort := "description"
		dir := "DESC"
		resp, err := http.Get(
			fmt.Sprintf("%s/v1/categories?page=%d&perPage=%d&sort=%s&dir=%s", ts.URL, page, perPage, sort, dir),
		)
		expecBody := fmt.Sprintf(
			`{"currentPage":%d,"perPage":%d,"total":%d,"isLast":%t,"items":[{"id":%d,"name":"%s","description":"%s"},{"id":%d,"name":"%s","description":"%s"}]}`,
			page, perPage, total, isLast, cate2.ID, command2.Name, command2.Description, cate1.ID, command1.Name, command1.Description,
		)
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expecBody, body)
	})

	t.Run("should return 200 when find categories fist page without filter", func(t *testing.T) {
		page := 1
		perPage := 1
		total := 2
		isLast := false
		resp, err := http.Get(
			fmt.Sprintf("%s/v1/categories?page=%d&perPage=%d", ts.URL, page, perPage),
		)
		expecBody := fmt.Sprintf(
			`{"currentPage":%d,"perPage":%d,"total":%d,"isLast":%t,"items":[{"id":%d,"name":"%s","description":"%s"}]}`,
			page, perPage, total, isLast, cate2.ID, command2.Name, command2.Description,
		)
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expecBody, body)
	})

	t.Run("should return 200 when find categories second page without filter", func(t *testing.T) {
		page := 2
		perPage := 1
		total := 2
		isLast := true
		resp, err := http.Get(
			fmt.Sprintf("%s/v1/categories?page=%d&perPage=%d", ts.URL, page, perPage),
		)
		expecBody := fmt.Sprintf(
			`{"currentPage":%d,"perPage":%d,"total":%d,"isLast":%t,"items":[{"id":%d,"name":"%s","description":"%s"}]}`,
			page, perPage, total, isLast, cate1.ID, command1.Name, command1.Description,
		)
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expecBody, body)
	})
}

func TestUpdateCategory(t *testing.T) {
	t.Cleanup(cleanUp)
	ts, app := runTestServer()
	defer ts.Close()

	t.Run("should return 200 when category is updated", func(t *testing.T) {
		command := category_usecase.CreateCategoryCommand{
			Name:        "test 1",
			Description: "desc fake",
		}
		_, output := app.useCases.Category.Create.Execute(command)
		newName := "new name"
		newDesc := "new description"
		data, _ := json.Marshal(map[string]any{
			"name":        newName,
			"description": newDesc,
		})
		req, err := http.NewRequest("PUT",
			fmt.Sprintf("%s/v1/categories/%d", ts.URL, output.ID),
			bytes.NewBuffer(data),
		)
		resp, _ := http.DefaultClient.Do(req)
		expecBody := fmt.Sprintf(`{"id":%d}`, output.ID)
		body := test.ReadRespBody(*resp)

		c, _ := app.useCases.Category.FindOne.Execute(output.ID)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expecBody, body)
		assert.Equal(t, newName, c.Name)
		assert.Equal(t, newDesc, c.Description)
		assert.Equal(t, output.ID, c.ID)
		assert.True(t, c.IsActive)
	})

	t.Run("should return 400 when category id is invalid", func(t *testing.T) {
		data, _ := json.Marshal(map[string]any{
			"name":        "new name",
			"description": "new description",
		})
		req, err := http.NewRequest("PUT",
			fmt.Sprintf("%s/v1/categories/%s", ts.URL, "9a"),
			bytes.NewBuffer(data),
		)
		resp, _ := http.DefaultClient.Do(req)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 404 when category id is less than one", func(t *testing.T) {
		data, _ := json.Marshal(map[string]any{
			"name":        "new name",
			"description": "new description",
		})
		req, err := http.NewRequest("PUT",
			fmt.Sprintf("%s/v1/categories/%d", ts.URL, 0),
			bytes.NewBuffer(data),
		)
		resp, _ := http.DefaultClient.Do(req)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
