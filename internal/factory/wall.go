package factory

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/archetypes"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components"
	dresolv "github.com/nelsonmarro/dhampyre-blood-vengeance/internal/resolv"

	"github.com/solarlune/resolv"
)

func CreateStaticWall(ecs *ecs.ECS, shape *resolv.ConvexPolygon) *donburi.Entry {
	entry := archetypes.Wall.Spawn(ecs) // Usamos el archetype StaticObject
	dresolv.SetShape(entry, shape)      // Adjuntamos la forma

	wallShape := components.Shape.Get(entry)
	if wallShape != nil {
		wallShape.Tags().Set(dresolv.TagSolidWall) // Asignamos la etiqueta aquí
	}
	return entry
}

func CreateObstacle(ecs *ecs.ECS, shape *resolv.ConvexPolygon) *donburi.Entry {
	entry := archetypes.Obstacle.Spawn(ecs) // Usamos el archetype StaticObject
	dresolv.SetShape(entry, shape)          // Adjuntamos la forma

	obstacleShape := components.Shape.Get(entry)
	if obstacleShape != nil {
		obstacleShape.Tags().Set(dresolv.TagObstacle) // Asignamos la etiqueta aquí
	}
	return entry
}
