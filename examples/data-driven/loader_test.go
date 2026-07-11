package datadriven

import "testing"

func TestLoadVillage(t *testing.T) {
	data, err := LoadMap("village")
	if err != nil {
		t.Fatal(err)
	}
	if data.ID != "village" {
		t.Fatalf("ID = %q, want village", data.ID)
	}
	if data.OnEnter == "" {
		t.Fatal("village must define an onEnter event")
	}
}
