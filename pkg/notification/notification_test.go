package notification

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotificationConstructor(t *testing.T) {
	noti := CreateNotification()

	assert.NotNil(t, noti)
	assert.Equal(t, 0, len(noti.errors))
	assert.Equal(t, 10, cap(noti.errors))
	assert.False(t, noti.HasErrors())
}

func TestNotificationAdd(t *testing.T) {
	noti := CreateNotification()

	err := errors.New("new error")
	noti.Add(err)

	assert.NotNil(t, noti)
	assert.Equal(t, 1, len(noti.errors))
	assert.Equal(t, err, noti.GetErrors()[0])
	assert.True(t, noti.HasErrors())
	assert.Equal(t, 10, cap(noti.errors))
}
