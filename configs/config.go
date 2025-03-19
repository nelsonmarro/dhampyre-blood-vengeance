package configs

type GameConfig struct {
	ScreenWidth                 int
	ScreenHeight                int
	TileSize                    int
	PlayerSize                  int
	JumpSpeed                   float64
	JumpAnimationDurationFrames float64
	Gravity                     float64
}

var C *GameConfig

func InitGameConfig() {
	C = &GameConfig{
		ScreenWidth:                 640,
		ScreenHeight:                480,
		TileSize:                    16,
		PlayerSize:                  128,
		JumpSpeed:                   -10,
		JumpAnimationDurationFrames: 10,
		Gravity:                     0.5,
	}
}
