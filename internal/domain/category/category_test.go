package category

import (
	"strings"
	"testing"

	"github.com.br/gibranct/admin-do-catalogo/pkg/notification"
	"github.com/stretchr/testify/assert"
)

func TestCategoryCreation(t *testing.T) {
	name := "Drinks"
	desc := "Drinks desc"
	isActive := true
	c := NewCategory(name, desc, isActive)

	n := notification.CreateNotification()

	c.Validate(n)

	assert.False(t, n.HasErrors())
	assert.Equal(t, name, c.Name)
	assert.Equal(t, desc, c.Description)
	assert.Equal(t, isActive, c.IsActive)
	assert.False(t, c.CreatedAt.IsZero())
	assert.False(t, c.UpdatedAt.IsZero())
	assert.True(t, c.DeletedAt.IsZero())
}

func TestCategoryDeactivate(t *testing.T) {
	name := "Drinks"
	desc := "Drinks desc"
	isActive := true
	c := NewCategory(name, desc, isActive)
	updatedAt := c.UpdatedAt

	c.Deactivate()

	assert.False(t, c.IsActive)
	assert.True(t, c.UpdatedAt.After(updatedAt))
	assert.False(t, c.DeletedAt.IsZero())
}

func TestCategoryActivate(t *testing.T) {
	name := "Drinks"
	desc := "Drinks desc"
	isActive := true
	c := NewCategory(name, desc, isActive)
	updatedAt := c.UpdatedAt

	c.Deactivate()
	c.Activate()

	assert.True(t, c.IsActive)
	assert.True(t, c.UpdatedAt.After(updatedAt))
	assert.True(t, c.DeletedAt.IsZero())
}

func TestCategoryUpdateToActive(t *testing.T) {
	name := "new Drinks"
	desc := "Drinks desc"
	isActive := true
	c := NewCategory("Drinks", "desc", false)
	updatedAt := c.UpdatedAt

	assert.False(t, c.IsActive)

	c.Update(name, desc, isActive)

	assert.Equal(t, name, c.Name)
	assert.Equal(t, desc, c.Description)
	assert.True(t, c.IsActive)
	assert.True(t, c.UpdatedAt.After(updatedAt))
	assert.True(t, c.DeletedAt.IsZero())
}

func TestCategoryUpdateToActiveWithInvalidName(t *testing.T) {
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
		c := NewCategory(v.name, "Drinks desc", true)

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

func TestCategoryUpdateToNotActive(t *testing.T) {
	name := "new Drinks"
	desc := "Drinks desc"
	isActive := false
	c := NewCategory("Drinks", "desc", true)
	updatedAt := c.UpdatedAt

	assert.True(t, c.IsActive)

	c.Update(name, desc, isActive)

	assert.Equal(t, name, c.Name)
	assert.Equal(t, desc, c.Description)
	assert.False(t, c.IsActive)
	assert.True(t, c.UpdatedAt.After(updatedAt))
	assert.False(t, c.DeletedAt.IsZero())
}
