package systems

import (
	"fmt"

	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components"
	dresolv "github.com/nelsonmarro/dhampyre-blood-vengeance/internal/resolv"
)

func ProjectileCollisionSystemFunc(ecs *ecs.ECS) {
	projectileQuery := donburi.NewQuery(filter.Contains(components.Projectile, components.Position, components.Shape))
	spaceQuery := donburi.NewQuery(filter.Contains(components.Space))

	var spaceEntry *donburi.Entry
	spaceEntry, _ = spaceQuery.FirstEntity(ecs.World)
	if spaceEntry == nil {
		return // No space to check collisions against
	}
	space := components.Space.Get(spaceEntry)

	projectileQuery.EachEntity(ecs.World, func(entry *donburi.Entry) {
		projectilePosition := components.Position.Get(entry)
		projectileShape := components.Shape.Get(entry)
		if projectilePosition == nil || projectileShape == nil {
			return
		}

		testShape := projectileShape.Clone()
		testShape.Move(projectilePosition.X, projectilePosition.Y)

		// Check for collisions with Obstacles
		collision := testShape.IntersectionTest(resolv.IntersectionTestSettings{
			TestAgainst: space.FilterShapes().ByTags(dresolv.TagObstacle),
			OnIntersect: func(set resolv.IntersectionSet) bool {
				fmt.Println("Projectile intersected with an obstacle!") // Added log
				// Move the test shape by the MTV to resolve the collision
				testShape.MoveVec(set.MTV)
				return true // Continue checking for other collisions
			},
		})

		fmt.Printf("Collision detected: %v\n", collision) // Added log

		if collision {
			// Projectile collided with an obstacle, so remove the projectile entity
			ecs.World.Remove(entry.Entity())
			fmt.Println("Projectile removed due to collision.") // Added log
		}
	})
}
