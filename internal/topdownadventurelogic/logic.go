// Package topdownadventurelogic contains deterministic rules shared by the
// top-down adventure lessons. It deliberately has no Ebitengine dependency.
package topdownadventurelogic

import "math"

// RoomSpec is the small, data-only contract for one dungeon room. The
// presentation can draw it however it likes, while tests can still verify a
// complete adventure route without opening a game window.
type RoomSpec struct {
	Name, Goal string
	Enemies    int
	NeedsKey   bool
	NeedsTools bool
	Boss       bool
}

// ValidDungeonRoute rejects incomplete room data before it becomes a stage.
func ValidDungeonRoute(rooms []RoomSpec) bool {
	if len(rooms) < 4 || !rooms[0].NeedsKey || !rooms[len(rooms)-1].Boss {
		return false
	}
	for _, room := range rooms {
		if room.Name == "" || room.Goal == "" || room.Enemies < 0 {
			return false
		}
	}
	return true
}

// RunGrade turns a completed run into an easy-to-explain replay goal.
func RunGrade(score, hp, frames int) string {
	if hp >= 4 && score >= 1800 && frames <= 1800 {
		return "S"
	}
	if hp >= 3 && score >= 1300 {
		return "A"
	}
	if hp >= 1 {
		return "B"
	}
	return "C"
}

type Vec struct{ X, Y float64 }
type Rect struct{ X, Y, W, H float64 }

func (a Rect) Intersects(b Rect) bool {
	return a.X < b.X+b.W && a.X+a.W > b.X && a.Y < b.Y+b.H && a.Y+a.H > b.Y
}

func Normalize(v Vec) Vec {
	l := math.Hypot(v.X, v.Y)
	if l == 0 {
		return Vec{}
	}
	return Vec{v.X / l, v.Y / l}
}

// AttackBox places one readable rectangular hit area in the facing direction.
func AttackBox(center, facing Vec, reach, width float64) Rect {
	d := Normalize(facing)
	if math.Abs(d.X) >= math.Abs(d.Y) {
		if d.X < 0 {
			return Rect{center.X - reach, center.Y - width/2, reach, width}
		}
		return Rect{center.X, center.Y - width/2, reach, width}
	}
	if d.Y < 0 {
		return Rect{center.X - width/2, center.Y - reach, width, reach}
	}
	return Rect{center.X - width/2, center.Y, width, reach}
}

type Fighter struct {
	Pos, Velocity    Vec
	HP, Invulnerable int
}

// Hurt applies damage and knockback once, then starts an invulnerable window.
func (f *Fighter) Hurt(damage int, source Vec, invulnerableFrames int) bool {
	if f.Invulnerable > 0 || f.HP <= 0 {
		return false
	}
	f.HP -= damage
	d := Normalize(Vec{f.Pos.X - source.X, f.Pos.Y - source.Y})
	f.Velocity = Vec{d.X * 5, d.Y * 5}
	f.Invulnerable = invulnerableFrames
	return true
}

func (f *Fighter) Tick() {
	if f.Invulnerable > 0 {
		f.Invulnerable--
	}
	f.Pos.X += f.Velocity.X
	f.Pos.Y += f.Velocity.Y
	f.Velocity.X *= .82
	f.Velocity.Y *= .82
}

type RoomPhase int

const (
	RoomLocked RoomPhase = iota
	RoomFight
	RoomCleared
)

type Room struct {
	Phase   RoomPhase
	Enemies int
}

func (r *Room) Enter() {
	if r.Phase == RoomLocked {
		r.Phase = RoomFight
	}
}
func (r *Room) EnemyDefeated() {
	if r.Phase != RoomFight || r.Enemies <= 0 {
		return
	}
	r.Enemies--
	if r.Enemies == 0 {
		r.Phase = RoomCleared
	}
}

type BossPhase int

const (
	BossGuard BossPhase = iota
	BossDash
	BossStorm
	BossDefeated
)

func PhaseForHP(hp, maxHP int) BossPhase {
	if hp <= 0 {
		return BossDefeated
	}
	ratio := float64(hp) / float64(maxHP)
	if ratio > .66 {
		return BossGuard
	}
	if ratio > .33 {
		return BossDash
	}
	return BossStorm
}
