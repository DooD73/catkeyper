package lockmanager

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	KeyA      uint16 = 0
	KeyC      uint16 = 8
	KeyT      uint16 = 17
	KeyShift  uint16 = 56
	KeyRShift uint16 = 60

	EventKeyDown      uint32 = 10
	EventKeyUp        uint32 = 11
	EventFlagsChanged uint32 = 12

	ShiftMask uint64 = 1 << 17
)

const sequenceWindow = 2 * time.Second

type clockFunc func() time.Time

type ManagerOption func(*Manager)

func WithClock(clock clockFunc) ManagerOption {
	return func(m *Manager) {
		m.clock = clock
	}
}

type Manager struct {
	locked           atomic.Bool
	suppressedEvents atomic.Uint64
	state            *State
	clock            clockFunc
}

func New(opts ...ManagerOption) *Manager {
	m := &Manager{
		state: NewState(),
		clock: time.Now,
	}
	for _, opt := range opts {
		opt(m)
	}
	m.state.clock = m.clock
	return m
}

func (m *Manager) IsLocked() bool {
	return m.locked.Load()
}

func (m *Manager) SetLocked(locked bool) {
	m.locked.Store(locked)
	m.state.Reset()
}

func (m *Manager) SuppressedCount() uint64 {
	return m.suppressedEvents.Load()
}

func (m *Manager) KeyboardDecision(keycode uint16, flags uint64, eventType uint32) (suppress bool, unlocked bool) {
	if !m.locked.Load() {
		return false, false
	}

	if m.state.ProcessEvent(keycode, flags, eventType) {
		m.locked.Store(false)
		unlocked = true
	}

	m.suppressedEvents.Add(1)
	return true, unlocked
}

type State struct {
	mu          sync.Mutex
	pressed     map[uint16]bool
	sequence    string
	sequenceTTL time.Time
	clock       clockFunc
}

func NewState() *State {
	return &State{
		pressed: make(map[uint16]bool),
		clock:   time.Now,
	}
}

func (s *State) ProcessEvent(keycode uint16, flags uint64, eventType uint32) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.clock()
	shiftHeld := flags&ShiftMask != 0 || s.pressed[KeyShift] || s.pressed[KeyRShift]

	switch eventType {
	case EventFlagsChanged:
		if keycode == KeyShift || keycode == KeyRShift {
			s.pressed[keycode] = flags&ShiftMask != 0
			if !s.pressed[KeyShift] && !s.pressed[KeyRShift] && flags&ShiftMask == 0 {
				s.resetLocked()
			}
		}
		return false
	case EventKeyUp:
		delete(s.pressed, keycode)
		return false
	case EventKeyDown:
		s.pressed[keycode] = true
	default:
		return false
	}

	if !shiftHeld {
		s.resetLocked()
		return false
	}

	if s.pressed[KeyC] && s.pressed[KeyA] && s.pressed[KeyT] {
		s.resetLocked()
		return true
	}

	if now.After(s.sequenceTTL) {
		s.sequence = ""
	}
	s.sequenceTTL = now.Add(sequenceWindow)

	switch keycode {
	case KeyC:
		s.sequence = "C"
	case KeyA:
		if s.sequence == "C" {
			s.sequence = "CA"
		} else {
			s.sequence = ""
		}
	case KeyT:
		if s.sequence == "CA" {
			s.resetLocked()
			return true
		}
		s.sequence = ""
	default:
		s.sequence = ""
	}

	return false
}

func (s *State) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resetLocked()
}

func (s *State) resetLocked() {
	s.sequence = ""
	s.sequenceTTL = time.Time{}
	for k := range s.pressed {
		if k != KeyShift && k != KeyRShift {
			delete(s.pressed, k)
		}
	}
}

// EventTypeName maps macOS CGEvent type values to stable names for logging.
func EventTypeName(eventType uint32) string {
	switch eventType {
	case EventKeyDown:
		return "key_down"
	case EventKeyUp:
		return "key_up"
	case EventFlagsChanged:
		return "flags_changed"
	case 0xFFFFFFFE:
		return "tap_disabled_by_timeout"
	case 0xFFFFFFFF:
		return "tap_disabled_by_user_input"
	default:
		return "unknown"
	}
}
