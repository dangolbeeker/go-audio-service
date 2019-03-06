package mix

import (
	"go-audio-service/snd"
	"sync"
)

type channelStruct struct {
	channel *Channel
	input   <-chan *snd.Samples
	buffer  []snd.Sample
}

// Mixer allows the mixing of different channels
type Mixer struct {
	mtx        sync.Mutex
	channels   []*channelStruct
	samplerate uint32
	gain       float32
	output     snd.Input
	done       chan struct{}
	running    bool
}

// NewMixer creates a new Mixer instance
func NewMixer(samplerate uint32) *Mixer {
	return &Mixer{
		samplerate: samplerate,
		gain:       1.0,
		done:       make(chan struct{}),
		running:    false,
	}
}

func (m *Mixer) addChannel(ch *Channel) chan<- *snd.Samples {
	samplesCh := make(chan *snd.Samples, 20)
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.channels = append(m.channels, &channelStruct{
		channel: ch,
		input:   samplesCh,
	})
	return samplesCh
}

// SetOutput sets the next filter in the output chain
func (m *Mixer) SetOutput(out snd.Input) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.output = out
	if out != nil && !m.running {
		m.startWorker()
		m.running = true
	}
}

// SetGain sets the master gain value
func (m *Mixer) SetGain(gain float32) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.gain = gain
}

// Gain returns the master gain value
func (m *Mixer) Gain() float32 {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return m.gain
}

// Stop stops the mixer
func (m *Mixer) Stop() {
	if m.running {
		m.done <- struct{}{}
	}
}

// GetChannel returns a new channel connected to this Mixer
func (m *Mixer) GetChannel() *Channel {
	ch := NewChannel(m.samplerate)
	ch.out = m.addChannel(ch)
	return ch
}

const bufSize = 256

func (m *Mixer) startWorker() {
	go func() {
		buffer := make([]snd.Sample, bufSize)
		samples := &snd.Samples{
			SampleRate: m.samplerate,
			Frames:     buffer,
		}
		sampleIdx := 0
		for {
			select {
			case <-m.done:
				m.mtx.Lock()
				defer m.mtx.Unlock()
				m.running = false
				return
			default:
			}

			m.mtx.Lock()
			for _, channel := range m.channels {
				if len(channel.buffer) < bufSize {
					select {
					case newSamples := <-(channel.input):
						channel.buffer = append(channel.buffer, newSamples.Frames...)
					default:
					}
				}
			}

			for _, channel := range m.channels {
				if len(channel.buffer) > 0 {
					buffer[sampleIdx].L += channel.buffer[0].L
					buffer[sampleIdx].R += channel.buffer[0].R
					channel.buffer = channel.buffer[1:]
				}
			}
			sampleIdx++
			if sampleIdx == bufSize {
				err := m.output.Write(samples)
				for i := 0; i < bufSize; i++ {
					buffer[i].L = 0.0
					buffer[i].R = 0.0
				}
				if err != nil {
					m.mtx.Unlock()
					panic(err)
				}
				sampleIdx = 0
			}
			m.mtx.Unlock()
		}
	}()
}
