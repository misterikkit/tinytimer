package main

import (
	"time"
)

type state uint8
type event uint8

// valid states
const (
	BOOT state = iota
	IDLE
	COUNTDOWN
	TIMERPOP
)

// events
const (
	ANIMATION_DONE event = iota
	TIMER_2M
	TIMER_10M
	CANCEL
)

type game struct {
	state      state
	animation  handle
	pollInputs func()
}

func NewGame() game {
	load := newLoader(White, time.Now(), time.Now().Add(scaleDuration(time.Second/2)))
	return game{
		state:      BOOT,
		animation:  handle{&load.f, load.update},
		pollInputs: func() {},
	}
}

func (g *game) update(now time.Time) {
	g.pollInputs()
	animationDone := g.animation.update(now)
	if animationDone {
		g.event(ANIMATION_DONE)
	}
}

func (g *game) event(e event) {
	switch e {
	case ANIMATION_DONE:
		g.animationDone()
	case TIMER_2M:
		g.startTimer(scaleDuration(2 * time.Minute))
	case TIMER_10M:
		g.startTimer(scaleDuration(10 * time.Second))
		// g.startTimer(scaleDuration(10 * time.Minute))
	case CANCEL:
		g.cancelTimer()
	}
}

func (g *game) animationDone() {
	switch g.state {
	case BOOT:
		g.toIdle(scaleDuration(2 * time.Second))
	case COUNTDOWN:
		g.state = TIMERPOP
		pop := newFlasher(Red, time.Now().Add(scaleDuration(2*time.Second)))
		g.animation = handle{&pop.f, pop.update}
	case TIMERPOP:
		g.toIdle(scaleDuration(1 * time.Second))
	}
}

func (g *game) startTimer(d time.Duration) {
	switch g.state {
	// Allow boot and timer pop animations to be interrupted.
	case BOOT:
		fallthrough
	case TIMERPOP:
		fallthrough
	case IDLE:
		g.state = COUNTDOWN
		// Don't scale this duration because it has been scaled in the caller.
		load := newLoader(Black, time.Now(), time.Now().Add(d))
		load.bg = K8SBlue
		g.animation = handle{&load.f, load.update}
	}
}

func (g *game) cancelTimer() {
	switch g.state {
	case COUNTDOWN:
		g.toIdle(scaleDuration(1 * time.Second))
	}
}

// toIdle fades from current animation to idle animation
func (g *game) toIdle(d time.Duration) {
	g.state = IDLE
	spin := newSpinner()
	// Don't scale this duration because it has been scaled in the caller.
	fade := newFader(time.Now(), time.Now().Add(d))
	fade.from = g.animation
	fade.to = handle{&spin.f, spin.update}
	g.animation = handle{&fade.f, fade.update}
}
