package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"github.com/kumagi/EbiShowcase/internal/cameralab"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"github.com/kumagi/EbiShowcase/internal/uilab"
)

const (
	width, height = 480, 720
	tide          = 0
	flame         = 1
	leaf          = 2
	modeMap       = 0
	modeBattle    = 1
	partyLimit    = 3
	timeLimit     = 100 * 60
)

type species struct {
	name  string
	kind  int
	maxHP int
	color color.RGBA
}

var speciesBook = []species{
	{"REEFLET", tide, 60, color.RGBA{72, 173, 218, 255}},
	{"MOSSHELL", leaf, 48, color.RGBA{91, 184, 106, 255}},
	{"CINDERFIN", flame, 50, color.RGBA{230, 102, 69, 255}},
	{"CLOUDRAY", tide, 52, color.RGBA{119, 180, 224, 255}},
	{"REEF LORD", tide, 82, color.RGBA{190, 105, 216, 255}},
}

type moveData struct {
	name     string
	power    int
	accuracy int
	maxUses  int
}

var moves = []moveData{{"QUICK FIN", 10, 100, 6}, {"AFFINITY BURST", 18, 80, 3}}
var typeNames = []string{"TIDE", "FLAME", "LEAF"}

var matchup = [3][3]float64{
	{1, 2, 0.5},
	{0.5, 1, 2},
	{2, 0.5, 1},
}

var regionNames = []string{"TIDEPOOL", "EMBER COVE", "KELP FOREST"}
var encounterTables = [3][2]int{{1, 0}, {2, 1}, {3, 2}}

type monster struct {
	speciesID int
	hp        int
	exp       int
}

type saveData struct {
	PartySpecies []int   `json:"party_species"`
	PartyExp     []int   `json:"party_exp"`
	BoxSpecies   []int   `json:"box_species"`
	Dex          [4]bool `json:"dex"`
	Visits       [3]int  `json:"region_visits"`
	Badges       [3]bool `json:"region_badges"`
	BestFrames   int     `json:"best_frames"`
}

var saveSlot []byte

type sparkle struct {
	x, y, vx, vy float64
	life         int
	c            color.RGBA
}

type game struct {
	party               []monster
	box                 []monster
	dex                 [4]bool
	visits              [3]int
	badges              [3]bool
	currentRegion       int
	active              int
	mode                int
	wildSpecies, wildHP int
	moveUses            [2]int
	orbs                int
	turns, frames       int
	mustSwitch          bool
	clear, over         bool
	rng                 *rand.Rand
	message             string
	savePreview         string
	actionFrame         int
	actionKind          int
	hitFlash, shake     int
	pendingCapture      bool
	bestFrames          int
	sparkles            []sparkle
	scene               *ebiten.Image
	audio               *audio.Context
	gate                audiolab.Gate
	pulse               *shaderlab.Pulse
	cam                 cameralab.State
	badge               *ebiten.Image
}

func newGame() *game {
	g := &game{
		party:   []monster{{speciesID: 0, hp: speciesBook[0].maxHP}},
		orbs:    8,
		rng:     rand.New(rand.NewSource(8707)),
		message: "Start at EMBER COVE: regional tables decide who appears.",
		scene:   ebiten.NewImage(width, height),
	}
	g.dex[0] = true
	g.audio = audio.NewContext(audiolab.SampleRate)
	g.pulse = shaderlab.NewPulse()
	g.cam = cameralab.State{Pos: cameralab.Vec{X: width / 2, Y: height / 2}, ViewW: width, ViewH: height}
	g.badge = ebiten.NewImage(20, 20)
	g.badge.Fill(color.RGBA{255, 222, 92, 255})
	if len(saveSlot) > 0 {
		g.savePreview = fmt.Sprintf("SAVE SLOT FOUND: %d bytes", len(saveSlot))
	}
	return g
}

func (g *game) Update() error {
	for i := len(g.sparkles) - 1; i >= 0; i-- {
		p := &g.sparkles[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += 0.05
		p.life--
		if p.life <= 0 {
			g.sparkles = append(g.sparkles[:i], g.sparkles[i+1:]...)
		}
	}
	if g.hitFlash > 0 {
		g.hitFlash--
	}
	if g.shake > 0 {
		g.shake--
	}
	if g.actionFrame > 0 {
		g.actionFrame--
		if g.actionFrame == 0 && g.pendingCapture {
			g.finishCapture()
		}
		return nil
	}
	if g.clear || g.over {
		if retryPressed() {
			best := g.bestFrames
			*g = *newGame()
			g.bestFrames = best
		}
		return nil
	}
	g.frames++
	if g.frames >= timeLimit {
		g.over = true
		g.message = "Expedition time ended before the bestiary was saved."
		return nil
	}
	choice := -1
	for i, key := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4} {
		if inpututil.IsKeyJustPressed(key) {
			choice = i
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) && g.mode == modeMap {
		g.save()
		return nil
	}
	if x, y, ok := pressPosition(); ok {
		if g.mode == modeMap {
			if y >= 540 && y < 620 {
				choice = min(2, x/160)
			} else if y >= 630 {
				g.save()
				return nil
			}
		} else if y >= 540 {
			row, col := (y-540)/78, x/240
			choice = row*2 + col
		}
	}
	if g.mode == modeMap {
		if choice >= 0 && choice < 3 {
			g.explore(choice)
		}
		return nil
	}
	if choice >= 0 && choice < 4 {
		g.battleAction(choice)
	}
	return nil
}

func (g *game) explore(region int) {
	entry := g.visits[region] % len(encounterTables[region])
	g.wildSpecies = encounterTables[region][entry]
	g.visits[region]++
	g.wildHP = 40
	g.moveUses = [2]int{moves[0].maxUses, moves[1].maxUses}
	g.mode = modeBattle
	g.currentRegion = region
	g.turns = 0
	g.mustSwitch = false
	g.message = fmt.Sprintf("%s encounter: wild %s! Weaken, then use an orb.", regionNames[region], speciesBook[g.wildSpecies].name)
}

func (g *game) battleAction(choice int) {
	if g.mustSwitch && choice != 3 {
		g.message = "The front fainted. SWITCH to a living party member first."
		return
	}
	switch choice {
	case 0, 1:
		g.useMove(choice)
	case 2:
		g.throwOrb()
	case 3:
		g.switchParty()
	}
}

func (g *game) useMove(index int) {
	if g.moveUses[index] == 0 {
		g.message = moves[index].name + " has no uses left in this encounter."
		return
	}
	g.moveUses[index]--
	g.turns++
	front := &g.party[g.active]
	data := moves[index]
	roll := g.rng.Intn(100)
	if roll >= data.accuracy {
		g.message = fmt.Sprintf("%s missed: roll %d >= %d.", data.name, roll, data.accuracy)
	} else {
		attackerType := speciesBook[front.speciesID].kind
		defenderType := speciesBook[g.wildSpecies].kind
		multiplier := matchup[attackerType][defenderType]
		damage := int(float64(data.power) * multiplier)
		g.wildHP = max(1, g.wildHP-damage)
		g.play(720)
		front.exp += 3
		g.message = fmt.Sprintf("%s x%.1f dealt %d. Wild HP %d/40.", data.name, multiplier, damage, g.wildHP)
		g.hitFlash = 12
		g.shake = 7
		g.burst(355, 235, speciesBook[g.wildSpecies].color, 18)
	}
	g.wildCounter()
	g.actionFrame = 24
	g.actionKind = 1
}

func (g *game) throwOrb() {
	if g.orbs == 0 {
		g.message = "No capture orbs remain."
		return
	}
	g.orbs--
	g.play(440)
	g.turns++
	chance := min(95, 25+(40-g.wildHP)*2)
	roll := g.rng.Intn(100)
	if roll >= chance {
		g.message = fmt.Sprintf("Capture failed: roll %d >= %d%%.", roll, chance)
		g.wildCounter()
		g.actionFrame = 38
		g.actionKind = 2
		g.burst(355, 235, color.RGBA{240, 185, 70, 255}, 12)
		return
	}
	g.pendingCapture = true
	g.actionFrame = 54
	g.actionKind = 2
	g.message = "The capture orb rocks... hold your breath!"
	return
}

func (g *game) finishCapture() {
	g.pendingCapture = false
	g.play(900)
	captured := monster{speciesID: g.wildSpecies, hp: speciesBook[g.wildSpecies].maxHP}
	if len(g.party) < partyLimit {
		g.party = append(g.party, captured)
		g.message = speciesBook[g.wildSpecies].name + " captured into the party!"
	} else {
		g.box = append(g.box, captured)
		g.message = speciesBook[g.wildSpecies].name + " captured into the storage box!"
	}
	if g.wildSpecies < len(g.dex) {
		g.dex[g.wildSpecies] = true
	}
	if g.wildSpecies == encounterTables[g.currentRegion][0] {
		g.badges[g.currentRegion] = true
		g.message += " Region research badge earned!"
	}
	g.burst(355, 235, color.RGBA{255, 222, 92, 255}, 42)
	g.shake = 12
	for i := range g.party {
		g.party[i].exp += 20
		g.evolveIfReady(i)
	}
	g.mode = modeMap
	g.active = min(g.active, len(g.party)-1)
	if g.badgeCount() == 3 && g.starterEvolved() {
		g.message += " Three regions complete: press SAVE to finish."
	} else {
		g.message += " EXP shared; check growth and choose the next region."
	}
	if g.orbs == 0 && g.dexCount() < 4 {
		g.over = true
		g.message = "Every orb was used before all four species were registered."
	}
}

func (g *game) play(freq float64) {
	g.gate.Arm(true)
	g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Sine, freq, .065)).Play()
}

func (g *game) burst(x, y float64, c color.RGBA, amount int) {
	for i := 0; i < amount; i++ {
		a := g.rng.Float64() * math.Pi * 2
		speed := 0.8 + g.rng.Float64()*3
		g.sparkles = append(g.sparkles, sparkle{x, y, math.Cos(a) * speed, math.Sin(a)*speed - 0.7, 20 + g.rng.Intn(24), c})
	}
}

func (g *game) badgeCount() int {
	count := 0
	for _, earned := range g.badges {
		if earned {
			count++
		}
	}
	return count
}

func (g *game) wildCounter() {
	if g.wildHP <= 0 {
		return
	}
	front := &g.party[g.active]
	wildType := speciesBook[g.wildSpecies].kind
	frontType := speciesBook[front.speciesID].kind
	multiplier := matchup[wildType][frontType]
	damage := int(8 * multiplier)
	front.hp = max(0, front.hp-damage)
	g.message += fmt.Sprintf(" Counter x%.1f dealt %d.", multiplier, damage)
	if front.hp == 0 {
		if g.aliveCount() == 0 {
			g.over = true
			g.message = "The entire field party fainted."
			return
		}
		g.mustSwitch = true
		g.message += " Front fainted: SWITCH is required."
	}
}

func (g *game) switchParty() {
	if len(g.party) < 2 {
		g.message = "No reserve partner has been captured yet."
		return
	}
	for step := 1; step <= len(g.party); step++ {
		next := (g.active + step) % len(g.party)
		if g.party[next].hp > 0 {
			g.active = next
			if g.mustSwitch {
				g.mustSwitch = false
				g.message = "Forced switch completed without another counter."
				return
			}
			g.turns++
			g.message = "Planned switch spent a turn."
			g.wildCounter()
			return
		}
	}
}

func (g *game) evolveIfReady(index int) {
	m := &g.party[index]
	if m.speciesID == 0 && level(m.exp) >= 3 {
		m.speciesID = 4
		m.hp = max(m.hp, speciesBook[4].maxHP/2)
		g.message += " REEFLET evolved into REEF LORD!"
	}
}

func level(exp int) int { return 1 + exp/30 }

func (g *game) starterEvolved() bool {
	for _, m := range g.party {
		if m.speciesID == 4 {
			return true
		}
	}
	return false
}

func (g *game) aliveCount() int {
	count := 0
	for _, m := range g.party {
		if m.hp > 0 {
			count++
		}
	}
	return count
}

func (g *game) dexCount() int {
	count := 0
	for _, caught := range g.dex {
		if caught {
			count++
		}
	}
	return count
}

func (g *game) save() {
	if g.badgeCount() < 3 || !g.starterEvolved() {
		g.message = "SAVE locked: earn 3 region badges and evolve REEFLET first."
		return
	}
	if g.bestFrames == 0 || g.frames < g.bestFrames {
		g.bestFrames = g.frames
	}
	payload := saveData{Dex: g.dex, Visits: g.visits, Badges: g.badges, BestFrames: g.bestFrames}
	for _, m := range g.party {
		payload.PartySpecies = append(payload.PartySpecies, m.speciesID)
		payload.PartyExp = append(payload.PartyExp, m.exp)
	}
	for _, m := range g.box {
		payload.BoxSpecies = append(payload.BoxSpecies, m.speciesID)
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		g.message = "Could not serialize save data."
		return
	}
	saveSlot = append(saveSlot[:0], encoded...)
	g.savePreview = fmt.Sprintf("SAVE SLOT: %d bytes / PARTY %d / BOX %d", len(saveSlot), len(g.party), len(g.box))
	g.clear = true
	g.message = "Bestiary, party, box, EXP, and region visits serialized!"
}

func (g *game) Draw(screen *ebiten.Image) {
	g.scene.Clear()
	g.drawScene(g.scene)
	op := &ebiten.DrawImageOptions{}
	if g.shake > 0 {
		op.GeoM.Translate(float64((g.frames%3-1)*3), float64(((g.frames/2)%3-1)*2))
	}
	screen.DrawImage(g.scene, op)
}

func (g *game) drawScene(screen *ebiten.Image) {
	screen.Fill(color.RGBA{9, 20, 39, 255})
	g.drawTitle(screen)
	g.drawEffectBadge(screen)
	best := "--"
	if g.bestFrames > 0 {
		best = fmt.Sprintf("%02d", g.bestFrames/60)
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("BADGES %d/3 DEX %d/4 PARTY %d/3 ORBS %d TIME %02d BEST %s", g.badgeCount(), g.dexCount(), len(g.party), g.orbs, max(0, 100-g.frames/60), best), 42, 44)
	g.drawRoster(screen)
	if g.mode == modeMap {
		g.drawMap(screen)
	} else {
		g.drawBattle(screen)
	}
	if g.clear || g.over {
		title := "EXPEDITION SAVED!"
		if g.over {
			title = "EXPEDITION FAILED!"
		}
		vector.DrawFilledRect(screen, 38, 260, 404, 180, color.RGBA{5, 13, 28, 247}, false)
		ebitenutil.DebugPrintAt(screen, title, 157, 304)
		ebitenutil.DebugPrintAt(screen, g.message, 58, 342)
		ebitenutil.DebugPrintAt(screen, g.savePreview, 92, 374)
		ebitenutil.DebugPrintAt(screen, "TAP / ENTER TO RESTART", 139, 409)
	}
	for _, p := range g.sparkles {
		alpha := uint8(min(255, p.life*10))
		vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 3, color.RGBA{p.c.R, p.c.G, p.c.B, alpha}, false)
	}
}

func (g *game) drawTitle(screen *ebiten.Image) {
	const label = "EBI MONSTERS EXPEDITION"
	if face, err := uilab.Face("en", 16); err == nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(160, 6)
		text.Draw(screen, label, face, op)
		return
	}
	ebitenutil.DebugPrintAt(screen, label, 160, 18)
}

func (g *game) drawEffectBadge(screen *ebiten.Image) {
	if g.pulse == nil || !g.pulse.Available() {
		return
	}
	fx := ebiten.NewImage(20, 20)
	if !g.pulse.Draw(fx, g.badge, float32(g.frames)*.08) {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(width-34, 10)
	screen.DrawImage(fx, op)
}

func (g *game) drawRoster(screen *ebiten.Image) {
	for i := 0; i < partyLimit; i++ {
		x := 8 + i*158
		fill := color.RGBA{34, 55, 75, 255}
		label := "EMPTY PARTY SLOT"
		detail := "—"
		if i < len(g.party) {
			m := g.party[i]
			s := speciesBook[m.speciesID]
			fill = color.RGBA{s.color.R / 2, s.color.G / 2, s.color.B / 2, 255}
			label = s.name
			detail = fmt.Sprintf("%s HP%d LV%d", typeNames[s.kind], m.hp, level(m.exp))
			if i == g.active && g.mode == modeBattle {
				label = "FRONT: " + label
			}
		}
		vector.DrawFilledRect(screen, float32(x), 70, 148, 62, fill, false)
		ebitenutil.DebugPrintAt(screen, label, x+8, 82)
		ebitenutil.DebugPrintAt(screen, detail, x+8, 108)
	}
}

func (g *game) drawMap(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "REGIONAL ENCOUNTER TABLES", 166, 158)
	for i, name := range regionNames {
		x := 12 + i*156
		vector.DrawFilledRect(screen, float32(x), 190, 144, 250, color.RGBA{28, 53, 70, 255}, false)
		ebitenutil.DebugPrintAt(screen, name, x+34, 209)
		first := speciesBook[encounterTables[i][0]].name
		second := speciesBook[encounterTables[i][1]].name
		ebitenutil.DebugPrintAt(screen, "TABLE", x+54, 247)
		ebitenutil.DebugPrintAt(screen, "A: "+first, x+13, 278)
		ebitenutil.DebugPrintAt(screen, "B: "+second, x+13, 306)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("VISITS %d", g.visits[i]), x+39, 344)
		badge := "BADGE: --"
		if g.badges[i] {
			badge = "BADGE: STAR"
		}
		ebitenutil.DebugPrintAt(screen, badge, x+33, 362)
		next := encounterTables[i][g.visits[i]%2]
		ebitenutil.DebugPrintAt(screen, "NEXT", x+55, 386)
		ebitenutil.DebugPrintAt(screen, speciesBook[next].name, x+34, 411)
	}
	ebitenutil.DebugPrintAt(screen, g.message, 43, 477)
	for i, name := range regionNames {
		x := i * 160
		vector.DrawFilledRect(screen, float32(x+8), 540, 144, 75, color.RGBA{52, 91, 122, 255}, false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d %s", i+1, name), x+25, 572)
	}
	ready := g.badgeCount() == 3 && g.starterEvolved()
	fill := color.RGBA{55, 61, 75, 255}
	if ready {
		fill = color.RGBA{75, 166, 111, 255}
	}
	vector.DrawFilledRect(screen, 25, 635, 430, 58, fill, false)
	ebitenutil.DebugPrintAt(screen, "SAVE BESTIARY + PARTY + BOX [S]", 130, 659)
	if g.savePreview != "" {
		ebitenutil.DebugPrintAt(screen, g.savePreview, 100, 704)
	}
}

func (g *game) drawBattle(screen *ebiten.Image) {
	wild := speciesBook[g.wildSpecies]
	front := g.party[g.active]
	frontSpecies := speciesBook[front.speciesID]
	regionColors := [...]color.RGBA{{28, 70, 91, 255}, {83, 45, 46, 255}, {33, 76, 55, 255}}
	vector.DrawFilledRect(screen, 20, 150, 440, 305, regionColors[g.currentRegion], false)
	bob := math.Sin(float64(g.frames)*0.12) * 4
	frontX, wildX := 125.0, 355.0
	if g.actionKind == 1 && g.actionFrame > 10 {
		frontX += float64(24-g.actionFrame) * 4
	}
	trackatlas.DrawCentered(screen, trackatlas.Species(front.speciesID), frontX, 260+bob, 108)
	wildSize := 124.0
	if g.hitFlash > 0 && g.hitFlash%4 < 2 {
		wildSize = 136
	}
	trackatlas.DrawCentered(screen, trackatlas.Species(g.wildSpecies), wildX, 235-bob, wildSize)
	if g.actionKind == 2 && g.actionFrame > 0 {
		progress := 1 - float64(g.actionFrame)/54
		if progress < 0 {
			progress = 0
		}
		orbX := 130 + (355-130)*progress
		orbY := 260 - math.Sin(progress*math.Pi)*125
		vector.DrawFilledCircle(screen, float32(orbX), float32(orbY), 12, color.RGBA{250, 196, 70, 255}, false)
		vector.StrokeCircle(screen, float32(orbX), float32(orbY), 12, 3, color.RGBA{255, 245, 195, 255}, false)
	}
	ebitenutil.DebugPrintAt(screen, "FRONT "+frontSpecies.name, 75, 327)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s HP%d LV%d", typeNames[frontSpecies.kind], front.hp, level(front.exp)), 82, 351)
	ebitenutil.DebugPrintAt(screen, "WILD "+wild.name, 313, 315)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s HP%d/40", typeNames[wild.kind], g.wildHP), 320, 341)
	chance := min(95, 25+(40-g.wildHP)*2)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("CAPTURE %d%%   TURN %d   ORBS %d", chance, g.turns, g.orbs), 142, 390)
	ebitenutil.DebugPrintAt(screen, g.message, 43, 426)
	labels := []string{
		fmt.Sprintf("[1] QUICK %d/%d", g.moveUses[0], moves[0].maxUses),
		fmt.Sprintf("[2] BURST %d/%d", g.moveUses[1], moves[1].maxUses),
		"[3] CAPTURE ORB",
		"[4] SWITCH FRONT",
	}
	for i, label := range labels {
		row, col := i/2, i%2
		x, y := col*240+6, row*78+540
		fill := color.RGBA{51, 84, 122, 255}
		if i == 2 {
			fill = color.RGBA{185, 126, 57, 255}
		}
		if g.mustSwitch && i != 3 {
			fill = color.RGBA{53, 55, 66, 255}
		}
		vector.DrawFilledRect(screen, float32(x), float32(y), 228, 68, fill, false)
		ebitenutil.DebugPrintAt(screen, label, x+35, y+28)
	}
}

func pressPosition() (int, int, bool) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return x, y, true
	}
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		return x, y, true
	}
	return 0, 0, false
}

func retryPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Ebi Monsters — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
