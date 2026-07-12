// Package tracklayout is the shared catalog for the 応用編 (15 genre tracks)
// texture atlas: one sheet of 海老天-themed sprites that every track game can
// SubImage from. Frame size and name order are the single source of truth for
// cmd/gen-track-atlas and internal/trackatlas.
package tracklayout

// Cell size and grid width.
const (
	FrameW = 48
	FrameH = 48
	Cols   = 8
)

// Sprite is one named cell in the atlas.
type Sprite struct {
	Name string
	Row  int
	Col  int
}

// Sprites lists every cell in row-major order. Index = Row*Cols+Col.
var Sprites = []Sprite{
	// Row 0 — characters
	{"hero", 0, 0},
	{"ally", 0, 1},
	{"npc", 0, 2},
	{"pet", 0, 3},
	{"fighter-p1", 0, 4},
	{"fighter-p2", 0, 5},
	{"slime", 0, 6},
	{"king-crab", 0, 7},
	// Row 1 — enemies
	{"ghost-patrol", 1, 0},
	{"ghost-chase", 1, 1},
	{"ghost-search", 1, 2},
	{"scout", 1, 3},
	{"leaf-guard", 1, 4},
	{"slug", 1, 5},
	{"swarm", 1, 6},
	{"boss-crab", 1, 7},
	// Row 2 — species + bomb kit
	{"species-0", 2, 0},
	{"species-1", 2, 1},
	{"species-2", 2, 2},
	{"species-3", 2, 3},
	{"species-evo", 2, 4},
	{"bomb", 2, 5},
	{"flame", 2, 6},
	{"capture-orb", 2, 7},
	// Row 3 — pickups
	{"pearl", 3, 0},
	{"coin", 3, 1},
	{"xp-gem", 3, 2},
	{"power-star", 3, 3},
	{"upgrade-blast", 3, 4},
	{"upgrade-cap", 3, 5},
	{"upgrade-spd", 3, 6},
	{"flag", 3, 7},
	// Row 4 — puzzle gems + props
	{"gem-red", 4, 0},
	{"gem-blue", 4, 1},
	{"gem-yellow", 4, 2},
	{"gem-green", 4, 3},
	{"gem-purple", 4, 4},
	{"gem-trash", 4, 5},
	{"peg", 4, 6},
	{"aura", 4, 7},
	// Row 5 — merge tiers (海老天 stacking)
	{"merge-1", 5, 0},
	{"merge-2", 5, 1},
	{"merge-3", 5, 2},
	{"merge-4", 5, 3},
	{"merge-5", 5, 4},
	{"merge-6", 5, 5},
	{"merge-7", 5, 6},
	{"pulse", 5, 7},
	// Row 6 — terrain
	{"tile-grass", 6, 0},
	{"tile-grass-dark", 6, 1},
	{"tile-cobble", 6, 2},
	{"tile-water", 6, 3},
	{"tile-wall", 6, 4},
	{"tile-crate", 6, 5},
	{"tile-wood", 6, 6},
	{"tile-stone", 6, 7},
	// Row 7 — more tiles + cards
	{"tile-glass", 7, 0},
	{"tile-lantern", 7, 1},
	{"tile-exit", 7, 2},
	{"tile-platform", 7, 3},
	{"tile-cell", 7, 4},
	{"card-attack", 7, 5},
	{"card-block", 7, 6},
	{"card-skill", 7, 7},
	// Row 8 — UI + bakery
	{"ui-btn", 8, 0},
	{"ui-btn-accent", 8, 1},
	{"ui-panel", 8, 2},
	{"ui-modal", 8, 3},
	{"route-rest", 8, 4},
	{"route-treasure", 8, 5},
	{"block-cell", 8, 6},
	{"bakery", 8, 7},
}

// SheetW and SheetH are full atlas pixel size.
func SheetW() int { return Cols * FrameW }
func SheetH() int {
	rows := 0
	for _, s := range Sprites {
		if s.Row+1 > rows {
			rows = s.Row + 1
		}
	}
	return rows * FrameH
}

// Rect returns the pixel rectangle of one cell.
func Rect(row, col int) (x, y, w, h int) {
	return col * FrameW, row * FrameH, FrameW, FrameH
}

// Find returns the sprite with the given name.
func Find(name string) (Sprite, bool) {
	for _, s := range Sprites {
		if s.Name == name {
			return s, true
		}
	}
	return Sprite{}, false
}

// Names returns every sprite name in atlas order.
func Names() []string {
	out := make([]string, len(Sprites))
	for i, s := range Sprites {
		out[i] = s.Name
	}
	return out
}
