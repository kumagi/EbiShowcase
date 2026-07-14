package uilab

import "testing"

func TestBilingualFacesLoad(t *testing.T) {
	for _, l := range []string{"ja", "en"} {
		f, e := Face(l, 18)
		if e != nil || f == nil {
			t.Fatalf("%s: %v", l, e)
		}
	}
}
