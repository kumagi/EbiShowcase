package savegame

import (
	"testing"
	"time"
)

func TestRoundTripAndCopyingStore(t *testing.T) {
	store := NewStore()
	now := time.Date(2026, 7, 14, 12, 0, 0, 0, time.UTC)
	if err := Autosave(store, "slot", struct {
		HP int `json:"hp"`
	}{7}, now); err != nil {
		t.Fatal(err)
	}
	raw, ok, err := store.Load("slot")
	if err != nil || !ok {
		t.Fatalf("load: %v %v", ok, err)
	}
	model, err := Decode(raw)
	if err != nil || model.UpdatedAt != now {
		t.Fatalf("model: %#v %v", model, err)
	}
	var payload struct {
		HP int `json:"hp"`
	}
	if err := model.Into(&payload); err != nil || payload.HP != 7 {
		t.Fatalf("payload: %#v %v", payload, err)
	}
	raw[0] = 'x'
	again, _, _ := store.Load("slot")
	if again[0] == 'x' {
		t.Fatal("store leaked mutable bytes")
	}
}
func TestRejectsUnknownVersion(t *testing.T) {
	if _, err := Decode([]byte(`{"version":99,"data":{}}`)); err == nil {
		t.Fatal("expected version error")
	}
}
