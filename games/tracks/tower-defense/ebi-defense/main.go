package main

import (
	_ "embed"

	"github.com/kumagi/EbiShowcase/internal/towerdefenseplay"
)

//go:embed assets/pearl-gate-battlefield.png
var battlefieldPNG []byte

//go:embed assets/defense-towers.png
var towersPNG []byte

//go:embed assets/defense-battle.png
var battlePNG []byte

func main() {
	towerdefenseplay.Run(towerdefenseplay.Config{Step: 8, Title: "EBI PEARL DEFENSE", BackgroundPNG: battlefieldPNG, TowersPNG: towersPNG, BattlePNG: battlePNG})
}
