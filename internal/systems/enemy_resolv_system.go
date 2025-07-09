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

// EnemyCollisionWithPlayerSystemFunc handles enemy collision with the player using resolv.
func EnemyResolvMovementSystemFunc(ecs *ecs.ECS) {
	world := ecs.World

	enemyQuery := donburi.NewQuery(filter.Contains(components.Position, components.Velocity, components.Enemy, components.Shape))
	playerQuery := donburi.NewQuery(filter.Contains(components.Player, components.Shape))

	spaceQuery := donburi.NewQuery(filter.Contains(components.Space))
	var spaceEntry *donburi.Entry
	spaceEntry, _ = spaceQuery.FirstEntity(ecs.World)

	if spaceEntry == nil {
		fmt.Println("Error: No Space entity found.")
		return
	}

	space := components.Space.Get(spaceEntry)

	playerEntry, _ := playerQuery.FirstEntity(world)
	if playerEntry == nil {
		return // No player found, nothing to collide with
	}
	playerShape := components.Shape.Get(playerEntry)
	if playerShape == nil {
		fmt.Println("Error: Player entity missing Shape component for EnemyCollisionWithPlayerSystemFunc.")
		return
	}

	enemyQuery.EachEntity(world, func(entry *donburi.Entry) {
		position := components.Position.Get(entry)
		velocity := components.Velocity.Get(entry)
		shape := components.Shape.Get(entry)

		if position == nil || velocity == nil || shape == nil {
			fmt.Println("Error: Enemy entity missing required components for EnemyCollisionWithPlayerSystemFunc.")
			return
		}

		// Calculate potential movement based on the current velocity (set by the enemy AI)
		moveVec := resolv.Vector{X: velocity.X, Y: velocity.Y}

		// Create a test shape for the potential new position
		testShape := shape.Clone()
		testShape.MoveVec(moveVec)

		// Collision test with the Player
		testShape.IntersectionTest(resolv.IntersectionTestSettings{
			TestAgainst: space.FilterShapes().ByTags(dresolv.TagEnemy), // Directly test against the player's shape
			OnIntersect: func(set resolv.IntersectionSet) bool {
				// Move the test shape by the MTV to resolve the collision
				testShape.MoveVec(set.MTV)

				// Obtener el componente X del MTV
				mtvX := set.MTV.X
				mtvY := set.MTV.Y

				// Resolve horizontal velocity
				if velocity.X > 0 && mtvX < 0 {
					velocity.X = 0
				} else if velocity.X < 0 && mtvX > 0 {
					velocity.X = 0
				}

				// Resolve vertical velocity (if you have it)
				if velocity.Y > 0 && mtvY < 0 {
					velocity.Y = 0
				} else if velocity.Y < 0 && mtvY > 0 {
					velocity.Y = 0
				}

				return true
			},
		})

		// Apply the resolved movement
		shape.SetPositionVec(testShape.Position())
		// Use the same offset calculation as the player, assuming similar shape creation
		position.X = testShape.Position().X + float64(shape.Bounds().Width()/2-50)
		position.Y = testShape.Position().Y + float64(shape.Bounds().Height()/2-25)
	})
}
