package worklogTimer

import (
	tea "github.com/charmbracelet/bubbletea"
	"sync/atomic"
	"time"
)

// WorklogTimer is an extension on stopwatch.Model component that can be used to

var lastID int64

func nextID() int {
	return int(atomic.AddInt64(&lastID, 1))
}

// TickMsg is a message sent on every timer tick.
type TickMsg struct {
	// ID is the identifier of the stopwatch that sends the message. This makes
	// it possible to determine which stopwatch a tick belongs to when there
	// are multiple stopwatches running.
	//
	// Note, however, that a stopwatch will reject ticks from other
	// stopwatches, so it's safe to flow all TickMsgs through all stopwatches
	// and have them still behave appropriately.
	ID  int
	tag int
}

// StartStopMsg is sent when the stopwatch should start or stop.
type StartStopMsg struct {
	ID      int
	running bool
}

// ResetMsg is sent when the stopwatch should reset.
type ResetMsg struct {
	ID int
}

// WorklogTimer for the stopwatch component.
type WorklogTimer struct {
	d       time.Duration
	id      int
	tag     int
	running bool

	// How long to wait before every tick. Defaults to 1 second.
	Interval time.Duration
}

func (w WorklogTimer) SetDuration(d time.Duration) WorklogTimer {
	w.d = d
	return w
}

// NewWithInterval creates a new stopwatch with the given timeout and tick
// interval.
func NewWithInterval(interval time.Duration) WorklogTimer {
	return WorklogTimer{
		Interval: interval,
		id:       nextID(),
	}
}

// New creates a new stopwatch with 1 s interval.
func New() WorklogTimer {
	return NewWithInterval(time.Second)
}

// ID returns the unique ID of the WorklogTimer.
func (w WorklogTimer) ID() int {
	return w.id
}

// Init starts the stopwatch.
func (w WorklogTimer) Init() tea.Cmd {
	return w.Start()
}

// Start starts the stopwatch.
func (w WorklogTimer) Start() tea.Cmd {
	return tea.Sequence(func() tea.Msg {
		return StartStopMsg{ID: w.id, running: true}
	}, tick(w.id, w.tag, w.Interval))
}

// Stop stops the stopwatch.
func (w WorklogTimer) Stop() tea.Cmd {
	return func() tea.Msg {
		return StartStopMsg{ID: w.id, running: false}
	}
}

// Toggle stops the stopwatch if it is running and starts it if it is stopped.
func (w WorklogTimer) Toggle() tea.Cmd {
	if w.Running() {
		return w.Stop()
	}
	return w.Start()
}

// Reset resets the stopwatch to 0.
func (w WorklogTimer) Reset() tea.Cmd {
	return func() tea.Msg {
		return ResetMsg{ID: w.id}
	}
}

// Running returns true if the stopwatch is running or false if it is stopped.
func (w WorklogTimer) Running() bool {
	return w.running
}

// Update handles the timer tick.
func (w WorklogTimer) Update(msg tea.Msg) (WorklogTimer, tea.Cmd) {
	switch msg := msg.(type) {
	case StartStopMsg:
		if msg.ID != w.id {
			return w, nil
		}
		w.running = msg.running
	case ResetMsg:
		if msg.ID != w.id {
			return w, nil
		}
		w.d = 0
	case TickMsg:
		if !w.running || msg.ID != w.id {
			break
		}

		// If a tag is set, and it's not the one we expect, reject the message.
		// This prevents the stopwatch from receiving too many messages and
		// thus ticking too fast.
		if msg.tag > 0 && msg.tag != w.tag {
			return w, nil
		}

		w.d += w.Interval
		w.tag++
		return w, tick(w.id, w.tag, w.Interval)
	}

	return w, nil
}

// Elapsed returns the time elapsed.
func (w WorklogTimer) Elapsed() time.Duration {
	return w.d
}

// View of the timer component.
func (w WorklogTimer) View() string {
	return w.d.String()
}

func tick(id int, tag int, d time.Duration) tea.Cmd {
	return tea.Tick(d, func(_ time.Time) tea.Msg {
		return TickMsg{ID: id, tag: tag}
	})
}
