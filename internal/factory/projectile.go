package factory

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/ganim8/v2"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/configs"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/archetypes"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components"
	dresolv "github.com/nelsonmarro/dhampyre-blood-vengeance/internal/resolv"
)

func CreateProjectile(ecs *ecs.ECS, animation *ganim8.Animation, owner *donburi.Entry, startX, startY, velocityX, velocityY float64, facingLeft bool, space *donburi.Entry) *donburi.Entry {
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

	if startX > float64(configs.C.ScreenWidth/2) {
		startX = startX/2 - 80
	}
	if startY > float64(configs.C.ScreenHeight/2) {
		startY = startY/2 - 80
	}
	obstacleShape := resolv.NewRectangleFromCorners(startX, startY, float64((animation.Sprite().W())), float64(animation.Sprite().H())) // Ajusta el tamaño
	obstacleShape.Tags().Set(dresolv.TagProjectile)
	dresolv.SetShape(projectile, obstacleShape)
	dresolv.Add(space, projectile)

	return projectile
}
