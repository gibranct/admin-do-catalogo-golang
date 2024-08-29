package video

import "errors"

type Rating uint8

const (
	ER Rating = iota
	L
	AGE_10
	AGE_12
	AGE_14
	AGE_16
	AGE_18
	UNKNOWN
)

func (r Rating) String() string {
	switch r {
	case ER:
		return "ER"
	case L:
		return "Livre"
	case AGE_10:
		return "10"
	case AGE_12:
		return "12"
	case AGE_14:
		return "14"
	case AGE_16:
		return "16"
	case AGE_18:
		return "18"

	}
	return "unknown"
}

func StringToRating(typeStr string) (Rating, error) {
	switch typeStr {
	case "ER":
		return ER, nil
	case "Livre":
		return L, nil
	case "10":
		return AGE_10, nil
	case "12":
		return AGE_12, nil
	case "14":
		return AGE_14, nil
	case "16":
		return AGE_16, nil
	case "18":
		return AGE_18, nil
	default:
		return UNKNOWN, errors.New("unknown type")
	}
}
