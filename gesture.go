/**
Gesture interaction with the device
https://github.com/openatx/uiautomator2#gesture-interaction-with-the-device
*/
package uiautomator

import (
	"fmt"
	"time"
)

type Position struct {
	X float32
	Y float32
}

func (pos *Position) String() string {
	return fmt.Sprintf("%s, %s", pos.X, pos.Y)
}

/*
Convert related position to absolute position
*/
func (ua *UIAutomator) rel2abs(rel *Position) *Position {
	if rel == nil {
		rel = &Position{}
	}

	abs := &Position{
		X: rel.X,
		Y: rel.Y,
	}
	size := &WindowSize{}

	if rel.X < 1 || rel.Y < 1 {
		if ua.size == nil {
			size, _ = ua.GetWindowSize()

			// Cache the window size
			ua.size = size
		}

		size = ua.size
	}

	if rel.X < 1 {
		abs.X = float32(size.Width) * rel.X
	}

	if rel.Y < 1 {
		abs.Y = float32(size.Height) * abs.Y
	}

	return abs
}

/*
Click on the screen
*/
func (ua *UIAutomator) Click(position *Position) error {
	if position.X < 0 || position.Y < 0 {
		return fmt.Errorf("Click: an invalid position %q", position)
	}

	abs := ua.rel2abs(position)

	return ua.post(
		&RPCOptions{
			Method: "click",
			Params: []interface{}{abs.X, abs.Y},
		},
		nil,
		nil,
	)
}

/*
Double click on the screen
*/
func (ua *UIAutomator) DbClick(position *Position, duration float32) error {
	if position.X < 0 || position.Y < 0 {
		return fmt.Errorf("DbClick: an invalid position %q", position)
	}

	abs := ua.rel2abs(position)

	// First click
	if err := ua.Click(abs); err != nil {
		return err
	}

	time.Sleep(time.Duration(duration*1000) * time.Millisecond)

	// Second click
	if err := ua.Click(abs); err != nil {
		return err
	}

	return nil
}

func (ua *UIAutomator) touch(action int, position *Position) error {
	return ua.post(
		&RPCOptions{
			Method: "injectInputEvent",
			Params: []interface{}{action, position.X, position.Y, 0},
		},
		nil,
		nil,
	)
}

func (ua *UIAutomator) touchDown(position *Position) error {
	return ua.touch(0, position)
}

func (ua *UIAutomator) touchUp(position *Position) error {
	return ua.touch(1, position)
}

func (ua *UIAutomator) touchMove(position *Position) error {
	return ua.touch(2, position)
}

/*
Long click on the screen
*/
func (ua *UIAutomator) LongClick(position *Position, duration float32) error {
	if position.X < 0 || position.Y < 0 {
		return fmt.Errorf("LongClick: an invalid position %q", position)
	}

	abs := ua.rel2abs(position)

	// Default duration is 0.5s
	if duration == 0 {
		duration = 0.5
	}

	if err := ua.touchDown(abs); err != nil {
		return err
	}

	time.Sleep(time.Duration(duration*1000) * time.Millisecond)

	if err := ua.touchUp(abs); err != nil {
		return err
	}
	return nil
}

/*
Swipe the screen
*/
func (ua *UIAutomator) Swipe(from *Position, to *Position, step int) error {
	if from.X < 0 || from.Y < 0 || to.X < 0 || to.Y < 0 {
		return fmt.Errorf("Swipe: invalid from(%s) -> to(%s)", from, to)
	}

	from = ua.rel2abs(from)
	to = ua.rel2abs(to)

	return ua.post(
		&RPCOptions{
			Method: "swipe",
			Params: []interface{}{from.X, from.Y, to.X, to.Y, step},
		},
		nil,
		nil,
	)
}

/*
Swipe by points, unlock the gesture login
*/
func (ua *UIAutomator) SwipePoints(points ...*Position) error {
	var positions []int

	for _, v := range points {
		abs := ua.rel2abs(v)
		positions = append(positions, int(abs.X), int(abs.Y))
	}

	return ua.post(
		&RPCOptions{
			Method: "swipePoints",
			Params: []interface{}{positions, 20},
		},
		nil,
		nil,
	)
}

/*
Swipe the screen
*/
func (ua *UIAutomator) Drag(start *Position, end *Position, duration float32) error {
	if start.X < 0 || start.Y < 0 || end.X < 0 || end.Y < 0 {
		return fmt.Errorf("Drag: invalid start(%s) -> end(%s)", start, end)
	}

	start = ua.rel2abs(start)
	end = ua.rel2abs(end)

	return ua.post(
		&RPCOptions{
			Method: "drag",
			Params: []interface{}{start.X, start.Y, end.X, end.Y, duration * 200},
		},
		nil,
		nil,
	)
}
