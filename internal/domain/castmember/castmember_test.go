package castmember_test

import (
	"strings"
	"testing"
	"time"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain/castmember"
	"github.com.br/gibranct/admin_do_catalogo/pkg/notification"
	"github.com/stretchr/testify/assert"
)

func TestCastMemberConstructor(t *testing.T) {
	name := "John Brown"
	aType := castmember.ACTOR
	cm := castmember.NewCastMember(name, aType)

	assert.NotNil(t, cm)
	assert.Equal(t, name, cm.Name)
	assert.NotNil(t, aType, cm.Type)
	assert.NotNil(t, cm.CreatedAt)
	assert.NotNil(t, cm.UpdatedAt)
}

func TestCastMemberUpdate(t *testing.T) {
	name := "John Brown"
	aType := castmember.ACTOR
	cm := castmember.NewCastMember("Kevin", castmember.DIRECTOR)

	updatedAt := cm.UpdatedAt
	time.Sleep(time.Millisecond)
	cm.Update(name, aType)

	assert.Equal(t, name, cm.Name)
	assert.Equal(t, aType, cm.Type)
	assert.NotNil(t, cm.CreatedAt)
	assert.NotNil(t, cm.UpdatedAt)
	assert.True(t, cm.UpdatedAt.After(updatedAt))
}

func TestCastMemberChangeType(t *testing.T) {
	aType := castmember.ACTOR
	cm := castmember.NewCastMember("Kevin", castmember.DIRECTOR)

	err := cm.ChangeType("actor")

	assert.Nil(t, err)
	assert.Equal(t, aType, cm.Type)
	assert.NotNil(t, cm.CreatedAt)
	assert.NotNil(t, cm.UpdatedAt)
}

func TestCastMemberChangeTypeWithUnknownType(t *testing.T) {
	aType := castmember.DIRECTOR
	cm := castmember.NewCastMember("Kevin", castmember.DIRECTOR)

	err := cm.ChangeType("publisher")

	assert.NotNil(t, err)
	assert.Equal(t, aType, cm.Type)
}

func TestCastMemberValidator(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "",
			expected: "'name' should not be empty",
		},
		{
			name:     "32",
			expected: "'name' must be between 3 and 255 characters",
		},
		{
			name:     strings.Repeat("a", 256),
			expected: "'name' must be between 3 and 255 characters",
		},
		{
			name:     strings.Repeat("a", 255),
			expected: "",
		},
	}
	for _, v := range tests {
		c := castmember.NewCastMember(v.name, castmember.ACTOR)

		n := notification.CreateNotification()

		c.Validate(n)

		assert.NotNil(t, c)
		if v.expected != "" {
			assert.Equal(t, v.expected, n.GetErrors()[0].Error())
		} else {
			assert.False(t, n.HasErrors())
		}
	}
}
