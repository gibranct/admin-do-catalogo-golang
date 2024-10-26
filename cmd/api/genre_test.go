package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	category_usecase "github.com.br/gibranct/admin_do_catalogo/internal/usecases/category"
	genre_usecase "github.com.br/gibranct/admin_do_catalogo/internal/usecases/genre"
	"github.com.br/gibranct/admin_do_catalogo/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateGenre(t *testing.T) {
	t.Cleanup(cleanUp)
	ts, app := runTestServer()
	defer ts.Close()

	t.Run("should return 201 when creation without categories is success", func(t *testing.T) {
		data, _ := json.Marshal(map[string]any{
			"name":        "Test",
			"categoryIds": []int{},
		})
		resp, err := http.Post(
			fmt.Sprintf("%s/v1/genres", ts.URL),
			conTypeApplicationJson,
			bytes.NewBuffer(data),
		)
		expecBody := `{"id":1}`
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, expecBody, body)
	})

	t.Run("should return 201 when creation with categories is success", func(t *testing.T) {
		command := category_usecase.CreateCategoryCommand{
			Name:        "test 1",
			Description: "desc fake",
		}
		_, output1 := app.useCases.Category.Create.Execute(command)
		data, _ := json.Marshal(map[string]any{
			"name":        "Test",
			"categoryIds": []int64{output1.ID},
		})
		resp, err := http.Post(
			fmt.Sprintf("%s/v1/genres", ts.URL),
			conTypeApplicationJson,
			bytes.NewBuffer(data),
		)
		expecBody := `{"id":2}`
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, expecBody, body)
	})

	t.Run("should return 400 when creation fails", func(t *testing.T) {
		data, _ := json.Marshal(map[string]any{
			"name":        "",
			"categoryIds": []int64{},
		})
		resp, err := http.Post(
			fmt.Sprintf("%s/v1/genres", ts.URL),
			conTypeApplicationJson,
			bytes.NewBuffer(data),
		)
		expecBody := `{"errors":["'name' should not be empty","'name' must be between 3 and 255 characters"],"message":"Could not save genre"}`
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, expecBody, body)
	})
}

func TestFindAllGenres(t *testing.T) {
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

	categoryIds := []int64{cate1.ID, cate2.ID}
	command3 := genre_usecase.CreateGenreCommand{
		Name:        "Genre 1",
		CategoryIds: &categoryIds,
	}

	_, genre := app.useCases.Genre.Create.Execute(command3)

	t.Run("should return all genres with categories", func(t *testing.T) {
		resp, err := http.Get(
			fmt.Sprintf("%s/v1/genres", ts.URL),
		)
		expecBody := fmt.Sprintf(
			`[{"id":%d,"name":"%s","active":true,"categoryIds":[%d,%d]}]`,
			genre.ID, command3.Name, cate1.ID, cate2.ID,
		)
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Regexp(t, expecBody, body)
	})
}

func TestDeleteGenreById(t *testing.T) {
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

	categoryIds := []int64{cate1.ID, cate2.ID}
	command3 := genre_usecase.CreateGenreCommand{
		Name:        "Genre 1",
		CategoryIds: &categoryIds,
	}

	_, genre := app.useCases.Genre.Create.Execute(command3)

	t.Run("should delete genre", func(t *testing.T) {
		url := fmt.Sprintf("%s/v1/genres/%d", ts.URL, genre.ID)
		deleteReq, _ := http.NewRequest(http.MethodDelete, url, nil)
		resp, err := http.DefaultClient.Do(deleteReq)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}
