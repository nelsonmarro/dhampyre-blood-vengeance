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

		// Create a test shape for the potential new position (using the current position)
		testShape := projectileShape.Clone()
		testShape.SetPosition(projectilePosition.X-float64(projectileShape.Bounds().Width()/2-60), projectilePosition.Y-float64(projectileShape.Bounds().Height()/2+60))

		// Check for collisions with Obstacles
		collision := testShape.IntersectionTest(resolv.IntersectionTestSettings{
			TestAgainst: space.FilterShapes().ByTags(dresolv.TagObstacle),
			OnIntersect: func(set resolv.IntersectionSet) bool {
				fmt.Println("Projectile intersected with an obstacle!") // Added log
				// Move the test shape by the MTV to resolve the collision (optional for removal)
				testShape.MoveVec(set.MTV)
				return true // Continue checking for other collisions
			},
		})

		if collision {
			// Projectile collided with an obstacle, so remove the projectile entity
			ecs.World.Remove(entry.Entity())
		}
	})
}
