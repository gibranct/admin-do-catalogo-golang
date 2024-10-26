package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain/castmember"
	castmember_usecase "github.com.br/gibranct/admin_do_catalogo/internal/usecases/castmember"
	"github.com.br/gibranct/admin_do_catalogo/pkg/test"
	"github.com/stretchr/testify/assert"
)

const conTypeApplicationJson = "application/json"

func TestCreateCastMember(t *testing.T) {
	t.Cleanup(cleanUp)
	ts, _ := runTestServer()
	defer ts.Close()

	t.Run("should return 201 when creation is success", func(t *testing.T) {
		data, _ := json.Marshal(map[string]any{
			"name": "Test",
			"type": "actor",
		})
		resp, err := http.Post(
			fmt.Sprintf("%s/v1/cast-members", ts.URL),
			conTypeApplicationJson,
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
			"name": "",
			"type": "actor",
		})
		resp, err := http.Post(
			fmt.Sprintf("%s/v1/cast-members", ts.URL),
			conTypeApplicationJson,
			bytes.NewBuffer(data),
		)
		expecBody := `{"errors":["'name' should not be empty","'name' must be between 3 and 255 characters"],"message":"Could not save cast member"}`
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, expecBody, body)
	})
}

func TestFindAllCastMembers(t *testing.T) {
	t.Cleanup(cleanUp)
	ts, app := runTestServer()
	defer ts.Close()
	command1 := castmember_usecase.CreateCastMemberCommand{
		Name: "test 1 b",
		Type: castmember.ACTOR,
	}
	command2 := castmember_usecase.CreateCastMemberCommand{
		Name: "test 1 a",
		Type: castmember.DIRECTOR,
	}
	_, castm1 := app.useCases.CastMember.Create.Execute(command1)
	_, castm2 := app.useCases.CastMember.Create.Execute(command2)

	t.Run("should return 200 when find all cast members without filter sorted by name ASC", func(t *testing.T) {
		page := 1
		perPage := 2
		total := 2
		isLast := true
		sort := "name"
		dir := "ASC"
		resp, err := http.Get(
			fmt.Sprintf("%s/v1/cast-members?page=%d&perPage=%d&sort=%s&dir=%s", ts.URL, page, perPage, sort, dir),
		)
		expecBody := fmt.Sprintf(
			`{"currentPage":%d,"perPage":%d,"total":%d,"isLast":%t,"items":[{"id":%d,"name":"%s","type":"%s"},{"id":%d,"name":"%s","type":"%s"}]}`,
			page, perPage, total, isLast, castm2.ID, command2.Name, command2.Type.String(), castm1.ID, command1.Name, command1.Type.String(),
		)
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expecBody, body)
	})

	t.Run("should return 200 when find all cast members without filter sorted by name DESC", func(t *testing.T) {
		page := 1
		perPage := 2
		total := 2
		isLast := true
		sort := "name"
		dir := "DESC"
		resp, err := http.Get(
			fmt.Sprintf("%s/v1/cast-members?page=%d&perPage=%d&sort=%s&dir=%s", ts.URL, page, perPage, sort, dir),
		)
		expecBody := fmt.Sprintf(
			`{"currentPage":%d,"perPage":%d,"total":%d,"isLast":%t,"items":[{"id":%d,"name":"%s","type":"%s"},{"id":%d,"name":"%s","type":"%s"}]}`,
			page, perPage, total, isLast, castm1.ID, command1.Name, command1.Type.String(), castm2.ID, command2.Name, command2.Type.String(),
		)
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expecBody, body)
	})

	t.Run("should return 200 when find cast members fist page without filter", func(t *testing.T) {
		page := 1
		perPage := 1
		total := 2
		isLast := false
		resp, err := http.Get(
			fmt.Sprintf("%s/v1/cast-members?page=%d&perPage=%d", ts.URL, page, perPage),
		)
		expecBody := fmt.Sprintf(
			`{"currentPage":%d,"perPage":%d,"total":%d,"isLast":%t,"items":[{"id":%d,"name":"%s","type":"%s"}]}`,
			page, perPage, total, isLast, castm2.ID, command2.Name, command2.Type.String(),
		)
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expecBody, body)
	})

	t.Run("should return 200 when find cast members second page without filter", func(t *testing.T) {
		page := 2
		perPage := 1
		total := 2
		isLast := true
		resp, err := http.Get(
			fmt.Sprintf("%s/v1/cast-members?page=%d&perPage=%d", ts.URL, page, perPage),
		)
		expecBody := fmt.Sprintf(
			`{"currentPage":%d,"perPage":%d,"total":%d,"isLast":%t,"items":[{"id":%d,"name":"%s","type":"%s"}]}`,
			page, perPage, total, isLast, castm1.ID, command1.Name, command1.Type.String(),
		)
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expecBody, body)
	})
}

func TestUpdateCastMember(t *testing.T) {
	t.Cleanup(cleanUp)
	ts, app := runTestServer()
	defer ts.Close()

	t.Run("should return 200 when cast members is updated", func(t *testing.T) {
		command := castmember_usecase.CreateCastMemberCommand{
			Name: "test 1",
			Type: castmember.ACTOR,
		}
		_, output := app.useCases.CastMember.Create.Execute(command)
		newName := "new name"
		newType := "actor"
		data, _ := json.Marshal(map[string]any{
			"name": newName,
			"type": newType,
		})
		req, err := http.NewRequest("PUT",
			fmt.Sprintf("%s/v1/cast-members/%d", ts.URL, output.ID),
			bytes.NewBuffer(data),
		)
		resp, _ := http.DefaultClient.Do(req)
		expecBody := fmt.Sprintf(`{"id":%d}`, output.ID)
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expecBody, body)
	})

	t.Run("should return 400 when cast member id is invalid", func(t *testing.T) {
		data, _ := json.Marshal(map[string]any{
			"name": "new name",
			"type": "actor",
		})
		req, err := http.NewRequest("PUT",
			fmt.Sprintf("%s/v1/cast-members/%s", ts.URL, "9a"),
			bytes.NewBuffer(data),
		)
		resp, _ := http.DefaultClient.Do(req)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 404 when cast member id is less than one", func(t *testing.T) {
		data, _ := json.Marshal(map[string]any{
			"name": "new name",
			"type": "actor",
		})
		req, err := http.NewRequest("PUT",
			fmt.Sprintf("%s/v1/cast-members/%d", ts.URL, 0),
			bytes.NewBuffer(data),
		)
		resp, _ := http.DefaultClient.Do(req)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
