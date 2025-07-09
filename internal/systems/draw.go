package systems

import (
	"image"
	"image/color"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/configs"
	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/ganim8/v2"
)

// Pre-allocate a fallback image to avoid creating it every frame
var fallbackImage *ebiten.Image

func DrawSystemFunc(ecs *ecs.ECS, screen *ebiten.Image) {
	world := ecs.World

	query := donburi.NewQuery(filter.Contains(components.Position, components.Sprite))

	// Initialize the fallback image if it hasn't been already
	if fallbackImage == nil {
		rect := image.Rect(0, 0, configs.C.PlayerSize, configs.C.PlayerSize)
		fallbackImage = ebiten.NewImageFromImage(rect)
		// You might want to fill it with a default color
		fallbackImage.Fill(color.White)
	}

	query.EachEntity(world, func(entry *donburi.Entry) {
		positionComponent := components.Position.Get(entry)
		spriteComponent := components.Sprite.Get(entry)

		if spriteComponent.Animation != nil {
			animWidth, animHeight := spriteComponent.Animation.Size()
			x := positionComponent.X
			y := positionComponent.Y
			ox := 0.5 // Horizontal center
			oy := 0.5 // Vertical center
			angle := 0.0
			sx := 1.0
			sy := 1.0

			if spriteComponent.Flipped {
				sx = -1.0
				x += float64(animWidth)
			}

			x += float64(animWidth) * ox * (sx - 1)
			y += float64(animHeight) * oy * (sy - 1)

			spriteComponent.Animation.Draw(screen, ganim8.DrawOpts(x, y, angle, sx, sy, ox, oy))
		} else {
			// Draw the pre-allocated fallback image
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(positionComponent.X-float64(configs.C.PlayerSize)/2, positionComponent.Y-float64(configs.C.PlayerSize)/2)
			screen.DrawImage(fallbackImage, op)
		}
	})
}
