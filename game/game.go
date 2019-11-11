package game

import (
	"time"

	"github.com/misterikkit/tinytimer/animation"
	"github.com/misterikkit/tinytimer/graphics"
	"github.com/misterikkit/tinytimer/hack"
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

type userInput interface {
	BtnCancel() bool
	Btn2Min() bool
	Btn10Min() bool
}

type Game struct {
	state     State
	Animation animation.Handle
	ui        userInput
}

func New(ui userInput) *Game {
	load := animation.NewLoader(graphics.White, time.Now(), time.Now().Add(hack.ScaleDuration(time.Second/2)))
	return &Game{
		state:     BOOT,
		Animation: animation.Handle{&load.Frame, load.Update},
		ui:        ui,
	}
}

func (g *Game) Update(now time.Time) {
	g.pollInputs()
	animationDone := g.Animation.Update(now)
	if animationDone {
		g.Event(ANIMATION_DONE)
	}
}

func (g *Game) pollInputs() {
	switch {
	case g.ui.BtnCancel():
		g.Event(CANCEL)
	case g.ui.Btn2Min():
		g.Event(TIMER_2M)
	case g.ui.Btn10Min():
		g.Event(TIMER_10M)
	}
}

func (g *Game) Event(e Event) {
	switch e {
	case ANIMATION_DONE:
		g.animationDone()
	case TIMER_2M:
		g.startTimer(hack.ScaleDuration(2 * time.Minute))
	case TIMER_10M:
		g.startTimer(hack.ScaleDuration(10 * time.Second))
		// g.startTimer(hack.ScaleDuration(10 * time.Minute))
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
		pop := animation.NewFlasher(graphics.Red, time.Now().Add(hack.ScaleDuration(2*time.Second)))
		g.Animation = animation.Handle{&pop.Frame, pop.Update}
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
		load := animation.NewLoader(graphics.Black, time.Now(), time.Now().Add(d))
		load.BG = graphics.CSIOrange
		g.Animation = animation.Handle{&load.Frame, load.Update}
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
	fade := animation.NewFader(time.Now(), time.Now().Add(d))
	fade.From = g.Animation
	fade.To = animation.Handle{&spin.Frame, spin.Update}
	g.Animation = animation.Handle{&fade.Frame, fade.Update}
}
