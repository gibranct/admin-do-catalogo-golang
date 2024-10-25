package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain/castmember"
	castmember_usecase "github.com.br/gibranct/admin_do_catalogo/internal/usecases/castmember"
	category_usecase "github.com.br/gibranct/admin_do_catalogo/internal/usecases/category"
	genre_usecase "github.com.br/gibranct/admin_do_catalogo/internal/usecases/genre"
	"github.com.br/gibranct/admin_do_catalogo/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateVideo(t *testing.T) {
	t.Cleanup(cleanUp)
	ts, app := runTestServer()
	defer ts.Close()

	t.Run("should return 201 when full video creation is success", func(t *testing.T) {
		_, category := app.useCases.Category.Create.Execute(category_usecase.CreateCategoryCommand{
			Name:        "dummy name",
			Description: "dummy desc",
		})
		categoryIds := []int64{category.ID}
		_, member := app.useCases.CastMember.Create.Execute(castmember_usecase.CreateCastMemberCommand{
			Name: "dummy name",
			Type: castmember.ACTOR,
		})

		_, genre := app.useCases.Genre.Create.Execute(genre_usecase.CreateGenreCommand{
			Name:        "genre dummy",
			CategoryIds: &categoryIds,
		})
		data, _ := json.Marshal(map[string]any{
			"title":        "dummy title",
			"description":  "dummy desc",
			"yearLaunched": 2025,
			"duration":     120.0,
			"opened":       false,
			"published":    true,
			"rating":       "Livre",
			"categoryIds":  categoryIds,
			"genreIds":     []int{int(genre.ID)},
			"memberIds":    []int{int(member.ID)},
		})
		resp, err := http.Post(
			fmt.Sprintf("%s/v1/videos", ts.URL),
			"application/json",
			bytes.NewBuffer(data),
		)
		expectedBody := `{"id":1}`
		body := test.ReadRespBody(*resp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, expectedBody, body)
	})
}
