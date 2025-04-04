package ringcliSpinner

import (
	"fmt"
	"time"
)

type Spinner struct {
	isAnimating bool
	cursorIndex int
	cursor      string
	timer       *time.Ticker
	progress    chan bool
}

func NewSpinner(frames string) *Spinner {

	spinner := Spinner{
		isAnimating: false,
		cursorIndex: 0,
		cursor:      frames,
	}

	return &spinner
}

func (s *Spinner) SetCursor(frames string) {

	s.cursor = frames
}

func (s *Spinner) Start() {

	if s.isAnimating {
		return
	}

	s.timer = time.NewTicker(50 * time.Millisecond)
	s.isAnimating = true
	s.progress = make(chan bool)
	go func() {
		for {
			select {
			case <-s.progress:
				s.isAnimating = false
				return
			case <-s.timer.C:
				s.cursorIndex += 1
				if s.cursorIndex >= len(s.cursor) {
					s.cursorIndex = 0
				}

				fmt.Printf("\x1B[1D")
				fmt.Printf(string(s.cursor[s.cursorIndex]))
			}
		}
	}()
}

func (s *Spinner) Stop() {

	if s.isAnimating {
		s.progress <- true
	}
}

func (s *Spinner) IsAnimating() bool {

	return s.isAnimating
}
