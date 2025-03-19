package systems

import (
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func MovementSystemFunc(ecs *ecs.ECS) {
	world := ecs.World

	query := donburi.NewQuery(filter.Contains(components.Position, components.Velocity))
	query.EachEntity(world, func(entry *donburi.Entry) {
		position := components.Position.Get(entry)
		velocity := components.Velocity.Get(entry)

		position.X += velocity.X
		position.Y += velocity.Y
	})
}
