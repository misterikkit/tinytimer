package input

type multi struct{ ab, ac, bc, abc bool }

type Manager struct {
	handlers            map[Event][]Handler
	btnA, btnB, btnC    func() bool
	lastA, lastB, lastC bool
	lastMulti           multi
	done                chan struct{}
}

type Handler = func(Event)

// Event is any input behavior that can be subscribed to.
type Event uint8

// Available events
const (
	A_Rise Event = iota
	A_Fall
	B_Rise
	B_Fall
	C_Rise
	C_Fall
	AB_Rise
	AB_Fall
	AC_Rise
	AC_Fall
	BC_Rise
	BC_Fall
	ABC_Rise
	ABC_Fall
)

// NewManager initializes a new input manager.
func NewManager(a, b, c func() bool) *Manager {
	return &Manager{
		handlers: make(map[Event][]Handler),
		btnA:     a,
		btnB:     b,
		btnC:     c,
	}
}

// Poll generates events based on current & past state, then updates internal state.
func (m *Manager) Poll() {
	a, b, c := m.btnA(), m.btnB(), m.btnC()
	if a && !m.lastA {
		m.invoke(A_Rise)
	}
	if !a && m.lastA {
		m.invoke(A_Fall)
	}
	if b && !m.lastB {
		m.invoke(B_Rise)
	}
	if !b && m.lastB {
		m.invoke(B_Fall)
	}
	if c && !m.lastC {
		m.invoke(C_Rise)
	}
	if !c && m.lastC {
		m.invoke(C_Fall)
	}

	curr := multi{
		ab:  a && b,
		ac:  a && c,
		bc:  b && c,
		abc: a && b && c,
	}

	// Only set lastMulti values when the edge is detected, rather than just
	// assigning curr to them. This is because the rising edge condition is not the
	// inverse of the falling edge condition.
	if curr.ab && !m.lastMulti.ab {
		m.invoke(AB_Rise)
		m.lastMulti.ab = true
	}
	if curr.ac && !m.lastMulti.ac {
		m.invoke(AC_Rise)
		m.lastMulti.ac = true
	}
	if curr.bc && !m.lastMulti.bc {
		m.invoke(BC_Rise)
		m.lastMulti.bc = true
	}
	if curr.abc && !m.lastMulti.abc {
		m.invoke(ABC_Rise)
		m.lastMulti.abc = true
	}
	// Falling multi-edges are a special case. They don't trigger when the combo
	// ceases to be high, but when all of the combo's buttons are low.
	if m.lastMulti.ab && !a && !b {
		m.invoke(AB_Fall)
		m.lastMulti.ab = false
	}
	if m.lastMulti.ac && !a && !c {
		m.invoke(AC_Fall)
		m.lastMulti.ac = false
	}
	if m.lastMulti.bc && !b && !c {
		m.invoke(BC_Fall)
		m.lastMulti.bc = false
	}
	if m.lastMulti.abc && !a && !b && !c {
		m.invoke(ABC_Fall)
		m.lastMulti.abc = false
	}

	// TODO: see if a defer will work here
	m.lastA, m.lastB, m.lastC = a, b, c
}

// AddHandler registers a callback for the given event.
func (m *Manager) AddHandler(e Event, h Handler) {
	m.handlers[e] = append(m.handlers[e], h)
}

// invoke calls all registered handlers for the given event.
func (m *Manager) invoke(e Event) {
	for _, f := range m.handlers[e] {
		f(e)
	}
}
