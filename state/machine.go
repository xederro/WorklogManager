package state

const (
	LOGIN = iota
	TICKETS
	WORKLOG
)

type State struct {
	state int
	new   bool
}

func NewState() *State {
	return &State{
		state: LOGIN,
		new:   true,
	}
}

func (s *State) Login() {
	s.new = true
	s.state = TICKETS
}

func (s *State) LogWork() {
	s.new = true
	s.state = WORKLOG
}

func (s *State) Logged() {
	s.new = false
	s.state = TICKETS
}

func (s *State) GetState() int {
	return s.state
}

func (s *State) IsNew() bool {
	n := s.new
	s.new = false
	return n
}
