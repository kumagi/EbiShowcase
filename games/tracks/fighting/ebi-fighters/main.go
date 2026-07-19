package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"github.com/kumagi/EbiShowcase/internal/cameralab"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"github.com/kumagi/EbiShowcase/internal/uilab"
)

const width, height = 480, 720

const (
	attackNone = iota
	attackJab
	attackHeavy
)

const (
	phaseReady = iota
	phaseStartup
	phaseActive
	phaseRecovery
)

const (
	aiIdle = iota
	aiApproach
	aiRetreat
	aiGuard
)

type moveSpec struct {
	name                          string
	activeStart, activeEnd, total int
	damage, chip, guardDamage     int
	hitstun, blockstun            int
	reach, lunge, push            float64
}

var moveSpecs = [...]moveSpec{
	{},
	{
		name: "JAB", activeStart: 8, activeEnd: 12, total: 30,
		damage: 8, guardDamage: 16, hitstun: 14, blockstun: 6,
		reach: 92, lunge: 10, push: 10,
	},
	{
		name: "HEAVY", activeStart: 16, activeEnd: 22, total: 48,
		damage: 20, chip: 2, guardDamage: 38, hitstun: 24, blockstun: 12,
		reach: 128, lunge: 24, push: 20,
	},
}

type fighter struct {
	x                                       float64
	hp, guardHP                             int
	attackKind, attackTick, stun, blockstun int
	guard, attackConnected                  bool
	guardAge, moveDir, walkTick             int
}
type spark struct{ x, y, vx, vy, life float64 }
type game struct {
	p, ai                              fighter
	frame, round, pWins, aiWins        int
	hitstop, shake, streak, bestStreak int
	roundOver, roundWinner             int
	aiPlan, aiPlanTicks                int
	sparks                             []spark
	message                            string
	matchOver                          bool
	rng                                *rand.Rand
	audio                              *audio.Context
	gate                               audiolab.Gate
	pulse                              *shaderlab.Pulse
	cam                                cameralab.State
	badge                              *ebiten.Image
}

func newGame() *game {
	prepareFightingArt()
	b := ebiten.NewImage(20, 20)
	b.Fill(color.RGBA{255, 100, 80, 255})
	g := &game{round: 1, rng: rand.New(rand.NewSource(4207)), audio: audiolab.Context(), pulse: shaderlab.NewPulse(), cam: cameralab.State{ViewW: width, ViewH: height}, badge: b}
	g.resetRound()
	return g
}
func (g *game) resetRound() {
	g.p = fighter{x: 120, hp: 100, guardHP: 100}
	g.ai = fighter{x: 360, hp: 100, guardHP: 100}
	g.frame = 0
	g.hitstop, g.shake = 0, 0
	g.roundOver, g.roundWinner = 0, 0
	g.aiPlan, g.aiPlanTicks = aiIdle, 0
	g.sparks = nil
	g.message = fmt.Sprintf("ROUND %d — FIGHT!", g.round)
}
func (g *game) Update() error {
	if g.matchOver {
		if any() {
			best, streak := g.bestStreak, g.streak
			*g = *newGame()
			g.bestStreak, g.streak = best, streak
		}
		return nil
	}
	if g.hitstop > 0 {
		g.hitstop--
		return nil
	}
	g.frame++
	g.cam.Pos = cameralab.Vec{X: (g.p.x + g.ai.x) / 2, Y: 360}
	if g.shake > 0 {
		g.shake--
	}
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
	if g.roundOver > 0 {
		g.roundOver--
		if g.roundOver == 0 {
			if g.pWins >= 2 || g.aiWins >= 2 {
				g.matchOver = true
			} else {
				g.round++
				g.resetRound()
			}
		}
		return nil
	}
	g.tickFighter(&g.p)
	g.tickFighter(&g.ai)
	move, jab, heavy, guard := readPlayerInput()
	if g.p.canAct() {
		switch {
		case heavy:
			g.startAttack(&g.p, attackHeavy)
		case jab:
			g.startAttack(&g.p, attackJab)
		default:
			g.setGuard(&g.p, guard)
			if !g.p.guard {
				g.moveFighter(&g.p, move, 3.2)
			}
		}
	}
	g.updateAI()
	g.recoverGuard(&g.p)
	g.recoverGuard(&g.ai)
	g.p.x = max(25, min(g.ai.x-42, g.p.x))
	g.ai.x = min(455, max(g.p.x+42, g.ai.x))
	g.updateCombat()
	g.p.x = max(25, min(g.ai.x-42, g.p.x))
	g.ai.x = min(455, max(g.p.x+42, g.ai.x))
	seconds := 45 - g.frame/60
	if g.p.hp <= 0 || g.ai.hp <= 0 || seconds <= 0 {
		if g.p.hp > g.ai.hp {
			g.roundWinner = 1
			g.pWins++
			g.streak++
			if g.streak > g.bestStreak {
				g.bestStreak = g.streak
			}
			g.message = "EBI WINS THE ROUND!"
		} else {
			g.roundWinner = -1
			g.aiWins++
			g.streak = 0
			g.message = "RIVAL WINS THE ROUND!"
		}
		// Keep the finished state on screen so the generated KO pose and the
		// hit/hurt boxes can be inspected before the next round begins.
		g.roundOver = 90
	}
	return nil
}

func readPlayerInput() (move int, jab, heavy, guard bool) {
	left := ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft)
	right := ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight)
	jab = inpututil.IsKeyJustPressed(ebiten.KeyJ) || inpututil.IsKeyJustPressed(ebiten.KeyX)
	heavy = inpututil.IsKeyJustPressed(ebiten.KeyU) || inpututil.IsKeyJustPressed(ebiten.KeyZ)
	guard = ebiten.IsKeyPressed(ebiten.KeyK) || ebiten.IsKeyPressed(ebiten.KeyC)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y >= 610 {
			switch min(4, x/(width/5)) {
			case 0:
				left = true
			case 1:
				right = true
			case 4:
				guard = true
			}
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y >= 610 {
			switch min(4, x/(width/5)) {
			case 2:
				jab = true
			case 3:
				heavy = true
			}
		}
	}
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y >= 610 {
			switch min(4, x/(width/5)) {
			case 0:
				left = true
			case 1:
				right = true
			case 4:
				guard = true
			}
		}
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y >= 610 {
			switch min(4, x/(width/5)) {
			case 2:
				jab = true
			case 3:
				heavy = true
			}
		}
	}
	if left != right {
		if left {
			move = -1
		} else {
			move = 1
		}
	}
	return move, jab, heavy, guard
}

func (f *fighter) canAct() bool {
	return f.stun == 0 && f.blockstun == 0 && f.attackKind == attackNone
}

func (g *game) tickFighter(f *fighter) {
	f.moveDir = 0
	if f.stun > 0 {
		f.stun--
	}
	if f.blockstun > 0 {
		f.blockstun--
	}
	if !f.canAct() {
		f.guard = false
		f.guardAge = 0
	}
}

func (g *game) setGuard(f *fighter, guarding bool) {
	if !guarding || f.guardHP <= 0 {
		f.guard = false
		f.guardAge = 0
		return
	}
	f.guard = true
	f.guardAge++
}

func (g *game) recoverGuard(f *fighter) {
	if f.guard || !f.canAct() || f.guardHP >= 100 || g.frame%3 != 0 {
		return
	}
	f.guardHP++
}

func (g *game) moveFighter(f *fighter, direction int, speed float64) {
	f.moveDir = direction
	if direction == 0 {
		return
	}
	f.x += float64(direction) * speed
	f.walkTick++
}

func (g *game) startAttack(f *fighter, kind int) {
	if !f.canAct() || kind <= attackNone || kind >= len(moveSpecs) {
		return
	}
	f.guard = false
	f.guardAge = 0
	f.attackKind = kind
	f.attackTick = 0
	f.attackConnected = false
}

func attackPhase(f fighter) int {
	if f.attackKind == attackNone {
		return phaseReady
	}
	spec := moveSpecs[f.attackKind]
	switch {
	case f.attackTick < spec.activeStart:
		return phaseStartup
	case f.attackTick < spec.activeEnd:
		return phaseActive
	default:
		return phaseRecovery
	}
}

func (g *game) updateAI() {
	if !g.ai.canAct() {
		g.aiPlan, g.aiPlanTicks = aiIdle, 0
		return
	}
	if g.aiPlanTicks <= 0 {
		g.chooseAIPlan()
	}
	g.aiPlanTicks--
	switch g.aiPlan {
	case aiApproach:
		g.setGuard(&g.ai, false)
		g.moveFighter(&g.ai, -1, 2.25)
	case aiRetreat:
		g.setGuard(&g.ai, false)
		g.moveFighter(&g.ai, 1, 2.8)
	case aiGuard:
		g.setGuard(&g.ai, true)
	default:
		g.setGuard(&g.ai, false)
	}
}

func (g *game) chooseAIPlan() {
	distance := g.ai.x - g.p.x
	playerPhase := attackPhase(g.p)
	if g.p.attackKind != attackNone && playerPhase != phaseRecovery {
		reach := moveSpecs[g.p.attackKind].reach
		if distance < reach+20 {
			if g.ai.guardHP > 22 && g.rng.Intn(100) < 78 {
				g.aiPlan, g.aiPlanTicks = aiGuard, 10+g.rng.Intn(7)
			} else {
				g.aiPlan, g.aiPlanTicks = aiRetreat, 8+g.rng.Intn(5)
			}
			return
		}
	}
	if g.p.attackKind != attackNone && playerPhase == phaseRecovery && distance < moveSpecs[attackHeavy].reach+8 {
		kind := attackJab
		if distance > moveSpecs[attackJab].reach || g.rng.Intn(3) == 0 {
			kind = attackHeavy
		}
		g.startAttack(&g.ai, kind)
		g.aiPlan, g.aiPlanTicks = aiIdle, moveSpecs[kind].total
		g.message = "RIVAL PUNISHES THE RECOVERY!"
		return
	}
	if g.p.guard && distance < moveSpecs[attackHeavy].reach {
		g.startAttack(&g.ai, attackHeavy)
		g.aiPlan, g.aiPlanTicks = aiIdle, moveSpecs[attackHeavy].total
		return
	}
	switch {
	case distance > 145:
		g.aiPlan, g.aiPlanTicks = aiApproach, 8+g.rng.Intn(7)
	case distance < 70:
		g.aiPlan, g.aiPlanTicks = aiRetreat, 7+g.rng.Intn(6)
	default:
		roll := g.rng.Intn(100)
		switch {
		case roll < 28 && g.ai.guardHP > 18:
			g.aiPlan, g.aiPlanTicks = aiGuard, 8+g.rng.Intn(8)
		case roll < 60:
			g.startAttack(&g.ai, attackJab)
			g.aiPlan, g.aiPlanTicks = aiIdle, moveSpecs[attackJab].total
		case roll < 82:
			g.startAttack(&g.ai, attackHeavy)
			g.aiPlan, g.aiPlanTicks = aiIdle, moveSpecs[attackHeavy].total
		default:
			g.aiPlan, g.aiPlanTicks = aiRetreat, 6+g.rng.Intn(6)
		}
	}
}

func attackCanConnect(a, b fighter) bool {
	if a.attackKind == attackNone || a.attackConnected || attackPhase(a) != phaseActive {
		return false
	}
	return math.Abs(b.x-a.x) < moveSpecs[a.attackKind].reach
}

func (g *game) updateCombat() {
	playerConnects := attackCanConnect(g.p, g.ai)
	rivalConnects := attackCanConnect(g.ai, g.p)
	if playerConnects && rivalConnects {
		g.p.attackConnected, g.ai.attackConnected = true, true
		g.cancelAttack(&g.p)
		g.cancelAttack(&g.ai)
		g.p.x -= 8
		g.ai.x += 8
		g.message = "WEAPONS CLASH — BACK TO NEUTRAL!"
		g.hitstop, g.shake = 7, 4
		g.burst((g.p.x+g.ai.x)/2, 510, 16)
	} else {
		if playerConnects {
			g.resolveAttack(&g.p, &g.ai)
		}
		if rivalConnects {
			g.resolveAttack(&g.ai, &g.p)
		}
	}
	g.advanceAttack(&g.p)
	g.advanceAttack(&g.ai)
}

func (g *game) resolveAttack(a, b *fighter) {
	spec := moveSpecs[a.attackKind]
	a.attackConnected = true
	direction := 1.0
	if b.x < a.x {
		direction = -1
	}
	if b.guard && b.canAct() {
		guardDamage := spec.guardDamage
		perfect := b.guardAge > 0 && b.guardAge <= 6
		if perfect {
			guardDamage = max(1, guardDamage/2)
		}
		b.guardHP -= guardDamage
		b.hp -= spec.chip
		b.blockstun = spec.blockstun
		b.x += direction * spec.push * .55
		a.x -= direction * spec.push * .2
		if b.guardHP <= 0 {
			b.guardHP = 45
			b.guard = false
			b.guardAge = 0
			b.blockstun = 0
			b.stun = 34
			g.message = "GUARD BREAK — HEAVY PUNISH WINDOW!"
			g.shake = 7
		} else if perfect {
			g.message = "JUST GUARD — COUNTER NOW!"
		} else {
			g.message = fmt.Sprintf("BLOCKED %s — ATTACKER MUST RECOVER", spec.name)
		}
		g.hitstop = max(g.hitstop, 4)
	} else {
		counter := b.attackKind != attackNone
		damage := spec.damage
		stun := spec.hitstun
		if counter {
			damage += max(3, spec.damage/2)
			stun += 6
			g.message = "COUNTER HIT!"
		} else {
			g.message = fmt.Sprintf("%s CONNECTS!", spec.name)
		}
		b.hp -= damage
		b.stun = stun
		b.guard = false
		b.guardAge = 0
		g.cancelAttack(b)
		b.x += direction * spec.push
		a.x -= direction * spec.push * .15
		g.hitstop = max(g.hitstop, 5+a.attackKind*2)
		if a.attackKind == attackHeavy || counter {
			g.shake = 7
		}
	}
	g.gate.Arm(true)
	g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Noise, 175+float64(a.attackKind)*25, .08)).Play()
	g.burst((a.x+b.x)/2, 510, 9+a.attackKind*5)
}

func (g *game) advanceAttack(f *fighter) {
	if f.attackKind == attackNone {
		return
	}
	f.attackTick++
	if f.attackTick >= moveSpecs[f.attackKind].total {
		g.cancelAttack(f)
	}
}

func (g *game) cancelAttack(f *fighter) {
	f.attackKind = attackNone
	f.attackTick = 0
	f.attackConnected = false
}
func (g *game) burst(x, y float64, n int) {
	for i := 0; i < n; i++ {
		a := float64(i) * math.Pi * 2 / float64(n)
		g.sparks = append(g.sparks, spark{x, y, math.Cos(a) * float64(1+i%3), math.Sin(a) * float64(1+i%3), 24 + float64(i%8)})
	}
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 17, 34, 255})
	drawFightingArena(s)
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.frame)*2) * 5
	}
	vector.DrawFilledRect(s, 0, 0, 480, 132, color.RGBA{5, 11, 26, 205}, false)
	vector.StrokeLine(s, 0, 130, 480, 130, 3, color.RGBA{244, 185, 74, 145}, false)
	draw(s, g.p, "player", true, ox, g.roundWinner < 0 && (g.roundOver > 0 || g.matchOver))
	draw(s, g.ai, "rival", false, ox, g.roundWinner > 0 && (g.roundOver > 0 || g.matchOver))
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), float32(2+p.life/14), color.RGBA{255, 211, 62, 255}, true)
	}
	vector.DrawFilledRect(s, 30, 45, 190, 18, color.RGBA{57, 60, 76, 255}, false)
	vector.DrawFilledRect(s, 260, 45, 190, 18, color.RGBA{57, 60, 76, 255}, false)
	vector.DrawFilledRect(s, 30, 45, float32(190*max(g.p.hp, 0)/100), 18, color.RGBA{45, 225, 194, 255}, false)
	vector.DrawFilledRect(s, 450-float32(190*max(g.ai.hp, 0)/100), 45, float32(190*max(g.ai.hp, 0)/100), 18, color.RGBA{240, 75, 91, 255}, false)
	vector.StrokeRect(s, 28, 43, 194, 22, 2, color.RGBA{255, 255, 255, 120}, false)
	vector.StrokeRect(s, 258, 43, 194, 22, 2, color.RGBA{255, 255, 255, 120}, false)
	vector.DrawFilledRect(s, 30, 68, 190, 6, color.RGBA{22, 42, 69, 220}, false)
	vector.DrawFilledRect(s, 260, 68, 190, 6, color.RGBA{22, 42, 69, 220}, false)
	vector.DrawFilledRect(s, 30, 68, float32(190*max(g.p.guardHP, 0)/100), 6, color.RGBA{95, 190, 255, 255}, false)
	vector.DrawFilledRect(s, 450-float32(190*max(g.ai.guardHP, 0)/100), 68, float32(190*max(g.ai.guardHP, 0)/100), 6, color.RGBA{187, 120, 255, 255}, false)
	sec := max(0, 45-g.frame/60)
	if f, e := uilab.Face("en", 16); e == nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(55, 75)
		text.Draw(s, fmt.Sprintf("EBI %d  HP %03d      %02d      HP %03d  RIVAL %d", g.pWins, max(0, g.p.hp), sec, max(0, g.ai.hp), g.aiWins), f, op)
	} else {
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("EBI %d  HP %03d      %02d      HP %03d  RIVAL %d", g.pWins, max(0, g.p.hp), sec, max(0, g.ai.hp), g.aiWins), 55, 75)
	}
	ebitenutil.DebugPrintAt(s, g.message, 155, 115)
	if g.pulse.Available() {
		fx := ebiten.NewImage(20, 20)
		if g.pulse.Draw(fx, g.badge, float32(g.frame)*.1) {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(440, 12)
			s.DrawImage(fx, op)
		}
	}
	labels := []string{"LEFT", "RIGHT", "JAB", "HEAVY", "GUARD"}
	pressed := []bool{g.p.moveDir < 0, g.p.moveDir > 0, g.p.attackKind == attackJab, g.p.attackKind == attackHeavy, g.p.guard}
	vector.DrawFilledRect(s, 0, 610, width, 110, color.RGBA{4, 11, 25, 178}, false)
	for i, l := range labels {
		x := float32(3 + i*95)
		buttonColor := color.RGBA{70, 120 + uint8(i)*18, 200, 255}
		if pressed[i] {
			buttonColor = color.RGBA{255, 174, 72, 255}
		}
		vector.DrawFilledRect(s, x, 620, 90, 65, buttonColor, false)
		vector.StrokeRect(s, x, 620, 90, 65, 2, color.RGBA{255, 255, 255, 90}, false)
		ebitenutil.DebugPrintAt(s, l, int(x)+15, 650)
	}
	if g.matchOver {
		result := "MATCH LOST"
		if g.pWins > g.aiWins {
			result = "EBI WINS THE MATCH!"
		}
		overlay(s, fmt.Sprintf("%s\nWIN STREAK %d  BEST %d\n\nTAP / SPACE TO REMATCH", result, g.streak, g.bestStreak))
	}
}
func draw(s *ebiten.Image, f fighter, fighterName string, right bool, offset float64, ko bool) {
	x := f.x + offset
	if f.attackKind != attackNone {
		spec := moveSpecs[f.attackKind]
		progress := math.Sin(float64(f.attackTick)/float64(spec.total)*math.Pi) * spec.lunge
		if !right {
			progress = -progress
		}
		x += progress
	}
	if f.stun > 0 && !ko {
		recoil := math.Sin(float64(f.stun)*1.7) * 4
		if right {
			x -= recoil
		} else {
			x += recoil
		}
	}
	vector.DrawFilledCircle(s, float32(x), 588, 42, color.RGBA{3, 8, 18, 100}, true)
	if ko {
		drawFighterPose(s, fighterName+"-ko", x, 605, 270)
	} else if f.stun > 0 {
		drawFighterPose(s, fighterName+"-hurt", x, 605, 270)
	} else if f.attackKind != attackNone {
		frame := attackAnimationFrame(f)
		if f.attackKind == attackHeavy && attackPhase(f) == phaseActive {
			afterX := x - 14
			if !right {
				afterX = x + 14
			}
			drawFighterMotionFrame(s, fighterName, frame, afterX, 605, 270, .22)
		}
		drawFighterMotionFrame(s, fighterName, frame, x, 605, 270, 1)
	} else if f.moveDir != 0 {
		walkFrames := [...]int{0, 1, 2, 1}
		drawFighterMotionFrame(s, fighterName, walkFrames[(f.walkTick/5)%len(walkFrames)], x, 605, 270, 1)
	} else {
		drawFighterPose(s, fighterName+"-ready", x, 605, 270)
	}
	if f.guard {
		vector.StrokeCircle(s, float32(x), 500, 63, 6, color.RGBA{100, 205, 255, 235}, true)
	}
	state := fighterStateLabel(f)
	if state != "READY" {
		ebitenutil.DebugPrintAt(s, state, int(x)-28, 410)
	}
	// Hurt and attack rectangles remain a separate teaching overlay. The
	// generated pose can change without changing these collision rules.
	hurtColor := color.NRGBA{74, 224, 255, 170}
	if !right {
		hurtColor = color.NRGBA{236, 94, 187, 170}
	}
	vector.DrawFilledRect(s, float32(x-27), 432, 54, 158, color.NRGBA{hurtColor.R, hurtColor.G, hurtColor.B, 22}, false)
	vector.StrokeRect(s, float32(x-27), 432, 54, 158, 2, hurtColor, false)
	if f.attackKind != attackNone && attackPhase(f) == phaseActive {
		hitboxWidth := moveSpecs[f.attackKind].reach - 32
		dx := 18.0
		if !right {
			dx = -18 - hitboxWidth
		}
		vector.DrawFilledRect(s, float32(x+dx), 500, float32(hitboxWidth), 42, color.NRGBA{255, 210, 62, 28}, false)
		vector.StrokeRect(s, float32(x+dx), 500, float32(hitboxWidth), 42, 3, color.NRGBA{255, 210, 62, 235}, false)
	}
}

func attackAnimationFrame(f fighter) int {
	if f.attackKind == attackNone {
		return 0
	}
	spec := moveSpecs[f.attackKind]
	switch attackPhase(f) {
	case phaseStartup:
		return min(3, f.attackTick*4/max(1, spec.activeStart))
	case phaseActive:
		return 4 + min(1, (f.attackTick-spec.activeStart)*2/max(1, spec.activeEnd-spec.activeStart))
	default:
		return 6 + min(1, (f.attackTick-spec.activeEnd)*2/max(1, spec.total-spec.activeEnd))
	}
}

func fighterStateLabel(f fighter) string {
	if f.stun > 0 {
		return "HITSTUN"
	}
	if f.blockstun > 0 {
		return "BLOCKSTUN"
	}
	if f.attackKind == attackNone {
		if f.guard {
			return "GUARD"
		}
		return "READY"
	}
	switch attackPhase(f) {
	case phaseStartup:
		return "STARTUP"
	case phaseActive:
		return "ACTIVE"
	default:
		return "RECOVERY"
	}
}
func any() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, msg, 115, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Ebi Fighters — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
