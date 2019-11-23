package input

type Manager struct {
	handlers            map[Event][]Handler
	btnA, btnB, btnC    func() bool
	lastA, lastB, lastC bool
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

	type multi struct{ ab, ac, bc, abc bool }
	curr := multi{
		ab:  a && b,
		ac:  a && c,
		bc:  b && c,
		abc: a && b && c,
	}
	last := multi{
		ab:  m.lastA && m.lastB,
		ac:  m.lastA && m.lastC,
		bc:  m.lastB && m.lastC,
		abc: m.lastA && m.lastB && m.lastC,
	}

	if curr.ab && !last.ab {
		m.invoke(AB_Rise)
	}
	if !curr.ab && last.ab {
		m.invoke(AB_Fall)
	}
	if curr.ac && !last.ac {
		m.invoke(AC_Rise)
	}
	if !curr.ac && last.ac {
		m.invoke(AC_Fall)
	}
	if curr.bc && !last.bc {
		m.invoke(BC_Rise)
	}
	if !curr.bc && last.bc {
		m.invoke(BC_Fall)
	}
	if curr.abc && !last.abc {
		m.invoke(ABC_Rise)
	}
	if !curr.abc && last.abc {
		m.invoke(ABC_Fall)
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
