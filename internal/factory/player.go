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

func CreatePlayer(ecs *ecs.ECS, idleAnimation *ganim8.Animation, space *donburi.Entry) *donburi.Entry {
	player := archetypes.Player.Spawn(ecs)

	// Initial position
	initialPosition := &components.PositionComponent{X: float64((configs.C.ScreenWidth / 2) - 300), Y: float64(configs.C.ScreenHeight - (8 * configs.C.TileSize))}
	donburi.Add(player, components.Position, initialPosition)

	// Sprite
	donburi.Add(player, components.Sprite, &components.SpriteComponent{
		Animation:         idleAnimation,
		AnimationName:     "idle",
		AnimationFinished: false,
	})

	// Velocity
	donburi.Add(player, components.Velocity, &components.VelocityComponent{X: 0, Y: 0})

	// Playrt Input
	donburi.Add(player, components.PlayerInput, &components.PlayerInputComponent{})
	// Health
	donburi.Add(player, components.Health, &components.HealthComponent{Current: 100, Max: 100})

	// Magic
	donburi.Add(player, components.Magic, &components.MagicComponent{Current: 50, Max: 100})

	// Player Tag
	components.Player.SetValue(player, components.PlayerComponent{})

	// Create Resolv Shape
	x := float64((configs.C.ScreenWidth / 2) - 300)
	y := float64(configs.C.ScreenHeight - (8 * configs.C.TileSize))

	playerShape := resolv.NewRectangle(x, y, float64(configs.C.PlayerSize-50), float64(configs.C.PlayerSize-60))
	playerShape.Tags().Set(dresolv.TagPlayer)
	dresolv.SetShape(player, playerShape) // Use AsPolygon() to get a ConvexPolygon
	dresolv.Add(space, player)

	return player
}
