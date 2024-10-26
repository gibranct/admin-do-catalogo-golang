package mocks

import (
	"github.com.br/gibranct/admin_do_catalogo/internal/domain"
	"github.com.br/gibranct/admin_do_catalogo/internal/domain/castmember"
	"github.com/stretchr/testify/mock"
)

type CastMemberGatewayMock struct {
	mock.Mock
}

func (m *CastMemberGatewayMock) Create(c *castmember.CastMember) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *CastMemberGatewayMock) FindById(castMemberId int64) (*castmember.CastMember, error) {
	args := m.Called(castMemberId)
	return args.Get(0).(*castmember.CastMember), args.Error(1)
}

func (m *CastMemberGatewayMock) Update(c castmember.CastMember) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *CastMemberGatewayMock) FindAll(query domain.SearchQuery) (*domain.Pagination[castmember.CastMember], error) {
	args := m.Called(query)
	return args.Get(0).(*domain.Pagination[castmember.CastMember]), args.Error(1)

}

func (m *CastMemberGatewayMock) ExistsByIds(castMemberIds []int64) ([]int64, error) {
	args := m.Called(castMemberIds)
	return args.Get(0).([]int64), args.Error(1)
}

func (m *CastMemberGatewayMock) DeleteById(castMemberId int64) error {
	args := m.Called(castMemberId)
	return args.Error(0)
}
