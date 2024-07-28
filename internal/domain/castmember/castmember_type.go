package castmember

type CastMemberType uint8

const (
	ACTOR CastMemberType = iota
	DIRECTOR
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
