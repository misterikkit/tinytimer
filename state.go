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
)

// events
const (
	ANIMATION_DONE event = iota
)

type game struct {
	state     state
	animation handle
}

func NewGame() game {
	load := newLoader(K8SBlue, time.Now(), time.Now().Add(20*time.Second))
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
	}
}

func (g *game) animationDone() {
	switch g.state {
	case BOOT:
		g.state = IDLE
		spin := newSpinner()
		fade := newFader(time.Now(), time.Now().Add(1*time.Second))
		fade.from = g.animation
		fade.to = handle{&spin.f, spin.update}
		g.animation = handle{&fade.f, fade.update}
	}
}
