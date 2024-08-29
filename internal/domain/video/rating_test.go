package video

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	tests := []struct {
		rating         Rating
		expectedResult string
	}{
		{
			rating:         ER,
			expectedResult: "ER",
		},
		{
			rating:         L,
			expectedResult: "Livre",
		},
		{
			rating:         AGE_10,
			expectedResult: "10",
		},
		{
			rating:         AGE_12,
			expectedResult: "12",
		},
		{
			rating:         AGE_14,
			expectedResult: "14",
		},
		{
			rating:         AGE_16,
			expectedResult: "16",
		},
		{
			rating:         AGE_18,
			expectedResult: "18",
		},
		{
			rating:         UNKNOWN,
			expectedResult: "unknown",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedResult, test.rating.String())
	}
}

func TestStringToRating(t *testing.T) {
	tests := []struct {
		ratingString   string
		expectedResult Rating
		err            error
	}{
		{
			ratingString:   "ER",
			expectedResult: ER,
			err:            nil,
		},
		{
			ratingString:   "Livre",
			expectedResult: L,
			err:            nil,
		},
		{
			ratingString:   "10",
			expectedResult: AGE_10,
			err:            nil,
		},
		{
			ratingString:   "12",
			expectedResult: AGE_12,
			err:            nil,
		},
		{
			ratingString:   "14",
			expectedResult: AGE_14,
			err:            nil,
		},
		{
			ratingString:   "16",
			expectedResult: AGE_16,
			err:            nil,
		},
		{
			ratingString:   "18",
			expectedResult: AGE_18,
			err:            nil,
		},
		{
			ratingString:   "unknown",
			expectedResult: UNKNOWN,
			err:            errors.New("unknown type"),
		},
	}

	for _, test := range tests {
		rate, err := StringToRating(test.ratingString)
		assert.Equal(t, test.err, err)
		assert.Equal(t, test.expectedResult, rate)
	}
}
