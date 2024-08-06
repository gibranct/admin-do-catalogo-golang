package castmember

import "errors"

type CastMemberType uint8

const (
	ACTOR CastMemberType = iota
	DIRECTOR
	UNKNOWN
)

func (c CastMemberType) String() string {
	switch c {
	case ACTOR:
		return "actor"
	case DIRECTOR:
		return "director"

	}
	return "unknown"
}

func TypeFromString(typeStr string) (CastMemberType, error) {
	switch typeStr {
	case "actor":
		return ACTOR, nil
	case "director":
		return DIRECTOR, nil
	default:
		return UNKNOWN, errors.New("unknown type")
	}
}
