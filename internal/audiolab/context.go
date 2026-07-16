package audiolab

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

var (
	contextOnce sync.Once
	context     *audio.Context
)

// Context returns the single audio context owned by the running game.
//
// Ebitengine permits at most one audio.Context per process. A completed run
// often replaces its game value with *newGame(); constructing audio there a
// second time would panic and make every retry button appear unresponsive.
// Keeping the process-wide device here lets a new run reset gameplay state
// without attempting to recreate the browser audio device.
func Context() *audio.Context {
	contextOnce.Do(func() {
		context = audio.NewContext(SampleRate)
	})
	return context
}
