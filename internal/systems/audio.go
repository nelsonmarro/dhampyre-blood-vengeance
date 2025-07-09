package systems

import (
	"bytes"
	"io"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"

	"github.com/nelsonmarro/dhampyre-blood-vengeance/internal/components"
)

const (
	sampleRate = 48000
)

type audioStream interface {
	io.ReadSeeker
	Length() int64
}

func AudioSystemFunc(e *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(components.Audio)) // Simplified query
	query.Each(e.World, func(entry *donburi.Entry) {
		audioComp := components.Audio.Get(entry)

		if audioComp.Data == nil || audioComp.Format == "" {
			return
		}

		if audioComp.Player == nil {
			audioContext := audio.CurrentContext()
			if audioContext == nil {
				audioContext = audio.NewContext(sampleRate)
			}

			var s audioStream
			var err error

			reader := bytes.NewReader(audioComp.Data)

			switch audioComp.Format {
			case "mp3":
				s, err = mp3.DecodeWithSampleRate(sampleRate, reader)
			case "wav":
				s, err = wav.DecodeWithSampleRate(sampleRate, reader)
			default:
				log.Printf("Unsupported audio format: %s", audioComp.Format)
				return
			}

			if err != nil {
				log.Printf("Error decoding audio (%s): %v", audioComp.Format, err)
				return
			}

			if audioComp.Loop {
				audioComp.InfiniteLoop = audio.NewInfiniteLoop(s, s.Length())
				audioComp.Player, err = audioContext.NewPlayer(audioComp.InfiniteLoop)
			} else {
				audioComp.Player, err = audioContext.NewPlayer(s)
			}

			if err != nil {
				log.Printf("Error creating audio player: %v", err)
				return
			}
			audioComp.Player.SetVolume(audioComp.Volume)
		}

		if audioComp.Playing && !audioComp.Player.IsPlaying() {
			audioComp.Player.Play()
		} else if !audioComp.Playing {
			audioComp.Player.Pause()
			audioComp.Player.Close()
		} else if audioComp.Volume != audioComp.Player.Volume() {
			audioComp.Player.SetVolume(audioComp.Volume)
		}
	})
}
