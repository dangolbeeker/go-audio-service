package generators

import (
	"go-audio-service/snd"
)

type Rect struct {
	snd.BasicWritableProvider
	samplerate uint32
	freq       float32
	fm         *snd.BasicConnector
	am         *snd.BasicConnector
}

func NewRect(samplerate uint32, freq float32) *Rect {
	r := &Rect{
		samplerate: samplerate,
		freq:       freq,
	}
	r.InitBasicWritableProvider()
	r.fm = r.AddInput("fm", 0.0)
	r.am = r.AddInput("am", 0.0)
	return r
}

func (r *Rect) Read(samples *snd.Samples) {
	r.ReadStateless(samples, 0, snd.EmptyNoteState)
}

func (r *Rect) ReadStateless(samples *snd.Samples, freq float32, state *snd.NoteState) {
	length := len(samples.Frames)
	var v float32
	var max uint32
	if freq > 0 {
		max = uint32(float32(samples.SampleRate) / freq)
	} else {
		max = uint32(float32(samples.SampleRate) / r.freq)
	}

	fm := r.fm.ReadBuffered(samples.SampleRate, len(samples.Frames), 0, state)
	am := r.am.ReadBuffered(samples.SampleRate, len(samples.Frames), 0, state)

	current := state.Timecode % max
	for i := 0; i < length; i++ {
		max += uint32(fm.Frames[i].L)
		half := max / 2
		if state.On {
			if current < half {
				v = 0.5 + am.Frames[i].L
			} else {
				v = -0.5 + am.Frames[i].L
			}
		} else {
			v = 0
		}
		samples.Frames[i].L = v
		samples.Frames[i].R = v
		current++
		if current >= max {
			current = 0
		}
	}
}
