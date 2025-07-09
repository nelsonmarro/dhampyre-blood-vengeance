package components

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/yohamta/donburi"
)

// AudioComponent stores audio-related data for an entity.
type AudioComponent struct {
	Data         []byte
	Format       string // e.g., "mp3", "wav"
	Player       *audio.Player
	Loop         bool
	Playing      bool
	Volume       float64
	SampleRate   int
	InfiniteLoop *audio.InfiniteLoop
}

// Audio is the component key for AudioComponent.
var Audio = donburi.NewComponentType[AudioComponent]()
