package factory

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/configs"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/archetypes"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components" // Asegúrate de que este path es correcto
	dresolv "github.com/nelsonmarro/dhampyre-blood-vengeance/internal/resolv"
)

func CreateSpace(ecs *ecs.ECS) *donburi.Entry {
	space := archetypes.Space.Spawn(ecs) // Usar archetypes.Space en lugar de archetypes.ResolvSpace

	cfg := configs.C
	resolvSpaceData := resolv.NewSpace(cfg.ScreenWidth, cfg.ScreenHeight, cfg.PlayerSize, cfg.PlayerSize) // Convertir a float64
	components.Space.Set(space, resolvSpaceData)                                                          // Usar components.Space en lugar de components.ResolvSpaceComponent

	tiledBackgroundQuery := donburi.NewQuery(filter.Contains(components.TiledBackground)) // Usar consulta para obtener el componente
	tiledBackgroundEntry, ok := tiledBackgroundQuery.First(ecs.World)
	if ok {
		tiledBackground := components.TiledBackground.Get(tiledBackgroundEntry)
		xCount := configs.C.ScreenWidth / configs.C.TileSize

		for layerIndex, layer := range tiledBackground.Layers {
			for i, tileIndex := range layer {
				tileX := (i % xCount) * configs.C.TileSize
				tileY := (i / xCount) * configs.C.TileSize

				if tileIndex > 0 {
					if layerIndex == 0 { // Primera capa: Fondo principal (Obstáculos)
						if tileIndex == 33 { // Asumiendo que 33 es un obstáculo
							obstacle := resolv.NewRectangle(float64(tileX), float64(tileY), float64(configs.C.TileSize), float64(configs.C.TileSize)) // Convertir a float64
							obstacle.Tags().Set(dresolv.TagObstacle)
							resolvSpaceData.Add(obstacle)
						}
					} else if layerIndex == 1 { // Segunda capa: Plataformas
						tileIndices := []int{
							5, // Índice del tile para las plataformas flotantes - ¡VERIFICAR!
							3, // Índice del tile para las plataformas que se extienden - ¡VERIFICAR!
						}
						for _, platformTileIndex := range tileIndices {
							if tileIndex == platformTileIndex {
								platform := resolv.NewRectangle(float64(tileX), float64(tileY), float64(configs.C.TileSize), float64(configs.C.TileSize)) // Convertir a float64
								platform.Tags().Set(dresolv.TagPlatform)
								resolvSpaceData.Add(platform)
								break // No es necesario seguir verificando otros índices de plataforma para este tile
							}
						}
					}
				}
			}
		}
	}

	return space
}
