package tiled

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/configs"
)

type TiledBackground struct {
	TilesImage *ebiten.Image // Para el tile set gótico
	Layers     [][]int
}

type Platform struct {
	X    int // Coordenada X inicial en tiles
	Y    int // Coordenada Y inicial en tiles
	W    int // Ancho en tiles
	H    int // Alto en tiles (para plataformas que se extienden)
	Type int // 0 para plataforma flotante, 1 para plataforma que se extiende
}

const (
	PlatformTypeFloating        = 0
	PlatformTypeGroundExtending = 1
)

func NewTiledBackground() *TiledBackground {
	return &TiledBackground{
		Layers: [][]int{
			generateMainTiledLayer(),
			generatePlatformLayer(),
		},
	}
}

func (t *TiledBackground) DrawTileBackground(screen *ebiten.Image) {
	if t.TilesImage == nil || len(t.Layers) == 0 {
		return
	}

	tileWidth := configs.C.TileSize
	tileHeight := configs.C.TileSize
	tileXCount := t.TilesImage.Bounds().Dx() / tileWidth
	xCount := configs.C.ScreenWidth / tileWidth

	for _, layer := range t.Layers {
		for i, tileIndex := range layer {
			if tileIndex <= 0 {
				continue
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64((i%xCount)*tileWidth), float64((i/xCount)*tileHeight))

			sx := (tileIndex % tileXCount) * tileWidth
			sy := (tileIndex / tileXCount) * tileHeight
			screen.DrawImage(t.TilesImage.SubImage(image.Rect(sx, sy, sx+tileWidth, sy+tileHeight)).(*ebiten.Image), op)
		}
	}
}

func generateMainTiledLayer() []int {
	xCount := configs.C.ScreenWidth / configs.C.TileSize  // numero de filas donde dibujar los azulejos
	yCount := configs.C.ScreenHeight / configs.C.TileSize // numero de columnas donde dibujar los azulejos
	layer := make([]int, xCount*yCount)

	tileIndices := []int{ // ÍNDICES DE AZULEJOS DEFINIDOS POR EL USUARIO
		1,   // Suelo de piedra (índice 1)
		33,  // Paredes de piedra (índice 33)
		100, // Vacio (índice 100)
	}

	index := 0
	for y := 0; y < yCount; y++ {
		for x := 0; x < xCount; x++ {
			tileType := 2 // Por defecto, usar el índice de "vacio" (tileIndices[2] = 100)

			if y < 3 { // Las primeras 3 filas desde arriba
				tileType = 1 // Usar el índice de "paredes de piedra" (tileIndices[1] = 33)
			} else if y == 3 { // La cuarta fila desde arriba
				tileType = 0 // Usar el índice de "suelo de piedra" (tileIndices[0] = 1)
			} else if y >= yCount-4 { // Las últimas 4 filas desde abajo (últimas 4 incluyendo la fila yCount-4)
				if y == yCount-4 { // La cuarta fila desde abajo
					tileType = 0 // Usar el índice de "suelo de piedra" (tileIndices[0] = 1)
				} else { // Las últimas 3 filas desde abajo
					tileType = 1 // Usar el índice de "paredes de piedra" (tileIndices[1] = 33)
				}
			}

			layer[index] = tileIndices[tileType]
			index++
		}
	}

	return layer
}

func definePlatforms() []Platform {
	return []Platform{
		{X: 17, Y: 15, W: 5, H: 1, Type: PlatformTypeFloating},        // Plataforma flotante
		{X: 25, Y: 20, W: 4, H: 6, Type: PlatformTypeGroundExtending}, // Otra plataforma que se extiende
		// Puedes añadir más plataformas aquí
	}
}

func generatePlatformLayer() []int {
	xCount := configs.C.ScreenWidth / configs.C.TileSize
	yCount := configs.C.ScreenHeight / configs.C.TileSize
	layer := make([]int, xCount*yCount)

	tileIndices := []int{
		8,   // Índice del tile para las plataformas - ¡VERIFICAR!
		1,   // Índice de paredes de piedra
		100, // Índice de vacío
	}
	floatingPlatformTileIndex := tileIndices[0]
	groundExtendingPlatformTileIndex := tileIndices[1]

	platforms := definePlatforms()

	for _, platform := range platforms {
		// pixelX := platform.X * configs.C.TileSize
		// pixelY := platform.Y * configs.C.TileSize
		// pixelWidth := platform.W * configs.C.TileSize
		// pixelHeight := platform.H * configs.C.TileSize
		//
		// fmt.Printf("Plataforma en X:%d, Y:%d, Ancho:%d, Alto:%d (en píxeles: X:%d, Y:%d, Ancho:%d, Alto:%d)\n",
		// 	platform.X, platform.Y, platform.W, platform.H,
		// 	pixelX, pixelY, pixelWidth, pixelHeight)

		if platform.Type == PlatformTypeFloating {
			// Generar plataforma flotante (una sola fila)
			for x := 0; x < platform.W; x++ {
				tileX := platform.X + x
				tileY := platform.Y
				if tileX >= 0 && tileX < xCount && tileY >= 0 && tileY < yCount {
					index := tileY*xCount + tileX
					layer[index] = floatingPlatformTileIndex
				}
			}
		} else if platform.Type == PlatformTypeGroundExtending {
			// Generar plataforma que se extiende (con altura)
			for y := 0; y < platform.H; y++ {
				for x := 0; x < platform.W; x++ {
					tileX := platform.X + x
					tileY := platform.Y + y
					if tileX >= 0 && tileX < xCount && tileY >= 0 && tileY < yCount {
						index := tileY*xCount + tileX
						layer[index] = groundExtendingPlatformTileIndex
					}
				}
			}
		}
	}

	return layer
}
