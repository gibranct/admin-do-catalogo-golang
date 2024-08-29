package video

type MediaStatus uint8

const (
	PENDING MediaStatus = iota
	PROCESSING
	COMPLETED
)

func (ms MediaStatus) String() string {
	switch ms {
	case PENDING:
		return "PENDING"
	case PROCESSING:
		return "PROCESSING"
	case COMPLETED:
		return "COMPLETED"
	}
	return "unknown"
}
