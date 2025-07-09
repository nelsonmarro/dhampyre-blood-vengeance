package factory

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/ganim8/v2"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/archetypes"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components"
	dresolv "github.com/nelsonmarro/dhampyre-blood-vengeance/internal/resolv"
)

func CreateProjectile(ecs *ecs.ECS, animation *ganim8.Animation, owner *donburi.Entry, startX, startY, velocityX, velocityY float64, facingLeft bool, spaceEntry *donburi.Entry) *donburi.Entry {
	projectile := archetypes.Projectile.Spawn(ecs)

	initialPosition := &components.PositionComponent{X: startX, Y: startY}
	donburi.Add(projectile, components.Position, initialPosition)

	donburi.Add(projectile, components.Velocity, &components.VelocityComponent{X: velocityX, Y: velocityY})

	projectileOptions := &components.ProjectileComponent{Damage: 10, Owner: owner}
	donburi.Add(projectile, components.Projectile, projectileOptions)

	donburi.Add(projectile, components.Sprite, &components.SpriteComponent{
		Animation:      animation,
		AnimationName:  "fire",
		AnimationSpeed: 0.1,
		Flipped:        facingLeft,
	})

	donburi.Add(projectile, components.Audio, &components.AudioComponent{
		Data:       nil,
		Format:     "mp3",
		Loop:       false,
		Playing:    false,
		Volume:     0.8, // Adjust as needed
		SampleRate: 48000,
	})

	projectileWidth := float64(animation.Sprite().W())
	projectileHeight := float64(animation.Sprite().H())

	// Calculate top-left corner for the resolv shape
	shapeX := startX - projectileWidth/2
	shapeY := startY - projectileHeight/2

	obstacleShape := resolv.NewRectangle(shapeX, shapeY, projectileWidth, projectileHeight)
	obstacleShape.SetPosition(shapeX, shapeY)
	obstacleShape.Tags().Set(dresolv.TagProjectile)
	dresolv.SetShape(projectile, obstacleShape)

	space := components.Space.Get(spaceEntry)
	if space != nil {
		space.Add(obstacleShape) // Add the shape to the resolv space
	}

	return projectile
}
