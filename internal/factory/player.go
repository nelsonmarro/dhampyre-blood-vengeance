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

func CreatePlayer(ecs *ecs.ECS, idleAnimation *ganim8.Animation, space *donburi.Entry, initialX, initialY float64) *donburi.Entry {
	player := archetypes.Player.Spawn(ecs)

	// Initial position
	initialPosition := &components.PositionComponent{X: initialX, Y: initialY}
	donburi.Add(player, components.Position, initialPosition)

	// Sprite
	donburi.Add(player, components.Sprite, &components.SpriteComponent{
		Animation:     idleAnimation,
		AnimationName: "idle",
	})

	// Velocity
	donburi.Add(player, components.Velocity, &components.VelocityComponent{X: 0, Y: 0})

	// Audio
	donburi.Add(player, components.Audio, &components.AudioComponent{
		Data:       nil,
		Format:     "mp3",
		Loop:       false,
		Playing:    false,
		Volume:     1, // Adjust as needed
		SampleRate: 48000,
	})

	// Playrt Input
	donburi.Add(player, components.PlayerInput, &components.PlayerInputComponent{})
	// Health
	donburi.Add(player, components.Health, &components.HealthComponent{Current: 100, Max: 100})

	// Magic
	donburi.Add(player, components.Magic, &components.MagicComponent{Current: 50, Max: 100})

	// Player Tag
	components.Player.SetValue(player, components.PlayerComponent{})

	playerShape := resolv.NewRectangle(initialPosition.X, initialPosition.Y, float64(configs.C.PlayerSize-40), float64(configs.C.PlayerSize/2-12))
	playerShape.Tags().Set(dresolv.TagPlayer)
	dresolv.SetShape(player, playerShape) // Use AsPolygon() to get a ConvexPolygon
	dresolv.Add(space, player)

	return player
}
