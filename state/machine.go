package state

const (
	CHOOSING_TICKET = iota
	MEASURING_TIME
	WRITING_WORKLOG
	SENDING
)

type State struct {
	State uint8
	new   bool
}

func NewState() *State {
	return &State{
		new:   true,
		State: CHOOSING_TICKET,
	}
}

func (s *State) Choose() {
	s.State = MEASURING_TIME
	s.new = true
}

func (s *State) Log() {
	s.State = WRITING_WORKLOG
	s.new = true
}

func (s *State) Send() {
	s.State = SENDING
	s.new = true
}

func (s *State) IsNew() bool {
	n := s.new
	s.new = false
	return n
}
