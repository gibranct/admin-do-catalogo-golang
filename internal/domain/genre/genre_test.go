package genre

import (
	"strings"
	"testing"
	"time"

	"github.com.br/gibranct/admin-do-catalogo/pkg/notification"
	"github.com/stretchr/testify/assert"
)

func TestGenreCreation(t *testing.T) {
	name := "Drinks"
	isActive := true
	c := NewGenre(name)

	n := notification.CreateNotification()

	c.Validate(n)

	assert.False(t, n.HasErrors())
	assert.Equal(t, name, c.Name)
	assert.Equal(t, isActive, c.IsActive)
	assert.False(t, c.CreatedAt.IsZero())
	assert.False(t, c.UpdatedAt.IsZero())
	assert.Nil(t, c.DeletedAt)
}

func TestGenreDeactivate(t *testing.T) {
	name := "Drinks"
	c := NewGenre(name)
	updatedAt := c.UpdatedAt

	time.Sleep(1 * time.Millisecond)

	c.Deactivate()

	assert.False(t, c.IsActive)
	assert.True(t, c.UpdatedAt.After(updatedAt))
	assert.NotNil(t, c.DeletedAt)
}

func TestGenreActivate(t *testing.T) {
	name := "Drinks"
	c := NewGenre(name)
	updatedAt := c.UpdatedAt

	c.Deactivate()
	time.Sleep(1 * time.Millisecond)
	c.Activate()

	assert.True(t, c.IsActive)
	assert.True(t, c.UpdatedAt.After(updatedAt))
	assert.Nil(t, c.DeletedAt)
}

func TestGenreUpdateToActive(t *testing.T) {
	name := "new Drinks"
	c := NewGenre("Drinks")
	c.Deactivate()
	updatedAt := c.UpdatedAt

	assert.False(t, c.IsActive)

	time.Sleep(1 * time.Millisecond)

	c.Update(name)
	c.Activate()

	assert.Equal(t, name, c.Name)
	assert.True(t, c.IsActive)
	assert.True(t, c.UpdatedAt.After(updatedAt))
	assert.Nil(t, c.DeletedAt)
}

func TestGenreUpdateToActiveWithInvalidName(t *testing.T) {
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
		c := NewGenre(v.name)

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

func TestGenreUpdateToNotActive(t *testing.T) {
	name := "new Drinks"
	c := NewGenre("Drinks")
	updatedAt := c.UpdatedAt

	assert.True(t, c.IsActive)

	time.Sleep(1 * time.Millisecond)

	c.Update(name)
	c.Deactivate()

	assert.Equal(t, name, c.Name)
	assert.False(t, c.IsActive)
	assert.True(t, c.UpdatedAt.After(updatedAt))
	assert.NotNil(t, c.DeletedAt)
}
