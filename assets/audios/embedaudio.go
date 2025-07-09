package audios

import (
	_ "embed"
)

var (
	//go:embed attack.mp3
	Attack_mp3 []byte

	//go:embed gameover.mp3
	Gameover_mp3 []byte

	//go:embed main_scene.mp3
	MainScene_mp3 []byte

	//go:embed title_scene.mp3
	TitleScene_mp3 []byte

	//go:embed player_hit.mp3
	PlayerHit_mp3 []byte
)
