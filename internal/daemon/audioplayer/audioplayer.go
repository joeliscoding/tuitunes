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

var running = false
var paused = false
var volume = &effects.Volume{
	Streamer: nil,
	Base:     2,
	Volume:   0,
	Silent:   false,
}

var ctrl *beep.Ctrl
var Shutdown func() // stops daemon when playback finishes

var queue []string
var queueIndex int = 0

func Enqueue(file string) error {
	queue = append(queue, file)
	fmt.Printf("Added to queue: %s\n", file)

	// start playback if not already running
	if !running {
		running = true
		err := playQueue()
		if err != nil {
			return err
		}
	}
	return nil
}

func playQueue() error {
	for queueIndex < len(queue) {
		fmt.Printf("Now playing: %s\n", queue[queueIndex])
		file := queue[queueIndex]
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		fileExt := strings.ToLower(file[strings.LastIndex(file, "."):])

		switch fileExt {
		case ".mp3":
			err := playMP3(f)
			if err != nil {
				return err
			}
		case ".wav":
			err := playWAV(f)
			if err != nil {
				return err
			}
		case ".flac":
			err := playFLAC(f)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported file type: %s", fileExt)
		}

		queueIndex++
	}

	// shutdown daemon process after finishing queue
	if Shutdown != nil {
		Shutdown()
	}

	return nil
}

func playMP3(f *os.File) error {
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}
	defer streamer.Close()

	err = playStream(streamer, format)
	if err != nil {
		return err
	}
	return nil
}

func playWAV(f *os.File) error {
	streamer, format, err := wav.Decode(f)
	if err != nil {
		return err
	}
	defer streamer.Close()

	err = playStream(streamer, format)
	if err != nil {
		return err
	}
	return nil
}

func playFLAC(f *os.File) error {
	streamer, format, err := flac.Decode(f)
	if err != nil {
		return err
	}
	defer streamer.Close()

	err = playStream(streamer, format)
	if err != nil {
		return err
	}
	return nil
}

func playStream(streamer beep.StreamSeekCloser, format beep.Format) error {
	// TODO: add queue streamer, for seamless transitions between tracks in the queue
	// To achieve this, the next song in the queue should always be added to the beep queue streamer

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}
	volume.Streamer = ctrl

	done := make(chan bool)
	speaker.Play(beep.Seq(volume, beep.Callback(func() {
		done <- true
	})))

	<-done

	return nil
}

func TogglePause() {
	paused = !paused
	ctrl.Paused = paused
}

func AdjustVolume(delta float64) {
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
