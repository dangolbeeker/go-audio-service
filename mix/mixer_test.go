package mix

import (
	"go-audio-service/snd"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMixer(t *testing.T) {
	assert := assert.New(t)
	samples1 := &snd.Samples{
		SampleRate: 22000,
		Frames:     make([]snd.Sample, 5),
	}
	for i := 0; i < 5; i++ {
		samples1.Frames[i].L = 0.3
		samples1.Frames[i].R = 0.3
	}
	samples2 := &snd.Samples{
		SampleRate: 22000,
		Frames:     make([]snd.Sample, 5),
	}
	for i := 0; i < 5; i++ {
		samples2.Frames[i].L = 0.2
		samples2.Frames[i].R = 0.2
	}

	ch1 := NewChannel(22000)
	ch2 := NewChannel(22000)
	m := NewMixer(22000)
	ch1.SetMixer(m)
	ch2.SetMixer(m)
	assert.Nil(ch1.Write(samples1))
	assert.Nil(ch2.Write(samples2))
	buf := &snd.BufferedOutput{}
	m.SetOutput(buf)

	time.Sleep(500 * time.Millisecond)

	assert.True(m.running)
	m.Stop()
	assert.False(m.running)
	assert.Len(buf.Frames, 5)
	assert.Equal(float32(0.5), buf.Frames[0].L)
	assert.Equal(float32(0.5), buf.Frames[0].R)
}
