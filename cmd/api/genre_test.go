package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	category_usecase "github.com.br/gibranct/admin-do-catalogo/internal/usecases/category"
	"github.com.br/gibranct/admin-do-catalogo/pkg/test"
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
			"application/json",
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
			"application/json",
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
			"application/json",
			bytes.NewBuffer(data),
		)
		expecBody := `{"errors":["'name' should not be empty","'name' must be between 3 and 255 characters"],"message":"Could not save genre"}`
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, expecBody, body)
	})
}
