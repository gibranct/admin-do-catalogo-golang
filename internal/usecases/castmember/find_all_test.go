package castmember_usecase

import (
	"testing"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain"
	"github.com.br/gibranct/admin_do_catalogo/internal/domain/castmember"
	"github.com.br/gibranct/admin_do_catalogo/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFindAllCategories(t *testing.T) {
	gatewayMock := new(mocks.CastMemberGatewayMock)
	sut := DefaultListCastMembersUseCase{
		Gateway: gatewayMock,
	}
	query := domain.SearchQuery{
		Page:      1,
		PerPage:   10,
		Term:      "aa",
		Sort:      "name",
		Direction: "ASC",
	}
	castMembers := []*castmember.CastMember{
		castmember.NewCastMember("John 1", castmember.ACTOR),
		castmember.NewCastMember("John 2", castmember.ACTOR),
		castmember.NewCastMember("John 3", castmember.DIRECTOR),
	}
	pageCategories := &domain.Pagination[castmember.CastMember]{
		CurrentPage: 1,
		PerPage:     10,
		Total:       1,
		Items:       castMembers,
	}

	gatewayMock.On("FindAll", query).Return(pageCategories, nil)

	page, err := sut.Execute(query)

	assert.Nil(t, err)
	assert.Len(t, page.Items, 3)
	assert.Equal(t, pageCategories.CurrentPage, page.CurrentPage)
	assert.Equal(t, pageCategories.PerPage, page.PerPage)
	assert.Equal(t, pageCategories.Total, page.Total)
	for idx, item := range page.Items {
		assert.Equal(t, castMembers[idx].Name, item.Name)
		assert.Equal(t, castMembers[idx].Type.String(), item.Type)
	}
}

func TestFindAllCategoriesWhenPageIsInvalid(t *testing.T) {
	gatewayMock := new(mocks.CastMemberGatewayMock)
	sut := DefaultListCastMembersUseCase{
		Gateway: gatewayMock,
	}
	tests := []struct {
		expectedMsg string
		query       domain.SearchQuery
	}{
		{
			expectedMsg: "invalid page",
			query: domain.SearchQuery{
				Page:      -1,
				PerPage:   1,
				Term:      "aa",
				Sort:      "name",
				Direction: "ASC",
			},
		},
		{
			expectedMsg: "invalid page",
			query: domain.SearchQuery{
				Page:      0,
				PerPage:   1,
				Term:      "aa",
				Sort:      "name",
				Direction: "ASC",
			},
		},
		{
			expectedMsg: "perPage should be greater than zero",
			query: domain.SearchQuery{
				Page:      1,
				PerPage:   0,
				Term:      "aa",
				Sort:      "name",
				Direction: "ASC",
			},
		},
		{
			expectedMsg: "perPage should be greater than zero",
			query: domain.SearchQuery{
				Page:      1,
				PerPage:   -1,
				Term:      "aa",
				Sort:      "name",
				Direction: "ASC",
			},
		},
		{
			expectedMsg: "can only sort by 'name' and 'description'",
			query: domain.SearchQuery{
				Page:      1,
				PerPage:   10,
				Term:      "aa",
				Sort:      "nae",
				Direction: "ASC",
			},
		},
		{
			expectedMsg: "invalid direction",
			query: domain.SearchQuery{
				Page:      1,
				PerPage:   10,
				Term:      "aa",
				Sort:      "name",
				Direction: "AC",
			},
		},
	}

	for _, data := range tests {
		page, err := sut.Execute(data.query)

		assert.NotNil(t, err)
		assert.Equal(t, data.expectedMsg, err.Error())
		assert.Nil(t, page)
	}

}
