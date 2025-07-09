package archetypes

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/layers"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/tags"
)

var (
	// Player Archetype (con Tag)
	Player = newArchetype(
		tags.PlayerTag,
		components.Player,
		components.Position,
		components.Sprite,
		components.Velocity,
		components.Health,
		components.Magic,
		components.Shape, // Componente para resolv.Shape (renombrado desde ObjectComponent)
	)

	BackgroundAudio = newArchetype(
		components.Audio,
	)

	// Player Archetype (con Tag)
	Enemy = newArchetype(
		tags.EnemyTag,
		components.Enemy,
		components.Position,
		components.Sprite,
		components.Velocity,
		components.Health,
		components.Shape, // Componente para resolv.Shape (renombrado desde ObjectComponent)
	)

	Projectile = newArchetype(
		tags.Projectile,
		components.Position,
		components.Sprite,
		components.Velocity,
		components.Projectile,
		components.Shape, // Componente para resolv.Shape (renombrado desde ObjectComponent)
	)

	Wall = newArchetype(
		tags.Wall,
		components.Shape,
		components.Position, // ¡Añadimos el componente Position!
		// components.Sprite, // Omitiendo Sprite por ahora en StaticObject
	)

	Obstacle = newArchetype(
		tags.Obstacle,
		components.Shape, // ¡Añadimos el componente Position!
		// components.Sprite, // Omitiendo Sprite por ahora en StaticObject
	)

	// Space Archetype
	Space = newArchetype(
		components.Space,
	)
)

type archetype struct {
	components []donburi.IComponentType
}

func newArchetype(cs ...donburi.IComponentType) *archetype {
	return &archetype{
		components: cs,
	}
}

func (a *archetype) Spawn(ecs *ecs.ECS, cs ...donburi.IComponentType) *donburi.Entry {
	e := ecs.World.Entry(ecs.Create(
		layers.Default, // Usando layers.Default para la capa por defecto
		append(a.components, cs...)...,
	))
	return e
}
