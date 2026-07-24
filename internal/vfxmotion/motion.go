// Package vfxmotion contains the small, deterministic presentation models used
// by the twelve advanced visual-effects lessons. It deliberately has no
// Ebitengine dependency: rules and animation timing can be tested without a
// window, images, particles, or a GPU.
package vfxmotion

import "math"

// Clamp01 keeps normalized animation time inside its useful range.
func Clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// Lerp interpolates between two values.
func Lerp(from, to, t float64) float64 {
	return from + (to-from)*Clamp01(t)
}

// EaseOutCubic starts quickly and settles gently.
func EaseOutCubic(t float64) float64 {
	t = Clamp01(t)
	return 1 - math.Pow(1-t, 3)
}

// EaseInOutCubic eases at both ends and reaches exactly 0 and 1.
func EaseInOutCubic(t float64) float64 {
	t = Clamp01(t)
	if t < 0.5 {
		return 4 * t * t * t
	}
	return 1 - math.Pow(-2*t+2, 3)/2
}

// Tween is a frame-counted animation clock. Advance is the only mutation;
// Progress is safe to call from Draw and never changes the clock.
type Tween struct {
	Frame  int
	Frames int
}

func NewTween(frames int) Tween {
	if frames < 1 {
		frames = 1
	}
	return Tween{Frames: frames}
}

func (t *Tween) Advance() {
	if t.Frame < t.Frames {
		t.Frame++
	}
}

func (t Tween) Progress() float64 {
	return Clamp01(float64(t.Frame) / float64(max(1, t.Frames)))
}

func (t Tween) Done() bool { return t.Frame >= t.Frames }

// A01 — resolve a tap into an immutable fact before any effect is spawned.
type TapOutcome uint8

const (
	TapMiss TapOutcome = iota
	TapHit
)

type TapResult struct {
	Outcome    TapOutcome
	X, Y       float64
	ScoreDelta int
}

func ResolveTap(tapX, tapY, targetX, targetY, radius float64) TapResult {
	result := TapResult{Outcome: TapMiss, X: tapX, Y: tapY}
	if math.Hypot(tapX-targetX, tapY-targetY) <= radius {
		result.Outcome = TapHit
		result.X, result.Y = targetX, targetY
		result.ScoreDelta = 1
	}
	return result
}

// A02 — map a rule-owned grade to a data-only presentation recipe.
type Grade uint8

const (
	GradeMiss Grade = iota
	GradeOK
	GradePerfect
)

func JudgeMeter(marker, center float64) Grade {
	distance := math.Abs(marker - center)
	switch {
	case distance <= 10:
		return GradePerfect
	case distance <= 40:
		return GradeOK
	default:
		return GradeMiss
	}
}

type GradeRecipe struct {
	Label      string
	Score      int
	BurstCount int
	BurstSpeed float64
	Flash      float64
	Freeze     int
}

func RecipeForGrade(grade Grade) GradeRecipe {
	switch grade {
	case GradePerfect:
		return GradeRecipe{Label: "PERFECT", Score: 100, BurstCount: 32, BurstSpeed: 3.8, Flash: 0.65, Freeze: 4}
	case GradeOK:
		return GradeRecipe{Label: "OK", Score: 40, BurstCount: 16, BurstSpeed: 2.4, Flash: 0.2, Freeze: 2}
	default:
		return GradeRecipe{Label: "MISS", BurstCount: 8, BurstSpeed: 1.1}
	}
}

// A03 — a proxy outlives a removed gameplay entity and flies to a HUD target.
type Proxy struct {
	ID                     int
	FromX, FromY, ToX, ToY float64
	Tween                  Tween
}

func NewProxy(id int, fromX, fromY, toX, toY float64, frames int) Proxy {
	return Proxy{ID: id, FromX: fromX, FromY: fromY, ToX: toX, ToY: toY, Tween: NewTween(frames)}
}

func (p Proxy) Position() (float64, float64) {
	t := EaseInOutCubic(p.Tween.Progress())
	x := Lerp(p.FromX, p.ToX, t)
	y := Lerp(p.FromY, p.ToY, t) - math.Sin(t*math.Pi)*24
	return x, y
}

func (p *Proxy) Advance()  { p.Tween.Advance() }
func (p Proxy) Done() bool { return p.Tween.Done() }

// A04 — derive a readable pose from physics; the pose never feeds physics.
type FlightPose uint8

const (
	FlightGlide FlightPose = iota
	FlightFlap
	FlightDive
	FlightCrash
)

type PoseTransform struct {
	Rotation       float64
	ScaleX, ScaleY float64
}

func PoseForFlight(verticalSpeed float64, framesSinceFlap int, crashed bool) (FlightPose, PoseTransform) {
	if crashed {
		return FlightCrash, PoseTransform{Rotation: math.Pi * 0.48, ScaleX: 1.12, ScaleY: 0.82}
	}
	if framesSinceFlap < 7 {
		return FlightFlap, PoseTransform{Rotation: -0.28, ScaleX: 0.88, ScaleY: 1.16}
	}
	if verticalSpeed > 2.2 {
		return FlightDive, PoseTransform{Rotation: min(0.58, verticalSpeed*0.055), ScaleX: 1.08, ScaleY: 0.92}
	}
	return FlightGlide, PoseTransform{Rotation: max(-0.18, verticalSpeed*0.035), ScaleX: 1, ScaleY: 1}
}

// A05 — a bounded history belongs to presentation, not to the ball.
type Point struct{ X, Y float64 }

type Trail struct {
	points []Point
	limit  int
}

func NewTrail(limit int) Trail {
	if limit < 1 {
		limit = 1
	}
	return Trail{limit: limit}
}

func (t *Trail) Push(point Point) {
	if t.limit < 1 {
		t.limit = 1
	}
	if len(t.points) == t.limit {
		copy(t.points, t.points[1:])
		t.points[len(t.points)-1] = point
		return
	}
	t.points = append(t.points, point)
}

func (t Trail) Points() []Point {
	return append([]Point(nil), t.points...)
}

func (t *Trail) Clear() { t.points = t.points[:0] }

// A06 — copy the draw facts before an entity is deleted from gameplay.
type Snapshot struct {
	ID      int
	X, Y    float64
	W, H    float64
	Variant int
}

type Tombstone struct {
	Snapshot Snapshot
	Tween    Tween
}

func NewTombstone(snapshot Snapshot, frames int) Tombstone {
	return Tombstone{Snapshot: snapshot, Tween: NewTween(frames)}
}

func (t *Tombstone) Advance()         { t.Tween.Advance() }
func (t Tombstone) Done() bool        { return t.Tween.Done() }
func (t Tombstone) Progress() float64 { return EaseOutCubic(t.Tween.Progress()) }

// A07 — work out corners from neighboring cells instead of storing a sprite
// frame in the gameplay body.
type Cell struct{ X, Y int }

type SegmentKind uint8

const (
	SegmentStraight SegmentKind = iota
	SegmentCorner
)

type SegmentPose struct {
	Kind     SegmentKind
	Rotation float64
}

func PoseForSegment(previous, current, next Cell) SegmentPose {
	a := Cell{X: previous.X - current.X, Y: previous.Y - current.Y}
	b := Cell{X: next.X - current.X, Y: next.Y - current.Y}
	if a.X == b.X || a.Y == b.Y {
		if a.X != 0 {
			return SegmentPose{Kind: SegmentStraight, Rotation: math.Pi / 2}
		}
		return SegmentPose{Kind: SegmentStraight}
	}
	rotation := 0.0
	switch {
	case (a.X < 0 && b.Y < 0) || (b.X < 0 && a.Y < 0):
		rotation = math.Pi
	case (a.X < 0 && b.Y > 0) || (b.X < 0 && a.Y > 0):
		rotation = math.Pi / 2
	case (a.X > 0 && b.Y < 0) || (b.X > 0 && a.Y < 0):
		rotation = -math.Pi / 2
	}
	return SegmentPose{Kind: SegmentCorner, Rotation: rotation}
}

// A08 — gameplay emits typed facts; a presentation dispatcher drains them
// once and a budget caps optional work.
type EventKind uint8

const (
	EventShotFired EventKind = iota
	EventEnemyDestroyed
	EventPlayerDamaged
)

type Event struct {
	Kind     EventKind
	X, Y     float64
	Strength float64
}

type Queue struct{ events []Event }

func (q *Queue) Push(event Event) { q.events = append(q.events, event) }

func (q *Queue) Drain() []Event {
	events := append([]Event(nil), q.events...)
	q.events = q.events[:0]
	return events
}

func (q Queue) Len() int { return len(q.events) }

type Budget struct {
	Limit int
	Used  int
}

func (b *Budget) Take(wanted int) int {
	if wanted <= 0 {
		return 0
	}
	remaining := max(0, b.Limit-b.Used)
	granted := min(wanted, remaining)
	b.Used += granted
	return granted
}

func (b *Budget) Reset() { b.Used = 0 }

// A09 — plan a Sokoban move without mutating the board. The caller commits
// this plan once, while Draw interpolates from the captured endpoints.
type MovePlan struct {
	Allowed              bool
	PlayerFrom, PlayerTo Cell
	BoxIndex             int
	BoxFrom, BoxTo       Cell
}

func PlanSokobanMove(player, direction Cell, walls map[Cell]bool, boxes []Cell) MovePlan {
	target := Cell{X: player.X + direction.X, Y: player.Y + direction.Y}
	plan := MovePlan{PlayerFrom: player, PlayerTo: target, BoxIndex: -1}
	if walls[target] {
		return plan
	}
	for i, box := range boxes {
		if box != target {
			continue
		}
		beyond := Cell{X: target.X + direction.X, Y: target.Y + direction.Y}
		if walls[beyond] || containsCell(boxes, beyond) {
			return plan
		}
		plan.BoxIndex = i
		plan.BoxFrom, plan.BoxTo = box, beyond
		break
	}
	plan.Allowed = true
	return plan
}

func containsCell(cells []Cell, target Cell) bool {
	for _, cell := range cells {
		if cell == target {
			return true
		}
	}
	return false
}

// A10 — compare last and current physics state once; never infer landing from
// particles or from a pose.
type GroundEdges struct {
	Landed  bool
	TookOff bool
}

func DetectGroundEdges(wasGrounded, grounded bool) GroundEdges {
	return GroundEdges{Landed: !wasGrounded && grounded, TookOff: wasGrounded && !grounded}
}

type Locomotion uint8

const (
	LocomotionIdle Locomotion = iota
	LocomotionRun
	LocomotionRise
	LocomotionFall
	LocomotionLand
)

func PoseForPlatform(horizontalSpeed, verticalSpeed float64, grounded bool, landedForFrames int) Locomotion {
	if landedForFrames > 0 {
		return LocomotionLand
	}
	if !grounded && verticalSpeed < 0 {
		return LocomotionRise
	}
	if !grounded {
		return LocomotionFall
	}
	if math.Abs(horizontalSpeed) > 0.35 {
		return LocomotionRun
	}
	return LocomotionIdle
}

// A11 — damage is applied before this visual-only timeline starts.
type ReactionPhase uint8

const (
	ReactionHitStop ReactionPhase = iota
	ReactionFlash
	ReactionRecover
	ReactionDone
)

type Reaction struct {
	Frame         int
	HitStopFrames int
	FlashFrames   int
	RecoverFrames int
}

func NewReaction(hitStop, flash, recover int) Reaction {
	return Reaction{
		HitStopFrames: max(0, hitStop),
		FlashFrames:   max(0, flash),
		RecoverFrames: max(0, recover),
	}
}

func (r Reaction) TotalFrames() int {
	return r.HitStopFrames + r.FlashFrames + r.RecoverFrames
}

func (r Reaction) Phase() ReactionPhase {
	switch {
	case r.Frame < r.HitStopFrames:
		return ReactionHitStop
	case r.Frame < r.HitStopFrames+r.FlashFrames:
		return ReactionFlash
	case r.Frame < r.TotalFrames():
		return ReactionRecover
	default:
		return ReactionDone
	}
}

func (r *Reaction) Advance() {
	if r.Frame < r.TotalFrames() {
		r.Frame++
	}
}

func (r Reaction) Offset(strength float64) float64 {
	if r.Phase() == ReactionDone {
		return 0
	}
	remaining := 1 - float64(r.Frame)/float64(max(1, r.TotalFrames()))
	sign := 1.0
	if r.Frame%2 == 1 {
		sign = -1
	}
	return sign * strength * remaining
}

// A12 — a composite effect is a deterministic script of cues. Replaying the
// same script for the same frame emits the same commands.
type CueKind uint8

const (
	CueFreeze CueKind = iota
	CueFlash
	CueShockwave
	CueDissolve
	CueConfetti
)

type Cue struct {
	At       int
	Kind     CueKind
	Strength float64
}

type EffectScript struct {
	Frame int
	Cues  []Cue
}

func BombScript(strength float64) EffectScript {
	return EffectScript{Cues: []Cue{
		{At: 0, Kind: CueFreeze, Strength: 3},
		{At: 0, Kind: CueFlash, Strength: 0.85 * strength},
		{At: 1, Kind: CueShockwave, Strength: 2.2 * strength},
		{At: 4, Kind: CueDissolve, Strength: strength},
		{At: 12, Kind: CueConfetti, Strength: strength},
	}}
}

func (s EffectScript) Current() []Cue {
	var current []Cue
	for _, cue := range s.Cues {
		if cue.At == s.Frame {
			current = append(current, cue)
		}
	}
	return current
}

func (s *EffectScript) Advance() { s.Frame++ }

func (s EffectScript) Done() bool {
	if len(s.Cues) == 0 {
		return true
	}
	return s.Frame > s.Cues[len(s.Cues)-1].At
}
