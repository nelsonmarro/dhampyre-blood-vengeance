package components

import "github.com/yohamta/donburi"

var PlayerInput = donburi.NewComponentType[PlayerInputComponent]()

type PlayerInputComponent struct {
	MovingLeft  bool
	MovingRight bool
	Jumping     bool
	Attacking   bool
}
