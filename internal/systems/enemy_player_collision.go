package systems

import (
	"fmt"
	"time"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/assets/audios"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components"
	dresolv "github.com/nelsonmarro/dhampyre-blood-vengeance/internal/resolv"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/ganim8/v2"
)

func EnemyCollisionSystemFunc(ecs *ecs.ECS, playerDeathAnim *ganim8.Animation) {
	playerQuery := donburi.NewQuery(filter.Contains(components.Player, components.Position, components.Shape, components.Sprite))
	spaceQuery := donburi.NewQuery(filter.Contains(components.Space))

	var spaceEntry *donburi.Entry
	spaceEntry, _ = spaceQuery.FirstEntity(ecs.World)
	if spaceEntry == nil {
		return // No space to check collisions against
	}
	space := components.Space.Get(spaceEntry)

	var playerEntry *donburi.Entry
	playerEntry, _ = playerQuery.FirstEntity(ecs.World)
	if playerEntry == nil {
		return // No player to check collision for
	}
	player := components.Player.Get(playerEntry)
	playerPosition := components.Position.Get(playerEntry)
	playerShape := components.Shape.Get(playerEntry)
	playerSprite := components.Sprite.Get(playerEntry)
	playerAudio := components.Audio.Get(playerEntry)

	if player.IsDead {
		return // Player is already dead, so no need to check for collisions
	}

	if playerPosition == nil || playerShape == nil {
		return // Player missing position or shape
	}

	playerTestShape := playerShape.Clone()
	playerTestShape.SetPosition(playerPosition.X-float64(playerShape.Bounds().Width()/2+20), playerPosition.Y-float64(playerShape.Bounds().Height()/2)) // Adjust position if needed

	// Check for collisions with Enemies
	collision := playerTestShape.IntersectionTest(resolv.IntersectionTestSettings{
		TestAgainst: space.FilterShapes().ByTags(dresolv.TagEnemy),
		OnIntersect: func(set resolv.IntersectionSet) bool {
			fmt.Println("Player intersected with an enemy!")

			// Play Hit Sound
			playerAudio.Data = audios.PlayerHit_mp3
			playerAudio.Playing = true

			// Player dies
			player.IsDead = true
			components.Player.SetValue(playerEntry, *player)

			// Set dead animation
			playerSprite.AnimationName = "dead"
			playerSprite.Animation = playerDeathAnim // Adjust speed as needed
			components.Sprite.SetValue(playerEntry, *playerSprite)

			audioQuery := donburi.NewQuery(filter.Contains(components.Audio))
			audioQuery.Each(ecs.World, func(audio *donburi.Entry) {
				audioComp := components.Audio.Get(audio)
				audioComp.Playing = false
			})

			// Remove player's velocity to stop movement
			velocity := components.Velocity.Get(playerEntry)
			if velocity != nil {
				velocity.X = 0
				velocity.Y = 0
				components.Velocity.SetValue(playerEntry, *velocity)
			}

			deadChan := make(chan bool)                                        // Crea un canal para recibir las señales de la tarea
			go removePlayerFormWorld(ecs, playerEntry, *playerAudio, deadChan) // Inicia la tarea

			select {
			case <-deadChan:
				fmt.Println("Game Over!")
				playerAudio.Data = nil
				playerAudio.Playing = false
			default:
			}

			// Game ends (you might want to implement a proper game state manager)
			return true // Stop checking for further collisions as the player is dead
		},
	})

	fmt.Printf("Player Collision detected: %v\n", collision) // Added log
}

func removePlayerFormWorld(ecs *ecs.ECS, playerEntry *donburi.Entry, playerAudio components.AudioComponent, canal chan bool) {
	// Simula una tarea que podría tardar más de 200 milisegundos.
	time.Sleep(800 * time.Millisecond) // O reemplaza con tu lógica

	ecs.World.Remove(playerEntry.Entity())

	canal <- true // Envía una señal cuando la tarea finaliza.
}
