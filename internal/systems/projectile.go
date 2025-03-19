package systems

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components"
)

func ProjectileSystemFunc(ecs *ecs.ECS) {
	world := ecs.World
	projectileQuery := donburi.NewQuery(filter.Contains(components.Projectile, components.Position, components.Velocity, components.Sprite, components.Shape))

	projectileQuery.Each(world, func(entry *donburi.Entry) {
		position := components.Position.Get(entry)
		velocity := components.Velocity.Get(entry)
		shape := components.Shape.Get(entry)

		if position == nil || velocity == nil || shape == nil {
			return
		}

		// Update projectile position based on its velocity
		position.X += velocity.X
		position.Y += velocity.Y

		// Update the position of the resolv.Shape
		shape.SetPosition(position.X, position.Y)

		// Update the sprite animation
		// You can add logic here to remove the projectile if it goes off-screen
		// or collides with an enemy.
		// if position.X > float64(configs.C.ScreenWidth) || position.X < 0 ||
		//     position.Y > float64(configs.C.ScreenHeight) || position.Y < 0 {
		//         world.Remove(entry.Entity())
		// }
	})
}
