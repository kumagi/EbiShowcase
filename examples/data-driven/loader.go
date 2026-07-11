package datadriven

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed assets/**/*.json
var assets embed.FS

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type MapData struct {
	ID      string `json:"id"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Tiles   []int  `json:"tiles"`
	Spawn   Point  `json:"spawn"`
	OnEnter string `json:"onEnter"`
}

func LoadMap(id string) (MapData, error) {
	path := "assets/maps/" + id + ".json"
	b, err := assets.ReadFile(path)
	if err != nil {
		return MapData{}, fmt.Errorf("read map %q: %w", id, err)
	}

	var data MapData
	if err := json.Unmarshal(b, &data); err != nil {
		return MapData{}, fmt.Errorf("decode map %q: %w", id, err)
	}
	if len(data.Tiles) != data.Width*data.Height {
		return MapData{}, fmt.Errorf("map %q: got %d tiles, want %d", id, len(data.Tiles), data.Width*data.Height)
	}
	return data, nil
}
