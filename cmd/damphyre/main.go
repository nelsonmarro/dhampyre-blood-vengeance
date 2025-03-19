package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/configs"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/scenes"
)

type Game struct {
	currentScene *scenes.MainScene
}

func NewGame() *Game {
	g := &Game{
		currentScene: scenes.NewLevelScene(),
	}
	return g
}

func (g *Game) Update() error {
	return g.currentScene.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.currentScene.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return configs.C.ScreenWidth, configs.C.ScreenHeight
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(configs.C.ScreenWidth*2, configs.C.ScreenHeight*2)
	ebiten.SetWindowTitle("Dhampyre: Noche Carmesí")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
