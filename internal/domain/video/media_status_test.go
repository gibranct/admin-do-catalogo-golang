package video

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMediaStatusString(t *testing.T) {
	tests := []struct {
		mediaStatus MediaStatus
		expected    string
	}{
		{
			mediaStatus: PENDING,
			expected:    "PENDING",
		},
		{
			mediaStatus: PROCESSING,
			expected:    "PROCESSING",
		},
		{
			mediaStatus: COMPLETED,
			expected:    "COMPLETED",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.mediaStatus.String())
	}
}
