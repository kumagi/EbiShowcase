package shaderlab

import "testing"

func TestEmbeddedSourcesArePresentAndDocumentUniforms(t *testing.T) {
	for name, source := range map[string][]byte{
		"pulse": pulseSource, "palette": paletteSource, "distort": distortSource, "status": statusSource,
	} {
		if len(source) < 80 {
			t.Fatalf("%s shader source is unexpectedly empty", name)
		}
	}
	if string(paletteSource) == string(statusSource) {
		t.Fatal("distinct shader lessons must not share one source")
	}
}
