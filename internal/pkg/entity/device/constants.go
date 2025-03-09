package device

type (
	State string
)

const (
	Available State = "available"
	InUse     State = "in_use"
	Inactive  State = "inactive"
)

func (s State) String() string {
	return string(s)
}
