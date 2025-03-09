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

func StringToState(s string) *State {
	switch s {
	case "available":
		state := Available
		return &state
	case "in_use":
		state := InUse
		return &state
	case "inactive":
		state := Inactive
		return &state
	default:
		return nil
	}
}
