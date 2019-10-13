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
	state     state
	animation handle
}

func NewGame() game {
	load := newLoader(White, time.Now(), time.Now().Add(3*time.Second/2))
	return game{
		state:     BOOT,
		animation: handle{&load.f, load.update},
	}
}

func (g *game) update(now time.Time) {
	animationDone := g.animation.update(now)
	DisplayLEDs(*g.animation.f)
	if animationDone {
		g.event(ANIMATION_DONE)
	}
}

func (g *game) event(e event) {
	switch e {
	case ANIMATION_DONE:
		g.animationDone()
	case TIMER_2M:
		g.startTimer(2 * time.Minute)
	case TIMER_10M:
		g.startTimer(10 * time.Minute)
	}
}

func (g *game) animationDone() {
	switch g.state {
	case BOOT:
		g.toIdle()
	case COUNTDOWN:
		g.state = TIMERPOP
		pop := newLoader(Red, time.Now(), time.Now().Add(2*time.Second))
		g.animation = handle{&pop.f, pop.update}
	case TIMERPOP:
		g.toIdle()
	}
}

func (g *game) startTimer(d time.Duration) {
	switch g.state {
	case IDLE:
		g.state = COUNTDOWN
		load := newLoader(Black, time.Now(), time.Now().Add(d))
		load.bg = K8SBlue
		g.animation = handle{&load.f, load.update}
	}
}

// toIdle fades from current animation to idle animation
func (g *game) toIdle() {
	g.state = IDLE
	spin := newSpinner()
	fade := newFader(time.Now(), time.Now().Add(1*time.Second))
	fade.from = g.animation
	fade.to = handle{&spin.f, spin.update}
	g.animation = handle{&fade.f, fade.update}
}
