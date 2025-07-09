package systems

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/ganim8/v2"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/assets/audios"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/configs"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/factory"
)

func PlayerMovementSystemFunc(ecs *ecs.ECS, space *donburi.Entry, runAnimation, idleAnimation, jumpAnimation, attackAnimation, projectileAnimation *ganim8.Animation) {
	world := ecs.World

	playerQuery := donburi.NewQuery(filter.Contains(components.Position, components.Velocity, components.Sprite, components.Player, components.PlayerInput))

	playerQuery.EachEntity(world, func(entry *donburi.Entry) {
		velocity := components.Velocity.Get(entry)
		sprite := components.Sprite.Get(entry)
		player := components.Player.Get(entry)
		input := components.PlayerInput.Get(entry)

		if player.IsDead {
			return // Early return if player is dead
		}

		// Update Facing Direction
		if input.MovingLeft {
			player.FacingLeft = true
		} else if input.MovingRight {
			player.FacingLeft = false
		}

		// Set Flipped state
		if player.FacingLeft != sprite.Flipped {
			sprite.Flipped = player.FacingLeft
		}

		// Attack logic
		if input.Attacking && !player.IsAttacking && velocity.OnGround && velocity.X == 0 && sprite.AnimationName == "idle" {
			player.IsAttacking = true
			sprite.AnimationName = "attack"
			sprite.Animation = attackAnimation
			if sprite.Animation.IsEnd() {
				sprite.Animation.GoToFrame(1)
				sprite.Animation.Resume()
			}
			ownerPosition := components.Position.Get(entry)
			startX := ownerPosition.X + 8
			startY := ownerPosition.Y + 15
			velocityX := 5.0
			if player.FacingLeft {
				velocityX = -5
			}
			proejectile := factory.CreateProjectile(ecs, projectileAnimation, entry, startX, startY, velocityX, 0, player.FacingLeft, space)
			projectileAudio := components.Audio.Get(proejectile)
			projectileAudio.Playing = true
			projectileAudio.Data = audios.Attack_mp3

			return // Early return after attack
		}

		// Attack animation end transition
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

		// Movement and animation logic when not attacking
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

		// Animation transitions based on state (if not attacking)
		if !player.IsAttacking {
			if velocity.OnGround && sprite.AnimationName == "jump" {
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
	})
}
