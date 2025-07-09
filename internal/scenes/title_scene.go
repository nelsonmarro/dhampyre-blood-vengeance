package scenes

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/assets/audios"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/assets/fonts"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/configs"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/archetypes"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/state"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/systems"
)

type TitleScene struct {
	gameTitle       string
	pressEnterText  string
	titleFont       *text.GoTextFaceSource
	textFont        *text.GoTextFaceSource
	gameContext     state.GameContext // Change the type to GameContext
	backgroundImage *ebiten.Image
	once            sync.Once // Para asegurar que configure se llama solo una vez
	ECS             *ecs.ECS
}

// NewTitleScene creates a new TitleScene.
func NewTitleScene(gameContext state.GameContext) state.Scene {
	titleFont, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	textFont, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	world := donburi.NewWorld()
	return &TitleScene{
		ECS:            ecs.NewECS(world),
		gameTitle:      "DHAMPYRE: BLOOD VENGEANCE",
		pressEnterText: "PRESS ENTER TO START",
		titleFont:      titleFont,
		textFont:       textFont,
		gameContext:    gameContext,
	}
}

// Update handles input and logic for the title scene.
func (s *TitleScene) Update() error {
	s.once.Do(s.Configure)
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		audioQuery := donburi.NewQuery(filter.Contains(components.Audio))
		audioEntry, _ := audioQuery.First(s.ECS.World)
		audioComp := components.Audio.Get(audioEntry)
		audioComp.Playing = false
		mainScene := NewMainScene(s.gameContext)
		s.gameContext.SetScene(mainScene) // Set the stored MainScene
	}
	s.ECS.Update()
	return nil
}

func (s *TitleScene) Draw(screen *ebiten.Image) {
	const (
		normalFontSize = 12
		bigFontSize    = 28
	)

	screenWidth := configs.C.ScreenWidth
	screenHeight := configs.C.ScreenHeight

	screen.Fill(color.RGBA{0x00, 0x00, 0x00, 0xff})

	// Draw background image if loaded
	if s.backgroundImage != nil {
		op := &ebiten.DrawImageOptions{}
		// For a simple parallax effect (no movement yet), we can just draw it to fill the screen.
		// If the image is smaller than the screen, you might want to tile it or adjust the drawing.
		op.GeoM.Scale(float64(screenWidth)/float64(s.backgroundImage.Bounds().Dx()), float64(screenHeight)/float64(s.backgroundImage.Bounds().Dy()))
		screen.DrawImage(s.backgroundImage, op)
	}

	x := float64(screenWidth/2 - 200)
	y := float64(screenHeight/2 - 60)

	optitle := &text.DrawOptions{}
	optitle.GeoM.Translate(x, 70)
	optitle.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, s.gameTitle, &text.GoTextFace{
		Source: s.titleFont,
		Size:   bigFontSize,
	}, optitle)

	opstart := &text.DrawOptions{}
	opstart.GeoM.Translate(x+50, y+60)
	opstart.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, s.pressEnterText, &text.GoTextFace{
		Source: s.titleFont,
		Size:   bigFontSize,
	}, opstart)
}

// Configure is called when the scene is set as the active scene.
func (s *TitleScene) Configure() {
	fmt.Println("TitleScene configured")
	img, _, err := ebitenutil.NewImageFromFile("assets/CastleDungeonPack/Background/ParFull.png")
	if err != nil {
		log.Printf("Error loading ParFull.png: %v", err)
		return
	}
	s.backgroundImage = img

	audioEntity := archetypes.BackgroundAudio.Spawn(s.ECS)
	donburi.Add(audioEntity, components.Audio, &components.AudioComponent{
		Data:       audios.TitleScene_mp3,
		Format:     "mp3",
		Loop:       true,
		Playing:    true,
		Volume:     0.8, // Adjust as needed
		SampleRate: 48000,
	})
	s.ECS.AddSystem(systems.AudioSystemFunc)
}

// Dispose is called when the scene is no longer the active scene.
func (s *TitleScene) Dispose() {
	fmt.Println("TitleScene disposed")
}
