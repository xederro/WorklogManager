package state

// State machine for managing the state of the application.
//           ┌─────LogWork─────┐
//           │                 ▼
// O─New─►TICKETS           WORKLOG
//           ▲                 │
//           └─────Logged──────┘

const (
	TICKETS = iota
	WORKLOG
)

type State struct {
	state int
	new   bool
}

func New() *State {
	return &State{
		state: TICKETS,
		new:   true,
	}
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
