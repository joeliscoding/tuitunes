package audioplayer

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/wav"
)

var paused = false
var volume = &effects.Volume{
	Streamer: nil,
	Base:     2,
	Volume:   0,
	Silent:   false,
}

func AddSong(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	fileExt := strings.ToLower(file[strings.LastIndex(file, "."):])

	switch fileExt {
	case ".mp3":
		err := decodeMP3(f)
		if err != nil {
			return err
		}
	case ".wav":
		err := decodeWAV(f)
		if err != nil {
			return err
		}
	case ".flac":
		err := decodeFLAC(f)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported file type: %s", fileExt)
	}

	return nil
}

func decodeMP3(f *os.File) error {
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}
	defer streamer.Close()
	playAudio(streamer, format)
	return nil
}

func decodeWAV(f *os.File) error {
	streamer, format, err := wav.Decode(f)
	if err != nil {
		return err
	}
	defer streamer.Close()
	playAudio(streamer, format)
	return nil
}

func decodeFLAC(f *os.File) error {
	streamer, format, err := flac.Decode(f)
	if err != nil {
		return err
	}
	defer streamer.Close()
	playAudio(streamer, format)
	return nil
}

func playAudio(streamer beep.StreamSeekCloser, format beep.Format) error {
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}
	volume.Streamer = ctrl
	speaker.Play(volume)

	for {
		speaker.Lock()
		ctrl.Paused = paused
		speaker.Unlock()
	}
}

func PauseAudio() error {
	paused = !paused
	return nil
}

func ChangeVolume(delta float64) {
	if volume.Volume+delta > -15 && volume.Volume+delta < 0 {
		volume.Volume += delta
		fmt.Printf("Volume: %f\n", volume.Volume)
	} else if volume.Volume+delta <= -15 {
		volume.Volume = -15
		fmt.Printf("Volume: %f (min)\n", volume.Volume)
	} else if volume.Volume+delta >= 0 {
		volume.Volume = 0
		fmt.Printf("Volume: %f (max)\n", volume.Volume)
	}
}
