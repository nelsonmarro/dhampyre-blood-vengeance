package systems

import (
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/ganim8/v2"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components"
)

// EnemyMovementSystemFunc implements the enemy movement pattern: idle, move left, idle, move right, repeat.
func EnemyMovementSystemFunc(ecs *ecs.ECS, runAnimation, idleAnimation *ganim8.Animation) {
	world := ecs.World

	enemyQuery := donburi.NewQuery(filter.Contains(components.Position, components.Velocity, components.Sprite, components.Enemy))

	enemyQuery.EachEntity(world, func(entry *donburi.Entry) {
		position := components.Position.Get(entry)
		velocity := components.Velocity.Get(entry)
		sprite := components.Sprite.Get(entry)
		enemy := components.Enemy.Get(entry)

		if enemy.IsDead {
			return // Early return if enemy is dead
		}

		const moveDistance = 40.0
		const moveSpeed = 1
		const idleDuration = 1 * time.Second

		switch enemy.CurrentState {
		case components.EnemyStateIdle:
			velocity.X = 0
			if sprite.AnimationName != "eidle" {
				sprite.AnimationName = "eidle"
				sprite.Animation = idleAnimation
				// position.X = enemy.InitialPosition // Consider if you need to reset position here
			}
			if enemy.StateTimer == nil {
				enemy.StateTimer = time.NewTimer(idleDuration)
			}

			select {
			case <-enemy.StateTimer.C:
				components.Enemy.SetValue(entry, components.EnemyComponent{
					InitialPosition: enemy.InitialPosition,
					CurrentState:    components.EnemyStateMovingLeft,
					StateTimer:      nil,
					MovingLeftNext:  enemy.MovingLeftNext, // Preserve the next move direction
				})
			default:
				// Do nothing, still idling
			}

		case components.EnemyStateMovingLeft:
			velocity.X = -moveSpeed
			sprite.Flipped = true
			if sprite.AnimationName != "erun" {
				sprite.AnimationName = "erun"
				sprite.Animation = runAnimation
			}
			if position.X <= enemy.InitialPosition-moveDistance {
				components.Enemy.SetValue(entry, components.EnemyComponent{
					InitialPosition: enemy.InitialPosition,
					CurrentState:    components.EnemyStateReturningCenterLeft,
					StateTimer:      nil,
					MovingLeftNext:  enemy.MovingLeftNext,
				})
			}

		case components.EnemyStateReturningCenterLeft:
			if position.X < enemy.InitialPosition {
				velocity.X = moveSpeed // Corrected: Set velocity directly
				sprite.Flipped = false // Corrected: Should be true when returning
				if sprite.AnimationName != "erun" {
					sprite.AnimationName = "erun"
					sprite.Animation = runAnimation
				}
			} else {
				velocity.X = 0
				components.Enemy.SetValue(entry, components.EnemyComponent{
					InitialPosition: enemy.InitialPosition,
					CurrentState:    components.EnemyStateIdle,
					StateTimer:      nil,
					MovingLeftNext:  false, // Set to move right next
				})
			}
		}
		position.X += velocity.X // Apply velocity to position
	})
}
