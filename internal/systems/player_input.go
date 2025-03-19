package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components"
)

func PlayerInputSystemFunc(ecs *ecs.ECS) {
	world := ecs.World

	playerQuery := donburi.NewQuery(filter.Contains(components.Player, components.PlayerInput))

	playerQuery.Each(world, func(entry *donburi.Entry) {
		input := components.PlayerInput.Get(entry)

		input.MovingLeft = ebiten.IsKeyPressed(ebiten.KeyArrowLeft)
		input.MovingRight = ebiten.IsKeyPressed(ebiten.KeyArrowRight)
		input.Jumping = ebiten.IsKeyPressed(ebiten.KeySpace)
		input.Attacking = ebiten.IsKeyPressed(ebiten.KeyA)
	})
}
