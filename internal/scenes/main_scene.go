package scenes

import (
	"image/color"
	"log"
	"sync"
	"time"

	"github.com/solarlune/resolv"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/assets/audios"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/configs"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/archetypes"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components"
	dresolv "github.com/nelsonmarro/dhampyre-blood-vengeance/internal/resolv" // Importar nuestro resolv wrapper
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/sprites"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/state"
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
	World               donburi.World // Añadimos el world, aunque en PlatformerScene no está explícito aquí
	ECS                 *ecslib.ECS
	once                sync.Once // Para asegurar que configure se llama solo una vez
	TiledBackground     *tiled.TiledBackground
	backgroundImage     *ebiten.Image
	IdleAnimation       *ganim8.Animation // Usaremos ganim8.Animation
	RunAnimation        *ganim8.Animation // Usaremos ganim8.Animation
	JumpAnimation       *ganim8.Animation // Usaremos ganim8.Animation
	AttackAnimation     *ganim8.Animation // Usaremos ganim8.Animation
	DeadAnimation       *ganim8.Animation // Usaremos ganim8.Animation
	ProjectileAnimation *ganim8.Animation // Usaremos ganim8.Animation
	EnemyIdleAnimation  *ganim8.Animation // Usaremos ganim8.Animation
	EnemyRunAnimation   *ganim8.Animation // Usaremos ganim8.Animation
	EnemyDeadAnimation  *ganim8.Animation // Usaremos ganim8.Animation
	SpaceEntry          *donburi.Entry
	gameContext         state.GameContext
}

func NewMainScene(gameContext state.GameContext) state.Scene {
	world := donburi.NewWorld()
	s := &MainScene{
		World:           world,
		ECS:             ecs.NewECS(world), // Use the ECS instance passed as an argument
		TiledBackground: tiled.NewTiledBackground(),
		gameContext:     gameContext,
	}
	return s
}

func (s *MainScene) Update() error {
	s.once.Do(s.Configure) // Asegurar que configure se llama solo la primera vez
	s.ECS.Update()         // Actualizar el ECS (ejecutar sistemas)
	s.RunAnimation.Update()
	s.JumpAnimation.Update()
	s.IdleAnimation.Update()
	s.AttackAnimation.Update()
	s.ProjectileAnimation.Update()
	s.EnemyIdleAnimation.Update()
	s.DeadAnimation.Update()
	s.EnemyRunAnimation.Update()

	playerQuery := donburi.NewQuery(filter.Contains(components.Player))
	if playerQuery.Count(s.World) == 0 {
		audioQuery := donburi.NewQuery(filter.Contains(components.Audio))
		audioEntry, _ := audioQuery.First(s.ECS.World)
		audioComp := components.Audio.Get(audioEntry)
		audioComp.Playing = false
		endScene := NewEndScene(s.gameContext)
		s.gameContext.SetScene(endScene)
	}

	return nil
}

func (s *MainScene) Draw(screen *ebiten.Image) {
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
	s.TiledBackground.DrawTileBackground(screen) // Dibujar el fondo de azulejos ANTES que otros elementos
	s.ECS.DrawLayer(layers.Default, screen)      // Dibujar los elementos del juego (sprites, etc.)

	// playerQuery := donburi.NewQuery(filter.Contains(components.Player))
	// playerEntry, _ := playerQuery.First(*s.World) // Obtén la entidad del jugador (maneja el error si es necesario)
	// if playerEntry != nil {
	// 	playerShapeComponent := components.Shape.Get(playerEntry)
	// 	if playerShapeComponent != nil {
	// 		s.drawResolvShape(screen, playerShapeComponent, color.RGBA{0, 0, 255, 128}) // Dibuja la colisión del jugador en azul
	// 	}
	// }
	// airPlatformShape := resolv.NewRectangle(272, float64(configs.C.ScreenHeight/2-22), float64(configs.C.TileSize), float64(configs.C.TileSize*5))
	// s.drawResolvShape(screen, airPlatformShape, color.RGBA{0, 222, 0, 128})
	//
	// tileWidth := float64(configs.C.TileSize) // Obtén el tamaño de un tile
	//
	// platformHeight := float64(configs.C.ScreenHeight/2 - 150)
	// platformShape := resolv.NewRectangle(400, 300, tileWidth, platformHeight)
	// s.drawResolvShape(screen, platformShape, color.RGBA{255, 0, 0, 128})
}

func (s *MainScene) Configure() {
	img, _, err := ebitenutil.NewImageFromFile("assets/CastleDungeonPack/Background/ParFull.png")
	if err != nil {
		log.Printf("Error loading ParFull.png: %v", err)
		return
	}
	s.backgroundImage = img

	audioEntity := archetypes.BackgroundAudio.Spawn(s.ECS)
	donburi.Add(audioEntity, components.Audio, &components.AudioComponent{
		Data:       audios.MainScene_mp3,
		Format:     "mp3",
		Loop:       true,
		Playing:    true,
		Volume:     7.0, // Adjust as needed
		SampleRate: 48000,
	})

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

	deadImage, _, err := ebitenutil.NewImageFromFile("assets/player_sprites/Dead.png")
	if err != nil {
		log.Fatal("Error al cargar Dead.png:", err)
	}

	enemyIdleImage, _, err := ebitenutil.NewImageFromFile("assets/enemies/Skeleton/Idle.png")
	enemyRunImage, _, err := ebitenutil.NewImageFromFile("assets/enemies/Skeleton/Walk.png")
	enemyDeadImage, _, err := ebitenutil.NewImageFromFile("assets/enemies/Skeleton/Dead.png")

	idleGrid := ganim8.NewGrid(playerSize, playerSize, idleImage.Bounds().Dx(), idleImage.Bounds().Dy())
	runGrid := ganim8.NewGrid(playerSize, playerSize, runImage.Bounds().Dx(), runImage.Bounds().Dy())
	jumpGrid := ganim8.NewGrid(playerSize, playerSize, jumpImage.Bounds().Dx(), jumpImage.Bounds().Dy())
	attackGrid := ganim8.NewGrid(playerSize, playerSize, attackImage.Bounds().Dx(), attackImage.Bounds().Dy())
	deadGrid := ganim8.NewGrid(playerSize, playerSize, deadImage.Bounds().Dx(), deadImage.Bounds().Dy())

	projectileGrid := ganim8.NewGrid(52, 48, projectileImage.Bounds().Dx(), projectileImage.Bounds().Dy())

	enemyIdleGrid := ganim8.NewGrid(playerSize, playerSize, enemyIdleImage.Bounds().Dx(), enemyIdleImage.Bounds().Dy())
	runEnemyGrid := ganim8.NewGrid(playerSize, playerSize, enemyRunImage.Bounds().Dx(), enemyRunImage.Bounds().Dy())
	deadEnemyGrid := ganim8.NewGrid(playerSize, playerSize, enemyDeadImage.Bounds().Dx(), enemyDeadImage.Bounds().Dy())

	// Crear las animaciones
	s.EnemyIdleAnimation = ganim8.New(enemyIdleImage, enemyIdleGrid.Frames("1-7", 1), 150*time.Millisecond) // Asumiendo 4 frames de idle
	s.EnemyRunAnimation = ganim8.New(enemyRunImage, runEnemyGrid.Frames("1-7", 1), 150*time.Millisecond)    // Asumiendo 6 frames de correr
	s.EnemyDeadAnimation = ganim8.New(enemyDeadImage, deadEnemyGrid.Frames("1-3", 1), 250*time.Millisecond) // Asumiendo 6 frames de correr

	s.IdleAnimation = ganim8.New(idleImage, idleGrid.Frames("1-5", 1), 250*time.Millisecond)                   // Asumiendo 4 frames de idle
	s.RunAnimation = ganim8.New(runImage, runGrid.Frames("1-6", 1), 150*time.Millisecond)                      // Asumiendo 6 frames de correr
	s.DeadAnimation = ganim8.New(deadImage, deadGrid.Frames("1-8", 1), 200*time.Millisecond)                   // Asumiendo 6 frames de correr
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

	initialX := float64((configs.C.ScreenHeight / 2) - 210)
	initialY := float64(configs.C.ScreenHeight - (8 * configs.C.TileSize))
	factory.CreatePlayer(s.ECS, s.IdleAnimation, spaceEntry, initialX, initialY)

	factory.CreateEnemySkeleton(s.ECS, &components.PositionComponent{X: float64((configs.C.ScreenWidth / 2)), Y: float64(configs.C.ScreenHeight - (8 * configs.C.TileSize))},
		s.EnemyIdleAnimation, spaceEntry)
	factory.CreateEnemySkeleton(s.ECS, &components.PositionComponent{X: float64((configs.C.ScreenWidth/2 + 250)), Y: float64(configs.C.ScreenHeight - (8 * configs.C.TileSize))},
		s.EnemyIdleAnimation, spaceEntry)

	s.buildObstacles(spaceEntry)

	s.ECS.AddSystem(systems.PlayerInputSystemFunc)
	s.ECS.AddSystem(func(e *ecs.ECS) { // Pasar la referencia a la escena al sistema de movimiento
		systems.PlayerMovementSystemFunc(e, spaceEntry, s.RunAnimation, s.IdleAnimation, s.JumpAnimation, s.AttackAnimation, s.ProjectileAnimation)
	})
	s.ECS.AddSystem(systems.ResolvMovementSystemFunc)
	s.ECS.AddSystem(systems.ProjectileSystemFunc)
	s.ECS.AddSystem(systems.ProjectileCollisionSystemFunc)
	s.ECS.AddSystem(func(e *ecs.ECS) { // Pasar la referencia a la escena al sistema de movimiento
		systems.EnemyMovementSystemFunc(e, s.EnemyRunAnimation, s.EnemyIdleAnimation)
	})
	s.ECS.AddSystem(func(e *ecs.ECS) { // Pasar la referencia a la escena al sistema de movimiento
		systems.EnemyCollisionSystemFunc(e, s.DeadAnimation)
	})

	s.ECS.AddSystem(func(e *ecs.ECS) { // Pasar la referencia a la escena al sistema de movimiento
		systems.EnemyProjectileCollisionSystemFunc(e, s.EnemyDeadAnimation)
	})
	s.ECS.AddSystem(func(e *ecs.ECS) { // Pasar la referencia a la escena al sistema de movimiento
		systems.AudioSystemFunc(e)
	})

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
	platformHeight := float64(configs.C.ScreenHeight/2 - 150)

	// Crea el rectángulo de resolv con el ancho de un tile
	platformShape := resolv.NewRectangle(400, 300, tileWidth, platformHeight)

	platform := factory.CreateObstacle(s.ECS, platformShape)

	airPlatformShape := resolv.NewRectangle(272, float64(configs.C.ScreenHeight/2-22), float64(configs.C.TileSize), float64(configs.C.TileSize*5))
	airPlatform := factory.CreateObstacle(s.ECS, airPlatformShape)

	dresolv.Add(s.SpaceEntry, platform, airPlatform, wallLeft, wallRight, wallTop, wallBottom)
}

// Dispose is called when the scene is no longer the active scene.
func (s *MainScene) Dispose() {
}

func (s *MainScene) drawResolvShape(screen *ebiten.Image, shape *resolv.ConvexPolygon, color color.Color) {
	vector.DrawFilledRect(screen, float32(shape.Position().X), float32(shape.Position().Y), float32(shape.Bounds().Height()), float32(shape.Bounds().Width()), color, false)
}
