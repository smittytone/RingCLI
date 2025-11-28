package ringcliSpinner

import (
	"fmt"
	"time"
	config "ringcli/lib/config"
)

type Spinner struct {
	isAnimating bool
	cursorIndex int
	cursor      string
	timer       *time.Ticker
	progress    chan bool
	lastPos     int
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

	s.timer = time.NewTicker(150 * time.Millisecond)
	s.progress = make(chan bool)
	go func() {
		for {
			select {
			case <-s.progress:
				// Halted
				s.isAnimating = false
				return
			case <-s.timer.C:
				// Ticker fires
				if config.Config.OutputToText {
					fmt.Printf("\x1B7" + string(s.cursor[s.cursorIndex]) + "\x1B8")
					s.cursorIndex = (s.cursorIndex + 1) % 6
				} else {
					count := 0
					for pos, char := range s.cursor {
						if count == s.cursorIndex {
							fmt.Printf("\x1B7" + string(char) + "\x1B8")
							s.lastPos = pos
							break
						}

						count += 1
					}

					s.cursorIndex = (s.cursorIndex + 1) % 4
				}
			}
		}
	}()
	s.isAnimating = true
}

func (s *Spinner) Stop() {

	if s.isAnimating {
		s.progress <- true
		s.isAnimating = false
	}
}

func (s *Spinner) IsAnimating() bool {

	return s.isAnimating
}
