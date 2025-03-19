package scenes

import (
	"image/color"
	"log"
	"sync"
	"time"

	"github.com/solarlune/resolv"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/configs"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components"
	dresolv "github.com/nelsonmarro/dhampyre-blood-vengeance/internal/resolv" // Importar nuestro resolv wrapper
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/sprites"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/systems"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/tiled"

	"github.com/hajimehoshi/ebiten/v2" // Import resolv
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	ecslib "github.com/yohamta/donburi/ecs" // Usamos alias ecslib para evitar confusiones
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/ganim8/v2"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/factory"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/layers"
)

// MainScene representa la escena principal del juego (el nivel)
type MainScene struct {
	World               *donburi.World // Añadimos el world, aunque en PlatformerScene no está explícito aquí
	ECS                 *ecslib.ECS
	once                sync.Once // Para asegurar que configure se llama solo una vez
	TiledBackground     *tiled.TiledBackground
	IdleAnimation       *ganim8.Animation // Usaremos ganim8.Animation
	RunAnimation        *ganim8.Animation // Usaremos ganim8.Animation
	JumpAnimation       *ganim8.Animation // Usaremos ganim8.Animation
	AttackAnimation     *ganim8.Animation // Usaremos ganim8.Animation
	ProjectileAnimation *ganim8.Animation // Usaremos ganim8.Animation
	SpaceEntry          *donburi.Entry
}

func NewLevelScene() *MainScene {
	configs.InitGameConfig()
	world := donburi.NewWorld() // Inicializamos el world aquí
	s := &MainScene{
		World:           &world,               // Inicializamos el world aquí
		ECS:             ecslib.NewECS(world), // Inicializamos el ECS aquí también (redundante, pero para seguir el ejemplo)
		TiledBackground: tiled.NewTiledBackground(),
	}
	return s
}

func (s *MainScene) Update() error {
	s.once.Do(s.configure) // Asegurar que configure se llama solo la primera vez
	s.ECS.Update()         // Actualizar el ECS (ejecutar sistemas)
	s.RunAnimation.Update()
	s.JumpAnimation.Update()
	s.IdleAnimation.Update()
	s.AttackAnimation.Update()
	s.ProjectileAnimation.Update()
	return nil
}

func (s *MainScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)                     // Rellenar pantalla con negro (color de fondo)
	s.TiledBackground.DrawTileBackground(screen) // Dibujar el fondo de azulejos ANTES que otros elementos
	s.ECS.DrawLayer(layers.Default, screen)      // Dibujar los elementos del juego (sprites, etc.)

	playerQuery := donburi.NewQuery(filter.Contains(components.Player))
	playerEntry, _ := playerQuery.First(*s.World) // Obtén la entidad del jugador (maneja el error si es necesario)
	if playerEntry != nil {
		playerShapeComponent := components.Shape.Get(playerEntry)
		if playerShapeComponent != nil {
			s.drawResolvShape(screen, playerShapeComponent, color.RGBA{0, 0, 255, 128}) // Dibuja la colisión del jugador en azul
		}
	}
	airPlatformShape := resolv.NewRectangle(272, 185, 16, 16)
	s.drawResolvShape(screen, airPlatformShape, color.RGBA{0, 222, 0, 128})

	tileWidth := float64(configs.C.TileSize) // Obtén el tamaño de un tile
	platformHeight := float64(configs.C.ScreenHeight/2 - configs.C.TileSize)

	// Crea el rectángulo de resolv con el ancho de un tile
	platformShape := resolv.NewRectangle(390, 370, tileWidth, platformHeight)
	s.drawResolvShape(screen, platformShape, color.RGBA{255, 0, 0, 128})
}

func (s *MainScene) configure() {
	tilesImage := sprites.LoadTilesImg()
	s.TiledBackground.TilesImage = tilesImage
	playerSize := configs.C.PlayerSize

	idleImage, _, err := ebitenutil.NewImageFromFile("assets/player_sprites/Idle.png")
	if err != nil {
		log.Fatal("Error al cargar Idle.png:", err)
	}
	runImage, _, err := ebitenutil.NewImageFromFile("assets/player_sprites/Run.png")
	if err != nil {
		log.Fatal("Error al cargar Run.png:", err)
	}
	jumpImage, _, err := ebitenutil.NewImageFromFile("assets/player_sprites/Jump.png")
	if err != nil {
		log.Fatal("Error al cargar Jump.png:", err)
	}

	attackImage, _, err := ebitenutil.NewImageFromFile("assets/player_sprites/Attack_1.png")
	if err != nil {
		log.Fatal("Error al cargar Attack_1.png:", err)
	}

	projectileImage, _, err := ebitenutil.NewImageFromFile("assets/player_sprites/Blood_Charge_1.png")
	if err != nil {
		log.Fatal("Error al cargar Attack_1.png:", err)
	}

	idleGrid := ganim8.NewGrid(playerSize, playerSize, idleImage.Bounds().Dx(), idleImage.Bounds().Dy())
	runGrid := ganim8.NewGrid(playerSize, playerSize, runImage.Bounds().Dx(), runImage.Bounds().Dy())
	jumpGrid := ganim8.NewGrid(playerSize, playerSize, jumpImage.Bounds().Dx(), jumpImage.Bounds().Dy())
	attackGrid := ganim8.NewGrid(playerSize, playerSize, attackImage.Bounds().Dx(), attackImage.Bounds().Dy())
	projectileGrid := ganim8.NewGrid(52, 48, projectileImage.Bounds().Dx(), projectileImage.Bounds().Dy())

	// Crear las animaciones
	s.IdleAnimation = ganim8.New(idleImage, idleGrid.Frames("1-5", 1), 250*time.Millisecond)                   // Asumiendo 4 frames de idle
	s.RunAnimation = ganim8.New(runImage, runGrid.Frames("1-6", 1), 150*time.Millisecond)                      // Asumiendo 6 frames de correr
	s.ProjectileAnimation = ganim8.New(projectileImage, projectileGrid.Frames("1-3", 1), 230*time.Millisecond) // Asumiendo 6 frames de correr

	s.AttackAnimation = ganim8.New(attackImage, attackGrid.Frames("1-6", 1), 80*time.Millisecond, func(anim *ganim8.Animation, loops int) {
		anim.PauseAtEnd()
	}) // Asumiendo 6 frames de correr

	s.JumpAnimation = ganim8.New(jumpImage, jumpGrid.Frames("1-2", 1, "3-4", 1, "5-6", 1), map[string]time.Duration{
		"1-2": time.Millisecond * 160,
		"3-4": time.Millisecond * 160,
		"5-6": time.Millisecond * 160,
	}, func(anim *ganim8.Animation, loops int) {
		anim.PauseAtEnd()
	}) // Asumiendo 6 frames de salto
	spaceEntry := factory.CreateSpace(s.ECS)
	s.SpaceEntry = spaceEntry

	factory.CreatePlayer(s.ECS, s.IdleAnimation, spaceEntry)

	s.buildObstacles(spaceEntry)

	s.ECS.AddSystem(systems.PlayerInputSystemFunc)
	s.ECS.AddSystem(func(e *ecs.ECS) { // Pasar la referencia a la escena al sistema de movimiento
		systems.PlayerMovementSystemFunc(e, spaceEntry, s.RunAnimation, s.IdleAnimation, s.JumpAnimation, s.AttackAnimation, s.ProjectileAnimation)
	})
	s.ECS.AddSystem(systems.ResolvMovementSystemFunc)
	s.ECS.AddSystem(systems.ProjectileSystemFunc)
	s.ECS.AddSystem(systems.ProjectileCollisionSystemFunc)

	drawSystem := systems.DrawSystemFunc
	s.ECS.AddRenderer(layers.Default, drawSystem) // Renderer de dibujado en la capa por defecto
}

func (s *MainScene) buildObstacles(spaceEntry *donburi.Entry) {
	moWidth, moHeight := ebiten.WindowSize()
	screenHeight := float64(moHeight)
	screenWidth := float64(moWidth)
	wallThickness := 16.0 // Grosor de los bordes

	// Crear los bordes de la pantalla como paredes estáticas
	wallShapeLeft := resolv.NewRectangle(-float64(72), 0, wallThickness, screenHeight)
	wallLeft := factory.CreateStaticWall(s.ECS, wallShapeLeft)
	wallShapeRight := resolv.NewRectangle(screenWidth/2, 0, wallThickness, screenHeight)
	wallRight := factory.CreateStaticWall(s.ECS, wallShapeRight)

	wallShapeTop := resolv.NewRectangle(0, -wallThickness, screenWidth, wallThickness)
	wallTop := factory.CreateStaticWall(s.ECS, wallShapeTop)
	wallShapeBottom := resolv.NewRectangle(0, screenHeight/2-float64(configs.C.PlayerSize)+9, screenWidth, wallThickness)
	wallBottom := factory.CreateStaticWall(s.ECS, wallShapeBottom)

	tileWidth := float64(configs.C.TileSize) // Obtén el tamaño de un tile
	platformHeight := float64(configs.C.ScreenHeight/2 - configs.C.TileSize)

	// Crea el rectángulo de resolv con el ancho de un tile
	platformShape := resolv.NewRectangle(390, 370, tileWidth, platformHeight)
	platform := factory.CreateObstacle(s.ECS, platformShape)

	airPlatformShape := resolv.NewRectangle(272, 185, 16, 16)
	airPlatform := factory.CreateObstacle(s.ECS, airPlatformShape)

	dresolv.Add(s.SpaceEntry, platform, airPlatform, wallLeft, wallRight, wallTop, wallBottom)
}

func (s *MainScene) drawResolvShape(screen *ebiten.Image, shape *resolv.ConvexPolygon, color color.Color) {
	vector.DrawFilledRect(screen, float32(shape.Position().X), float32(shape.Position().Y), float32(shape.Bounds().Height()), float32(shape.Bounds().Width()), color, false)
}
