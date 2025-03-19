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

// ResolvMovementSystemFunc es el sistema encargado del movimiento del jugador usando resolv para colisiones.
func ResolvMovementSystemFunc(ecs *ecs.ECS) {
	playerQuery := donburi.NewQuery(filter.Contains(components.Player))
	var playerEntry *donburi.Entry
	playerEntry, _ = playerQuery.FirstEntity(ecs.World)

	if playerEntry == nil {
		fmt.Println("Error: No Player entity found.")
		return
	}

	spaceQuery := donburi.NewQuery(filter.Contains(components.Space))
	var spaceEntry *donburi.Entry
	spaceEntry, _ = spaceQuery.FirstEntity(ecs.World)

	if spaceEntry == nil {
		fmt.Println("Error: No Space entity found.")
		return
	}

	space := components.Space.Get(spaceEntry)
	position := components.Position.Get(playerEntry)
	velocity := components.Velocity.Get(playerEntry)
	shape := components.Shape.Get(playerEntry)

	if position == nil || velocity == nil || shape == nil {
		fmt.Println("Error: Player entity missing required components for ResolvMovementSystem.")
		return
	}

	// Calculate potential movement
	moveVec := resolv.Vector{X: velocity.X, Y: velocity.YSpeed}
	// Reset OnGround at the beginning of each frame
	velocity.OnGround = false

	// Perform ground detection
	checkVec := resolv.Vector{X: 0, Y: 5} // Aumentamos un poco para mayor seguridad
	isOnGround := shape.ShapeLineTest(resolv.ShapeLineTestSettings{
		Vector:      checkVec,
		TestAgainst: space.FilterShapes().ByTags(dresolv.TagSolidWall),
	})
	if isOnGround {
		velocity.OnGround = true
		if velocity.YSpeed > 0 {
			velocity.YSpeed = 0
		}
	}

	// Create a test shape for the potential new position
	testShape := shape.Clone()
	testShape.MoveVec(moveVec)

	// Collision test with Solid Walls
	testShape.IntersectionTest(resolv.IntersectionTestSettings{
		TestAgainst: space.FilterShapes().ByTags(dresolv.TagSolidWall),
		OnIntersect: func(set resolv.IntersectionSet) bool {
			// Move the test shape by the MTV to resolve the collision
			testShape.MoveVec(set.MTV)

			// Obtener el componente X del MTV
			mtvX := set.MTV.X
			mtvY := set.MTV.Y

			// Resolver la velocidad horizontal
			if velocity.X > 0 && mtvX < 0 {
				velocity.X = 0
			} else if velocity.X < 0 && mtvX > 0 {
				velocity.X = 0
			}

			// Resolver la velocidad vertical
			if velocity.YSpeed < 0 && mtvY > 0 { // Hitting something above
				velocity.YSpeed = 0
			}

			return true // Continue checking for other collisions
		},
	})

	// Collision test with Obstacles
	testShape.IntersectionTest(resolv.IntersectionTestSettings{
		TestAgainst: space.FilterShapes().ByTags(dresolv.TagObstacle),
		OnIntersect: func(set resolv.IntersectionSet) bool {
			// Move the test shape by the MTV to resolve the collision
			testShape.MoveVec(set.MTV)

			// Obtener el componente X del MTV
			mtvX := set.MTV.X
			mtvY := set.MTV.Y

			// Resolver la velocidad horizontal
			if velocity.X > 0 && mtvX < 0 {
				velocity.X = 0
			} else if velocity.X < 0 && mtvX > 0 {
				velocity.X = 0
			}

			// Resolver la velocidad vertical
			if velocity.YSpeed < 0 && mtvY > 0 { // Hitting something above
				velocity.YSpeed = 0
			} else if velocity.YSpeed > 0 && mtvY < 0 { // Landing on an obstacle
				velocity.YSpeed = 0
				velocity.OnGround = true
			}

			return true // Continue checking for other collisions
		},
	})

	// Apply the resolved movement
	shape.SetPositionVec(testShape.Position())
	position.X = testShape.Position().X + float64(shape.Bounds().Width()/2)
	position.Y = testShape.Position().Y + float64(shape.Bounds().Height()/2)
}
