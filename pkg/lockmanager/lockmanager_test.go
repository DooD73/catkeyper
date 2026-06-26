package lockmanager

import (
	"testing"
	"time"
)

type inputEvent struct {
	keycode   uint16
	flags     uint64
	eventType uint32
	at        time.Duration
}

func shiftDown(keycode uint16) inputEvent {
	return inputEvent{
		keycode:   keycode,
		flags:     ShiftMask,
		eventType: EventFlagsChanged,
	}
}

func shiftUp(keycode uint16) inputEvent {
	return inputEvent{
		keycode:   keycode,
		flags:     0,
		eventType: EventFlagsChanged,
	}
}

func keyDown(keycode uint16) inputEvent {
	return inputEvent{
		keycode:   keycode,
		flags:     ShiftMask,
		eventType: EventKeyDown,
	}
}

func keyDownNoShift(keycode uint16) inputEvent {
	return inputEvent{
		keycode:   keycode,
		flags:     0,
		eventType: EventKeyDown,
	}
}

func keyUp(keycode uint16) inputEvent {
	return inputEvent{
		keycode:   keycode,
		flags:     ShiftMask,
		eventType: EventKeyUp,
	}
}

func simulateEvents(t *testing.T, state *State, events []inputEvent) bool {
	t.Helper()

	base := time.Unix(0, 0)
	state.clock = func() time.Time {
		return base
	}

	var matched bool
	for _, ev := range events {
		if ev.at > 0 {
			base = base.Add(ev.at)
		}
		if state.ProcessEvent(ev.keycode, ev.flags, ev.eventType) {
			matched = true
		}
	}
	return matched
}

func simulateManager(t *testing.T, mgr *Manager, events []inputEvent) (suppress bool, unlocked bool) {
	t.Helper()

	base := time.Unix(0, 0)
	mgr.clock = func() time.Time { return base }
	mgr.state.clock = mgr.clock

	for _, ev := range events {
		if ev.at > 0 {
			base = base.Add(ev.at)
		}
		s, u := mgr.KeyboardDecision(ev.keycode, ev.flags, ev.eventType)
		if s {
			suppress = true
		}
		if u {
			unlocked = true
		}
	}
	return suppress, unlocked
}

func TestProcessEventSequentialChord(t *testing.T) {
	state := NewState()
	events := []inputEvent{
		shiftDown(KeyShift),
		keyDown(KeyC),
		keyDown(KeyA),
		keyDown(KeyT),
	}

	if !simulateEvents(t, state, events) {
		t.Fatal("expected sequential Shift+C+A+T to match")
	}
}

func TestProcessEventSimultaneousChord(t *testing.T) {
	state := NewState()
	events := []inputEvent{
		shiftDown(KeyShift),
		keyDown(KeyT),
		keyDown(KeyA),
		keyDown(KeyC),
	}

	if !simulateEvents(t, state, events) {
		t.Fatal("expected simultaneous Shift+C+A+T to match")
	}
}

func TestProcessEventWrongOrder(t *testing.T) {
	state := NewState()
	events := []inputEvent{
		shiftDown(KeyShift),
		keyDown(KeyC),
		keyDown(KeyT),
	}

	if simulateEvents(t, state, events) {
		t.Fatal("expected C then T without A to not match")
	}
}

func TestProcessEventPartialChord(t *testing.T) {
	tests := []struct {
		name   string
		events []inputEvent
	}{
		{
			name: "C only",
			events: []inputEvent{
				shiftDown(KeyShift),
				keyDown(KeyC),
			},
		},
		{
			name: "CA only",
			events: []inputEvent{
				shiftDown(KeyShift),
				keyDown(KeyC),
				keyDown(KeyA),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := NewState()
			if simulateEvents(t, state, tt.events) {
				t.Fatal("expected partial chord to not match")
			}
		})
	}
}

func TestProcessEventWrongCharMidSequence(t *testing.T) {
	state := NewState()
	events := []inputEvent{
		shiftDown(KeyShift),
		keyDown(KeyC),
		keyDown(42),
		keyDown(KeyA),
	}

	if simulateEvents(t, state, events) {
		t.Fatal("expected wrong character to break sequential progress before T")
	}
}

func TestProcessEventNoShift(t *testing.T) {
	state := NewState()
	events := []inputEvent{
		keyDownNoShift(KeyC),
		keyDownNoShift(KeyA),
		keyDownNoShift(KeyT),
	}

	if simulateEvents(t, state, events) {
		t.Fatal("expected chord without shift to not match")
	}
}

func TestProcessEventShiftReleaseResetsSequence(t *testing.T) {
	state := NewState()
	events := []inputEvent{
		shiftDown(KeyShift),
		keyDown(KeyC),
		shiftUp(KeyShift),
		shiftDown(KeyShift),
		keyDown(KeyA),
		keyDown(KeyT),
	}

	if simulateEvents(t, state, events) {
		t.Fatal("expected shift release to reset sequential progress")
	}
}

func TestProcessEventRightShiftAccepted(t *testing.T) {
	state := NewState()
	events := []inputEvent{
		shiftDown(KeyRShift),
		keyDown(KeyC),
		keyDown(KeyA),
		keyDown(KeyT),
	}

	if !simulateEvents(t, state, events) {
		t.Fatal("expected right shift to work for unlock chord")
	}
}

func TestProcessEventSequentialTimeoutBeforeA(t *testing.T) {
	state := NewState()
	events := []inputEvent{
		shiftDown(KeyShift),
		keyDown(KeyC),
		{keycode: KeyA, flags: ShiftMask, eventType: EventKeyDown, at: 2100 * time.Millisecond},
	}

	if simulateEvents(t, state, events) {
		t.Fatal("expected timeout between C and A to prevent match")
	}
}

func TestProcessEventSequentialTimeoutBeforeT(t *testing.T) {
	state := NewState()
	events := []inputEvent{
		shiftDown(KeyShift),
		keyDown(KeyC),
		keyDown(KeyA),
		keyUp(KeyC),
		keyUp(KeyA),
		{keycode: KeyT, flags: ShiftMask, eventType: EventKeyDown, at: 2100 * time.Millisecond},
	}

	if simulateEvents(t, state, events) {
		t.Fatal("expected timeout between CA and T to prevent match")
	}
}

func TestKeyboardDecisionUnlockedPassthrough(t *testing.T) {
	mgr := New()
	mgr.SetLocked(false)

	suppress, unlocked := mgr.KeyboardDecision(KeyC, ShiftMask, EventKeyDown)
	if suppress || unlocked {
		t.Fatalf("expected passthrough when unlocked, got suppress=%v unlocked=%v", suppress, unlocked)
	}
	if mgr.SuppressedCount() != 0 {
		t.Fatalf("expected suppressed count 0, got %d", mgr.SuppressedCount())
	}
}

func TestKeyboardDecisionLockedSuppresses(t *testing.T) {
	mgr := New()
	mgr.SetLocked(true)

	suppress, unlocked := mgr.KeyboardDecision(42, 0, EventKeyDown)
	if !suppress || unlocked {
		t.Fatalf("expected suppression when locked, got suppress=%v unlocked=%v", suppress, unlocked)
	}
	if mgr.SuppressedCount() != 1 {
		t.Fatalf("expected suppressed count 1, got %d", mgr.SuppressedCount())
	}
}

func TestKeyboardDecisionChordUnlocks(t *testing.T) {
	mgr := New()
	mgr.SetLocked(true)

	events := []inputEvent{
		shiftDown(KeyShift),
		keyDown(KeyC),
		keyDown(KeyA),
		keyDown(KeyT),
	}

	_, unlocked := simulateManager(t, mgr, events)
	if !unlocked {
		t.Fatal("expected chord to unlock keyboard")
	}
	if mgr.IsLocked() {
		t.Fatal("expected manager to be unlocked after chord")
	}
}

func TestSetLockedResetsChordProgress(t *testing.T) {
	mgr := New(WithClock(func() time.Time { return time.Unix(0, 0) }))
	mgr.SetLocked(true)

	events := []inputEvent{
		shiftDown(KeyShift),
		keyDown(KeyC),
		keyDown(KeyA),
	}
	simulateManager(t, mgr, events)

	mgr.SetLocked(true)

	_, unlocked := mgr.KeyboardDecision(KeyT, ShiftMask, EventKeyDown)
	if unlocked {
		t.Fatal("expected partial chord to be reset after SetLocked(true)")
	}
}

func TestEventTypeName(t *testing.T) {
	tests := []struct {
		eventType uint32
		want      string
	}{
		{EventKeyDown, "key_down"},
		{EventKeyUp, "key_up"},
		{EventFlagsChanged, "flags_changed"},
		{0xFFFFFFFE, "tap_disabled_by_timeout"},
		{0xFFFFFFFF, "tap_disabled_by_user_input"},
		{99, "unknown"},
	}

	for _, tt := range tests {
		if got := EventTypeName(tt.eventType); got != tt.want {
			t.Fatalf("EventTypeName(%d) = %q, want %q", tt.eventType, got, tt.want)
		}
	}
}
