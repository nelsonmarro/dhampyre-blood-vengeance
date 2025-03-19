package sprites

import (
	"image"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/configs"
)

func LoadSpriteSheet(filepath string, frameWidth, frameHeight int) ([]*ebiten.Image, error) {
	imgFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	eImg := ebiten.NewImageFromImage(img)
	frames := []*ebiten.Image{}
	cols := eImg.Bounds().Dx() / frameWidth
	rows := eImg.Bounds().Dy() / frameHeight

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			frame := eImg.SubImage(image.Rect(x*frameWidth, y*frameHeight, (x+1)*frameWidth, (y+1)*frameHeight)).(*ebiten.Image)
			frames = append(frames, frame)
		}
	}

	return frames, nil
}

func LoadTilesImg() *ebiten.Image {
	tilesImage, _, err := ebitenutil.NewImageFromFile("assets/CastleDungeonPack/Tileset/Tileset.png") // *** ¡VERIFICAR RUTA! ***
	if err != nil {
		log.Fatal("Error al cargar tile set gótico:", err)
	}
	return tilesImage
}

func LoadPlayerIdleSheet() []*ebiten.Image {
	idleFrames, err := LoadSpriteSheet("assets/player_sprites/Idle.png", configs.C.PlayerSize, configs.C.PlayerSize)
	if err != nil {
		log.Fatal("Error al cargar sprite sheet de idle:", err)
	}
	return idleFrames
}

func LoadPlayerRunSheet() []*ebiten.Image {
	// Cargar el sprite sheet de correr para el jugador
	runFrames, err := LoadSpriteSheet("assets/player_sprites/Run.png", configs.C.PlayerSize, configs.C.PlayerSize)
	if err != nil {
		log.Fatal("Error al cargar sprite sheet de correr:", err)
	}

	return runFrames
}

func LoadPlayerJumpSheet() []*ebiten.Image {
	jumpFrames, err := LoadSpriteSheet("assets/player_sprites/Jump.png", configs.C.PlayerSize, configs.C.PlayerSize)
	if err != nil {
		log.Fatal("Error al cargar sprite sheet de salto:", err)
	}
	return jumpFrames
}
