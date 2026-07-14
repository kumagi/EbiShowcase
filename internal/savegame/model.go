// Package savegame provides a small, versioned save envelope for Ebitengine
// games. Keep game rules out of this package: callers store their own JSON
// payload in Data and migrate it deliberately when a format changes.
package savegame

import (
	"encoding/json"
	"fmt"
	"time"
)

const CurrentVersion = 1

type Model struct {
	Version   int             `json:"version"`
	UpdatedAt time.Time       `json:"updated_at"`
	Data      json.RawMessage `json:"data"`
}

func New(payload any, now time.Time) (Model, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return Model{}, err
	}
	return Model{Version: CurrentVersion, UpdatedAt: now.UTC(), Data: data}, nil
}

func Decode(raw []byte) (Model, error) {
	var model Model
	if err := json.Unmarshal(raw, &model); err != nil {
		return Model{}, err
	}
	if model.Version != CurrentVersion {
		return Model{}, fmt.Errorf("savegame: unsupported version %d", model.Version)
	}
	if !json.Valid(model.Data) {
		return Model{}, fmt.Errorf("savegame: invalid payload")
	}
	return model, nil
}

func (m Model) Into(target any) error { return json.Unmarshal(m.Data, target) }
