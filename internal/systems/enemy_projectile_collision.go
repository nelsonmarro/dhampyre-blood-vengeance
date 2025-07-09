package systems

import (
	"fmt"
	"time"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components"
	dresolv "github.com/nelsonmarro/dhampyre-blood-vengeance/internal/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/ganim8/v2"
)

func EnemyProjectileCollisionSystemFunc(ecs *ecs.ECS, enemyDeathAnim *ganim8.Animation) {
	projectileQuery := donburi.NewQuery(filter.Contains(components.Projectile, components.Position, components.Shape))
	enemyQuery := donburi.NewQuery(filter.Contains(components.Enemy, components.Position, components.Shape, components.Sprite, components.Velocity))

	spaceQuery := donburi.NewQuery(filter.Contains(components.Space))
	var spaceEntry *donburi.Entry
	spaceEntry, _ = spaceQuery.FirstEntity(ecs.World)

	if spaceEntry == nil {
		fmt.Println("Error: No Space entity found.")
		return
	}
	var projectileEntry *donburi.Entry
	projectileEntry, _ = projectileQuery.FirstEntity(ecs.World)
	if projectileEntry == nil {
		return // No player to check collision for
	}
	proectilePosition := components.Position.Get(projectileEntry)
	proectileShape := components.Shape.Get(projectileEntry)

	if proectilePosition == nil || proectileShape == nil {
		return // proectile missing position or shape
	}

	proectileTestShape := proectileShape.Clone()
	proectileTestShape.SetPosition(proectilePosition.X-float64(proectileShape.Bounds().Width()/2+20), proectilePosition.Y-float64(proectileShape.Bounds().Height()/2)) // Adjust position if needed

	enemyQuery.EachEntity(ecs.World, func(entry *donburi.Entry) {
		enemy := components.Enemy.Get(entry)
		enemySprite := components.Sprite.Get(entry)
		enemyShape := components.Shape.Get(entry)
		enemyVelocity := components.Velocity.Get(entry)

		// check for collisions with proectile
		proectileCollision := proectileTestShape.Intersection(enemyShape)

		if !proectileCollision.IsEmpty() {
			enemy.IsDead = true
			enemySprite.AnimationName = "dead"
			enemySprite.Animation = enemyDeathAnim // Adjust speed as needed
			components.Sprite.SetValue(entry, *enemySprite)
			if enemyVelocity != nil {
				enemyVelocity.X = 0
				enemyVelocity.Y = 0
				components.Velocity.SetValue(entry, *enemyVelocity)
			}

			deadChan := make(chan bool)                               // Crea un canal para recibir las señales de la tarea
			go removeEnemyFormWorld(ecs, spaceEntry, entry, deadChan) // Inicia la tarea

			select {
			case <-deadChan:
				fmt.Println("EnemyDead")
				proectileTestShape.SetPosition(proectilePosition.X-float64(proectileShape.Bounds().Width()/2+20), proectilePosition.Y-float64(proectileShape.Bounds().Height()/2)) // Adjust position if needed
			default:
			}

		}
	})
}

func removeEnemyFormWorld(ecs *ecs.ECS, space, entry *donburi.Entry, canal chan bool) {
	// Simula una tarea que podría tardar más de 200 milisegundos.
	time.Sleep(700 * time.Millisecond) // O reemplaza con tu lógica

	dresolv.Remove(space, entry)
	ecs.World.Remove(entry.Entity())

	canal <- true // Envía una señal cuando la tarea finaliza.
}
