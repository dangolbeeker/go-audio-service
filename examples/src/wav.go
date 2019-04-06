package examples

import (
	"go-audio-service/filters"
	"go-audio-service/generators"
	"go-audio-service/notes"
	"go-audio-service/snd"
	"pixelext/nodes"
	"pixelext/services"

	"github.com/faiface/pixel/pixelgl"

	"github.com/faiface/pixel"
)

type wavExample struct {
	nodes.BaseNode
	multi *notes.NoteMultiplexer
	gain  snd.Readable
}

func (w *wavExample) Init() {
	txt := nodes.NewText("txt", "basic")
	txt.SetZeroAlignment(nodes.AlignmentTopLeft)
	txt.SetPos(pixel.V(20, 580))
	txt.Printf("Wav example")
	txt.Printf("\nPress space for sound\nPress Q to quit")
	w.AddChild(txt)

	samples, err := services.ResourceManager().LoadSample("samples/CYCdh_K4-Snr05.wav")
	if err != nil {
		panic(err)
	}
	sampleplayer := generators.NewSample(samples)

	w.multi = notes.NewNoteMultiplexer()
	w.multi.SetReadable(sampleplayer)

	gain := filters.NewGain(.6)
	gain.SetReadable(w.multi)
	w.gain = gain
}

func (w *wavExample) Mount() {
	GetOutput().SetReadable(w.gain)
	Start()
}

func (w *wavExample) Unmount() {
	Stop()
}

func (w *wavExample) Update(dt float64) {
	if nodes.Events().JustPressed(pixelgl.KeyQ) {
		SwitchScene("main")
	}
	if nodes.Events().JustPressed(pixelgl.KeySpace) {
		w.multi.SendNoteEvent(notes.NewNoteEvent(notes.Pressed, 880.0, .5))
	}
	if nodes.Events().JustReleased(pixelgl.KeySpace) {
		w.multi.SendNoteEvent(notes.NewNoteEvent(notes.Released, 880.0, .5))
	}
}

func init() {
	w := &wavExample{
		BaseNode: *nodes.NewBaseNode("wav"),
	}
	w.Self = w
	AddExample("Wav", w)
}
