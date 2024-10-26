package genre

import (
	"strings"
	"testing"
	"time"

	"github.com.br/gibranct/admin_do_catalogo/pkg/notification"
	"github.com/stretchr/testify/assert"
)

func TestGenreCreation(t *testing.T) {
	name := "Drinks"
	isActive := true
	g := NewGenre(name)

	n := notification.CreateNotification()

	g.Validate(n)

	assert.False(t, n.HasErrors())
	assert.Equal(t, name, g.Name)
	assert.Equal(t, isActive, g.IsActive)
	assert.False(t, g.CreatedAt.IsZero())
	assert.False(t, g.UpdatedAt.IsZero())
	assert.True(t, len(g.CategoryIds) == 0)
	assert.True(t, cap(g.CategoryIds) == 5)
	assert.Nil(t, g.DeletedAt)
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

func TestAddCategoryId(t *testing.T) {
	cId1 := int64(89)
	cId2 := int64(63)
	cId3 := int64(23)
	cId4 := int64(45)
	cId5 := int64(3)
	cId6 := int64(2)
	g := NewGenre("genre 1")

	assert.Equal(t, len(g.CategoryIds), 0)
	assert.Equal(t, cap(g.CategoryIds), 5)

	g.AddCategoryId(cId1)
	g.AddCategoryId(cId2)
	g.AddCategoryId(cId3)
	g.AddCategoryId(cId4)
	g.AddCategoryId(cId5)
	g.AddCategoryId(cId6)

	assert.Equal(t, len(g.CategoryIds), 6)
	assert.Equal(t, cap(g.CategoryIds), 10)
}

func TestAddCategoryIds(t *testing.T) {
	cId1 := int64(89)
	cId2 := int64(63)
	cId3 := int64(23)
	cId4 := int64(45)
	cId5 := int64(3)
	cId6 := int64(2)
	ids := []int64{cId1, cId2, cId3, cId4, cId5, cId6}
	g := NewGenre("genre 1")

	assert.Equal(t, len(g.CategoryIds), 0)
	assert.Equal(t, cap(g.CategoryIds), 5)

	g.AddCategoriesIds(ids)

	assert.Equal(t, len(g.CategoryIds), 6)
	assert.Equal(t, cap(g.CategoryIds), 10)
}

func TestRemoveCategoryId(t *testing.T) {
	cId1 := int64(89)
	cId2 := int64(63)
	cId3 := int64(23)
	cId4 := int64(45)
	cId5 := int64(3)
	cId6 := int64(2)
	ids := []int64{cId1, cId2, cId3, cId4, cId5, cId6}
	g := NewGenre("genre 1")

	assert.Equal(t, len(g.CategoryIds), 0)
	assert.Equal(t, cap(g.CategoryIds), 5)

	g.AddCategoriesIds(ids)

	assert.Equal(t, len(g.CategoryIds), 6)
	assert.Equal(t, cap(g.CategoryIds), 10)

	counter := len(g.CategoryIds)

	for _, cId := range g.CategoryIds {
		g.RemoveCategoryId(cId)
		counter--
		assert.Equal(t, len(g.CategoryIds), counter)
	}

	assert.Empty(t, g.CategoryIds)
}
