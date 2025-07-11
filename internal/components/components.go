package components

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv" // Import resolv
	"github.com/yohamta/donburi"
	"github.com/yohamta/ganim8/v2"
)

// Posición Component
var Position = donburi.NewComponentType[PositionComponent]()

type PositionComponent struct {
	X float64
	Y float64
}

// TiledBackground contiene la información del fondo de tiles.
type TiledBackgroundComponent struct {
	TilesImage *ebiten.Image
	Layers     [][]int
}

var TiledBackground = donburi.NewComponentType[TiledBackgroundComponent]()

// Sprite Component
var Sprite = donburi.NewComponentType[SpriteComponent]()

type SpriteComponent struct {
	Animation             *ganim8.Animation
	AnimationName         string
	AnimationSpeed        float64
	AnimationFinished     bool
	AnimationFrameCounter int
	Flipped               bool // Nuevo campo para rastrear el estado de flip
}

// Velocity Component
var Velocity = donburi.NewComponentType[VelocityComponent]()

type VelocityComponent struct {
	X        float64
	Y        float64
	YSpeed   float64
	OnGround bool
}

type PlayerComponent struct {
	IsAttacking bool
	FacingLeft  bool
	IsDead      bool
}

var Player = donburi.NewComponentType[PlayerComponent]()

type EnemyComponent struct {
	InitialPosition float64
	CurrentState    EnemyMovementState
	StateTimer      *time.Timer
	MovingLeftNext  bool // Add this line
	IsDead          bool
}

// EnemyMovementState represents the current movement state of the enemy.
type EnemyMovementState int

const (
	EnemyStateIdle                 EnemyMovementState = 0
	EnemyStateMovingLeft           EnemyMovementState = 1
	EnemyStateReturningCenterLeft  EnemyMovementState = 2
	EnemyStateMovingRight          EnemyMovementState = 3
	EnemyStateReturningCenterRight EnemyMovementState = 4
)

var Enemy = donburi.NewComponentType[EnemyComponent]()

type ProjectileComponent struct {
	Owner  *donburi.Entry
	Damage int
}

var Projectile = donburi.NewComponentType[ProjectileComponent]()

// Health Component
var Health = donburi.NewComponentType[HealthComponent]()

type HealthComponent struct {
	Current int
	Max     int
}

// Magic Component
var Magic = donburi.NewComponentType[MagicComponent]()

type MagicComponent struct {
	Current int
	Max     int
}

// Object Component para Resolv - MODIFIED: Now stores IShape
var (
	Shape       = donburi.NewComponentType[resolv.ConvexPolygon]()
	ShapeCircle = donburi.NewComponentType[resolv.Circle]()
)

// Space Component para Resolv Space
var Space = donburi.NewComponentType[resolv.Space]()
