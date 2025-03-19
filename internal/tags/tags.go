package tags

import "github.com/yohamta/donburi"

var (
	PlayerTag  = donburi.NewTag().SetName("Player")
	EnemyTag   = donburi.NewTag().SetName("Enemy")      // Tag para enemigos
	Wall       = donburi.NewTag().SetName("Wall")       // Tag para objetos estáticos (paredes, etc.)
	Obstacle   = donburi.NewTag().SetName("Obstacle")   // Tag para objetos estáticos (paredes, etc.)
	Projectile = donburi.NewTag().SetName("Projectile") // Tag para objetos estáticos (paredes, etc.)
	// PlatformTag        = donburi.NewTag().SetName("Platform")        // Comentando PlatformTag, no lo necesitamos ahora directamente
	// FloatingPlatformTag = donburi.NewTag().SetName("FloatingPlatform") // Comentando FloatingPlatformTag, no lo necesitamos ahora directamente
	// WallTag            = donburi.NewTag().SetName("Wall")            // Comentando WallTag,  podríamos usar StaticObjectTag en su lugar
	// RampTag            = donburi.NewTag().SetName("Ramp")            // Comentando RampTag, no lo necesitamos ahora directamente
)
