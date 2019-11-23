package game

import (
	"time"

	"github.com/misterikkit/tinytimer/animation"
	"github.com/misterikkit/tinytimer/graphics"
	"github.com/misterikkit/tinytimer/hack"
	"github.com/misterikkit/tinytimer/input"
)

type State uint8
type Event uint8

// valid states
const (
	BOOT State = iota
	IDLE
	COUNTDOWN
	TIMERPOP
)

// events
const (
	ANIMATION_DONE Event = iota
	TIMER_2M
	TIMER_10M
	CANCEL
)

// Game keeps track of the current state of the app. Calling it a game makes it more fun.
type Game struct {
	state     State
	Animation animation.Interface
	// TODO: Make animations reusable and store them here to avoid allocations.
}

// New creates a new game in the "boot animation" state.
func New(ui *input.Manager) *Game {
	g := &Game{
		state:     BOOT,
		Animation: animation.NewLoader(graphics.White, graphics.Black, time.Now(), time.Now().Add(hack.ScaleDuration(time.Second/2))),
	}
	ui.AddHandler(input.A_Fall, g.handleInput)
	ui.AddHandler(input.B_Fall, g.handleInput)
	ui.AddHandler(input.C_Fall, g.handleInput)
	return g
}

// Update checks inputs and updates game state & animations based on the current time.
func (g *Game) Update(now time.Time) {
	animationDone := g.Animation.Update(now)
	if animationDone {
		g.Event(ANIMATION_DONE)
	}
}

// handleInput converts input events to game events.
func (g *Game) handleInput(e input.Event) {
	switch e {
	case input.A_Fall:
		g.Event(CANCEL)
	case input.B_Fall:
		g.Event(TIMER_10M)
	case input.C_Fall:
		g.Event(TIMER_2M)
	}
}

// Event signals the game that an event has occurred. These are inputs into the
// state machine.
func (g *Game) Event(e Event) {
	switch e {
	case ANIMATION_DONE:
		g.animationDone()
	case TIMER_2M:
		g.startTimer(hack.ScaleDuration(2 * time.Minute))
	case TIMER_10M:
		// g.startTimer(hack.ScaleDuration(10 * time.Second))
		g.startTimer(hack.ScaleDuration(10 * time.Minute))
	case CANCEL:
		g.cancelTimer()
	}
}

func (g *Game) animationDone() {
	switch g.state {
	case BOOT:
		g.toIdle((2 * time.Second))
	case COUNTDOWN:
		g.state = TIMERPOP
		g.Animation = animation.NewFlasher(graphics.Red, time.Now().Add(hack.ScaleDuration(2*time.Second)))
	case TIMERPOP:
		g.toIdle((1 * time.Second))
	}
}

func (g *Game) startTimer(d time.Duration) {
	switch g.state {
	// Allow boot and timer pop animations to be interrupted.
	case BOOT:
		fallthrough
	case TIMERPOP:
		fallthrough
	case IDLE:
		g.state = COUNTDOWN
		// Don't scale this duration because it has been scaled in the caller.
		g.Animation = animation.NewLoader(graphics.Black, graphics.CSIOrange, time.Now(), time.Now().Add(d))
	}
}

func (g *Game) cancelTimer() {
	switch g.state {
	case COUNTDOWN:
		g.toIdle((1 * time.Second))
	}
}

// toIdle fades from current animation to idle animation
func (g *Game) toIdle(d time.Duration) {
	g.state = IDLE
	spin := animation.NewSpinner(graphics.K8SBlue)
	// Don't scale this duration because it has been scaled in the caller.
	g.Animation = animation.NewFader(g.Animation, spin, time.Now(), time.Now().Add(d))
}
