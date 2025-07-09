package factory

import (
	_ "image/png"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/configs"
	"github.com/solarlune/resolv"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/ganim8/v2"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/archetypes"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components" // Corrected import path for components
	dresolv "github.com/nelsonmarro/dhampyre-blood-vengeance/internal/resolv"
)

func CreateEnemySkeleton(ecs *ecs.ECS, position *components.PositionComponent, idleAnimation *ganim8.Animation, space *donburi.Entry) *donburi.Entry {
	enemy := archetypes.Enemy.Spawn(ecs)

	// Initial position
	initialPosition := position
	donburi.Add(enemy, components.Position, initialPosition)

	donburi.Add(enemy, components.Sprite, &components.SpriteComponent{
		Animation:     idleAnimation,
		AnimationName: "eidle",
	})

	donburi.Add(enemy, components.Velocity, &components.VelocityComponent{X: 0, Y: 0})

	donburi.Add(enemy, components.Health, &components.HealthComponent{Current: 100, Max: 100})

	donburi.Add(enemy, components.Enemy, &components.EnemyComponent{
		CurrentState:    components.EnemyStateIdle, // Initialize the state
		InitialPosition: initialPosition.X,
	})

	enemyShape := resolv.NewRectangle(initialPosition.X, initialPosition.Y, float64(configs.C.PlayerSize-40), float64(configs.C.PlayerSize/2-12))
	enemyShape.Tags().Set(dresolv.TagEnemy)
	dresolv.SetShape(enemy, enemyShape) // Use AsPolygon() to get a ConvexPolygon
	dresolv.Add(space, enemy)

	return enemy
}
