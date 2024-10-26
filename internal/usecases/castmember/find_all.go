package castmember_usecase

import (
	"github.com.br/gibranct/admin_do_catalogo/internal/domain"
	"github.com.br/gibranct/admin_do_catalogo/internal/domain/castmember"
)

type ListCastMembersOutput struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ListCastMembersUseCase interface {
	Execute(query domain.SearchQuery) (*domain.Pagination[ListCastMembersOutput], error)
}

type DefaultListCastMembersUseCase struct {
	Gateway castmember.CastMemberGateway
}

func (useCase *DefaultListCastMembersUseCase) Execute(query domain.SearchQuery) (*domain.Pagination[ListCastMembersOutput], error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}

	page, err := useCase.Gateway.FindAll(query)

	if err != nil {
		return nil, err
	}

	var outputs []*ListCastMembersOutput

	for _, item := range page.Items {
		output := &ListCastMembersOutput{
			ID:   item.ID,
			Name: item.Name,
			Type: item.Type.String(),
		}

		outputs = append(outputs, output)
	}

	return &domain.Pagination[ListCastMembersOutput]{
		Items:       outputs,
		CurrentPage: page.CurrentPage,
		PerPage:     page.PerPage,
		Total:       page.Total,
		IsLast:      page.IsLast,
	}, nil

}
