package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/configs"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/scenes"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/state"
)

type Game struct {
	currentScene state.Scene
	ecs          *ecs.ECS
	World        donburi.World
}

func NewGame() *Game {
	w := donburi.NewWorld()
	g := &Game{
		World: w,
		ecs:   ecs.NewECS(w),
	}
	titleScene := scenes.NewTitleScene(g) // Create the title scene
	g.SetScene(titleScene)
	fmt.Printf("Game instance created: %+v\n", g) // Add this line
	return g
}

func (g *Game) Update() error {
	if g.currentScene != nil {
		return g.currentScene.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.currentScene != nil {
		g.currentScene.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return configs.C.ScreenWidth, configs.C.ScreenHeight
}

func (g *Game) SetScene(newScene state.Scene) {
	if g.currentScene != nil {
		g.currentScene.Dispose()
		g.currentScene = nil
	}
	g.currentScene = newScene
}

func (g *Game) ECSManager() *ecs.ECS {
	return g.ecs
}

func (g *Game) GetWorld() donburi.World {
	return g.World
}

func main() {
	configs.InitGameConfig() // Call the configuration initialization here

	fmt.Printf("ScreenWidth from configs: %d\n", configs.C.ScreenWidth)
	fmt.Printf("ScreenHeight from configs: %d\n", configs.C.ScreenHeight)

	ebiten.SetWindowSize(configs.C.ScreenWidth*2, configs.C.ScreenHeight*2)
	ebiten.SetWindowTitle("Dhampyre: Noche Carmesí")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	// ebiten.SetFullscreen(true)
	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
