package resolv

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components"
)

var (
	TagPlayer     = resolv.NewTag("Player")
	TagSolidWall  = resolv.NewTag("SolidWall")
	TagObstacle   = resolv.NewTag("Obstacle")
	TagPlatform   = resolv.NewTag("Platform")
	TagProjectile = resolv.NewTag("Projectile")
)

func Add(spaceEntry *donburi.Entry, shapes ...*donburi.Entry) {
	space := components.Space.Get(spaceEntry)
	for _, shapeEntry := range shapes {
		shape := GetShape(shapeEntry)
		space.Add(shape)
	}
}

func SetShape(entry *donburi.Entry, obj resolv.IShape) {
	components.Shape.Set(entry, obj.(*resolv.ConvexPolygon))
}

func SetShapeCircle(entry *donburi.Entry, obj resolv.IShape) {
	components.ShapeCircle.Set(entry, obj.(*resolv.Circle))
}

func GetShape(entry *donburi.Entry) resolv.IShape {
	return components.Shape.Get(entry)
}
