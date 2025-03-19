package systems

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/ganim8/v2"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/configs"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/factory"
)

func PlayerMovementSystemFunc(ecs *ecs.ECS, space *donburi.Entry, runAnimation, idleAnimation, jumpAnimation, attackAnimation, projectileAnimation *ganim8.Animation) {
	world := ecs.World

	playerQuery := donburi.NewQuery(filter.Contains(components.Position, components.Velocity, components.Sprite, components.Player, components.PlayerInput))

	for entry := range playerQuery.Iter(world) {
		velocity := components.Velocity.Get(entry)
		// position := components.Position.Get(entry)
		sprite := components.Sprite.Get(entry)
		player := components.Player.Get(entry)
		input := components.PlayerInput.Get(entry)

		// 1. Update Facing Direction based on Input
		if input.MovingLeft {
			player.FacingLeft = true
		} else if input.MovingRight {
			player.FacingLeft = false
		}

		// 2. Set Flipped state based on Facing Direction
		if player.FacingLeft != sprite.Flipped {
			sprite.Flipped = player.FacingLeft
		}

		if input.Attacking && !player.IsAttacking && velocity.OnGround && velocity.X == 0 && sprite.AnimationName == "idle" {
			player.IsAttacking = true
			sprite.AnimationName = "attack"
			sprite.Animation = attackAnimation
			// Ensure the attack animation respects the last facing direction
			if player.FacingLeft != sprite.Flipped {
				sprite.Flipped = player.FacingLeft
			}
			if sprite.Animation.IsEnd() {
				sprite.Animation.GoToFrame(1)
				sprite.Animation.Resume()
			}
			// Create a new projectile
			ownerPosition := components.Position.Get(entry)
			startX := ownerPosition.X + 8
			startY := ownerPosition.Y + 15 // Temporarily raise the start Y
			var velocityX float64 = 5
			if player.FacingLeft {
				velocityX = -5
			}
			factory.CreateProjectile(ecs, projectileAnimation, entry, startX, startY, velocityX, 0, player.FacingLeft, space)
		}

		if sprite.AnimationName == "attack" && sprite.Animation.IsEnd() {
			player.IsAttacking = false
			if velocity.X != 0 {
				sprite.AnimationName = "run"
				sprite.Animation = runAnimation
			} else {
				sprite.AnimationName = "idle"
				sprite.Animation = idleAnimation
			}
		}

		if !player.IsAttacking {
			if input.MovingLeft {
				velocity.X = -2
				if sprite.AnimationName != "run" {
					sprite.AnimationName = "run"
					sprite.Animation = runAnimation
				}
			} else if input.MovingRight {
				velocity.X = 2
				if sprite.AnimationName != "run" {
					sprite.AnimationName = "run"
					sprite.Animation = runAnimation
				}
			} else {
				velocity.X = 0
				if velocity.OnGround && sprite.AnimationName != "idle" && sprite.AnimationName != "jump" {
					sprite.AnimationName = "idle"
					sprite.Animation = idleAnimation
				}
			}

			if input.Jumping && velocity.OnGround {
				velocity.YSpeed = configs.C.JumpSpeed
				velocity.OnGround = false
				sprite.AnimationName = "jump"
				sprite.Animation = jumpAnimation
			}
		}

		velocity.YSpeed += configs.C.Gravity

		// Transiciones de animación basadas en el estado (si no está atacando)
		if !player.IsAttacking {
			if velocity.OnGround && sprite.AnimationName == "jump" { // Transición de salto a correr/idle al aterrizar
				if velocity.X != 0 {
					sprite.AnimationName = "run"
					sprite.Animation = runAnimation
				} else {
					sprite.AnimationName = "idle"
					sprite.Animation = idleAnimation
				}
			} else if !velocity.OnGround && sprite.AnimationName != "jump" {
				sprite.AnimationName = "jump"
				sprite.Animation = jumpAnimation
			}
		}
		// fmt.Println("Animation Name:", sprite.AnimationName) // Para debug
	}
}
