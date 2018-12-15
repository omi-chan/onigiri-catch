//onigiri.go
package main

import (
	"fmt"
	"math"
	"time"
	"github.com/termbox-go"
)

type state struct {
	Player      []CollisionableMovableObject
	Onigiri     CollisionableMovableObject
	Shadows     []MovableObject
	TopLine     CollisionableObject
	BottomLine  CollisionableObject
	LeftLine    CollisionableObject
	RightLine   CollisionableObject
	Count       int
	End         bool
}

var (
	_temespan = 10
	_height   = 25
	_width    = 80
)

//timer event
func timerLoop(tch chan bool, s state) {
	for {
		tch <- true
		time.Sleep(time.Duration(_temespan) * time.Millisecond)
	}
}

//key events
func keyEventLoop(kch chan termbox.Key) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			kch <- ev.Key
		default:
		}
	}
}

//draw console
func update(s state) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawObj(s.Onigiri)
	for i := range s.Shadows {
		drawObj(s.Shadows[i])
	}
	drawObj(s.LeftLine)
	drawObj(s.RightLine)
	drawObj(s.TopLine)
	drawObj(s.BottomLine)
	drawLine(1, 0, "EXIT : ESC KEY")
	drawLine(20, 0, "← : LEFT")
	drawLine(30, 0, "→ : RIGHT")
	drawLine(41, 0, "↓ : LEFT JUMP")
	drawLine(56, 0, "↑ : RIGHT JUMP")
	for i := range s.Player {
		drawObj(s.Player[i])
	}
	if s.End == true {
		drawObj(NewCollisionableObject(_width/2-_width*3/20, _height*2/5, _width*3/10, 2, "="))
		drawObj(NewCollisionableObject(_width/2-_width*3/20, _height*1/5, _width*3/10, 2, "="))
		drawLine(_width/2-12, _height*8/30, fmt.Sprintf("🍙   G A M E  O V E R  🍙"))
		drawLine(_width/2-8, _height*10/30, fmt.Sprintf("RETRY : ENTER KEY"))
		drawLine(_width/2-7, _height*11/30, fmt.Sprintf("EXIT : ESC KEY"))
	}
	termbox.Flush()
}

//draw object
func drawObj(o Objective) {
	for w := 0; w < o.Size().Width; w++ {
		for h := 0; h < o.Size().Height; h++ {
			termbox.SetCell(o.Point().X+w, o.Point().Y+h,
				[]rune(o.Str())[0], termbox.ColorDefault, termbox.ColorDefault)
		}
	}
}

//drawLine
func drawLine(x, y int, str string) {
	runes := []rune(str)
	for i := 0; i < len(runes); i++ {
		termbox.SetCell(x+i, y, runes[i], termbox.ColorDefault, termbox.ColorDefault)
	}
}

// controller
func controller(s state, kch chan termbox.Key, tch chan bool) {
	var onigiriMaxTime = 9
	var onigiriTime = onigiriMaxTime
	for {
		select {
		case key := <-kch: //key event
			switch key {
			case termbox.KeyEsc, termbox.KeyCtrlC: //game end
				return
			case termbox.KeyEnter: //game retry
				start()
			case termbox.KeyArrowLeft:
				if s.End == false {
				for i := range s.Player {
					s.Player[i].Move(-3, 0)
				}
			}
				break
			case termbox.KeyArrowRight:
				if s.End == false {
				for i := range s.Player {
					s.Player[i].Move(3, 0)
				}
			}
				break
			case termbox.KeyArrowUp:
				var distanse = (_width-_width*1/5-1)-s.Player[0].point.X
				if s.End == false {
				for i := range s.Player {
				s.Player[i].Move(distanse, 0)
				}
			}
				break
			case termbox.KeyArrowDown:
				var distanse = s.Player[0].point.X
				if s.End == false {
				for i := range s.Player {
				s.Player[i].Move(-distanse, 0)
				}
			}
				break
			}
			s = updateStatus(s)
		case <-tch: //time event
			onigiriTime = onigiriTime - 1
			if onigiriTime < 0 {
				onigiriTime = onigiriMaxTime - int(math.Min(float64(s.Count), 4))
				s = onigiriMove(s)
			}
			break
		default:
			break
		}
		update(s)
	}
}

//onigiriMove
func onigiriMove(s state) state {
	s.Onigiri.Next()
	s = updateStatus(s)
	return s
}

//updateStatus
func updateStatus(s state) state {
	for i := range s.Player {
		if s.Onigiri.Collision(s.Player[i]) {
			s.Onigiri.Prev()
			s.Player = append(s.Player, NewCollisionableMovableObject(s.Onigiri.Point().X, s.Onigiri.Point().Y, 2, 1, "🍙", 0, 0))
			s.Onigiri.Stop()
			s.Count++
		}
	}
	if s.Onigiri.Collision(s.BottomLine) {
		s.Onigiri.Prev()
		s.Shadows = nextShadow(s.Shadows, s.Onigiri)
		s.End = true
	}
	return s
}

//initState
func initState() state {

	s := state{}
	_width, _height = termbox.Size()
	s.TopLine = NewCollisionableObject(0, 1, _width, 1, "-")
	s.BottomLine = NewCollisionableObject(0, _height-2, _width, 1, "-")
	s.LeftLine = NewCollisionableObject(0, 0, 1, _height, " ")
	s.RightLine = NewCollisionableObject(_width-1, 0, 1, _height, " ")
	s.Player = append(s.Player,NewCollisionableMovableObject(_width/2-_width*1/20, _height-4, _width*1/5, 2, "-", 0, 0))
	s.Onigiri = inirOnigiri()
	s.Shadows = append(s.Shadows, NewMovableObject(s.Onigiri.Point().X, s.Onigiri.Point().Y, s.Onigiri.Size().Width, s.Onigiri.Size().Height, string(" "), 0, 0))
	s.End = false
	return s
}

//nextShadow
func nextShadow(shadow []MovableObject, onigiri CollisionableMovableObject) []MovableObject {
	for i := len(shadow) - 1; i >= 0; i-- {
		shadow[i] = NewMovableObject(onigiri.Point().X, onigiri.Point().Y, 1, 1, shadow[i].Str(), 0, 0)
	}
	return shadow
}

//inirOnigiri
func inirOnigiri() CollisionableMovableObject {
	return NewCollisionableMovableObject(_width/2, _height/8, 2, 1, "🍙", 0, 1)
}

//start
func start() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	s := initState()

	kch := make(chan termbox.Key)
	tch := make(chan bool)
	go keyEventLoop(kch)
	go timerLoop(tch, s)
	controller(s, kch, tch)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	defer termbox.Close()
}