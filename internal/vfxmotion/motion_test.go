package vfxmotion

import (
	"math"
	"reflect"
	"testing"
)

func TestResolveTapReturnsAStableRuleResult(t *testing.T) {
	tests := []struct {
		name string
		x, y float64
		want TapResult
	}{
		{name: "inside target", x: 103, y: 96, want: TapResult{Outcome: TapHit, X: 100, Y: 100, ScoreDelta: 1}},
		{name: "on boundary", x: 110, y: 100, want: TapResult{Outcome: TapHit, X: 100, Y: 100, ScoreDelta: 1}},
		{name: "outside target", x: 111, y: 100, want: TapResult{Outcome: TapMiss, X: 111, Y: 100}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveTap(tt.x, tt.y, 100, 100, 10)

			if got != tt.want {
				t.Fatalf("ResolveTap() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestGradeRecipeKeepsScoringAndPresentationMappingExplicit(t *testing.T) {
	tests := []struct {
		marker float64
		grade  Grade
		score  int
		label  string
	}{
		{marker: 100, grade: GradePerfect, score: 100, label: "PERFECT"},
		{marker: 125, grade: GradeOK, score: 40, label: "OK"},
		{marker: 160, grade: GradeMiss, score: 0, label: "MISS"},
	}

	for _, tt := range tests {
		grade := JudgeMeter(tt.marker, 100)
		recipe := RecipeForGrade(grade)

		if grade != tt.grade || recipe.Score != tt.score || recipe.Label != tt.label {
			t.Fatalf("marker %.0f: grade=%v recipe=%+v", tt.marker, grade, recipe)
		}
	}
}

func TestProxyFinishesAfterTheGameplayEntityIsGone(t *testing.T) {
	proxy := NewProxy(42, 80, 220, 20, 20, 4)

	startX, startY := proxy.Position()
	for !proxy.Done() {
		proxy.Advance()
	}
	endX, endY := proxy.Position()

	if startX != 80 || startY != 220 {
		t.Fatalf("start = (%.1f, %.1f), want (80, 220)", startX, startY)
	}
	if math.Abs(endX-20) > 0.0001 || math.Abs(endY-20) > 0.0001 {
		t.Fatalf("end = (%.1f, %.1f), want (20, 20)", endX, endY)
	}
}

func TestFlightPoseIsDerivedWithoutChangingPhysics(t *testing.T) {
	tests := []struct {
		name      string
		vy        float64
		sinceFlap int
		crashed   bool
		want      FlightPose
	}{
		{name: "fresh flap", vy: -6, sinceFlap: 2, want: FlightFlap},
		{name: "gliding", vy: 0.5, sinceFlap: 20, want: FlightGlide},
		{name: "diving", vy: 4, sinceFlap: 20, want: FlightDive},
		{name: "crashed", vy: 0, sinceFlap: 20, crashed: true, want: FlightCrash},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := PoseForFlight(tt.vy, tt.sinceFlap, tt.crashed)

			if got != tt.want {
				t.Fatalf("PoseForFlight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrailEvictsOldSamplesAtItsPresentationBudget(t *testing.T) {
	trail := NewTrail(3)
	trail.Push(Point{X: 1})
	trail.Push(Point{X: 2})
	trail.Push(Point{X: 3})
	trail.Push(Point{X: 4})

	got := trail.Points()
	want := []Point{{X: 2}, {X: 3}, {X: 4}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Points() = %+v, want %+v", got, want)
	}
}

func TestTombstoneRetainsTheDeletedEntitySnapshot(t *testing.T) {
	snapshot := Snapshot{ID: 7, X: 20, Y: 30, W: 48, H: 18, Variant: 2}
	tombstone := NewTombstone(snapshot, 3)

	for !tombstone.Done() {
		tombstone.Advance()
	}

	if tombstone.Snapshot != snapshot {
		t.Fatalf("snapshot changed: got %+v, want %+v", tombstone.Snapshot, snapshot)
	}
	if tombstone.Progress() != 1 {
		t.Fatalf("progress = %.2f, want 1", tombstone.Progress())
	}
}

func TestSegmentPoseCoversStraightAndCornerShapes(t *testing.T) {
	tests := []struct {
		name    string
		a, b, c Cell
		kind    SegmentKind
	}{
		{name: "horizontal", a: Cell{0, 1}, b: Cell{1, 1}, c: Cell{2, 1}, kind: SegmentStraight},
		{name: "vertical", a: Cell{1, 0}, b: Cell{1, 1}, c: Cell{1, 2}, kind: SegmentStraight},
		{name: "corner", a: Cell{0, 1}, b: Cell{1, 1}, c: Cell{1, 0}, kind: SegmentCorner},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PoseForSegment(tt.a, tt.b, tt.c)

			if got.Kind != tt.kind {
				t.Fatalf("kind = %v, want %v", got.Kind, tt.kind)
			}
		})
	}
}

func TestEventQueueDrainsOnceAndBudgetCapsOptionalWork(t *testing.T) {
	var queue Queue
	queue.Push(Event{Kind: EventShotFired, X: 10, Y: 20})
	queue.Push(Event{Kind: EventEnemyDestroyed, X: 30, Y: 40})
	budget := Budget{Limit: 5}

	events := queue.Drain()
	first := budget.Take(4)
	second := budget.Take(4)

	if len(events) != 2 || queue.Len() != 0 {
		t.Fatalf("drain: events=%d remaining=%d, want 2 and 0", len(events), queue.Len())
	}
	if first != 4 || second != 1 || budget.Used != 5 {
		t.Fatalf("budget grants = (%d, %d), used=%d", first, second, budget.Used)
	}
}

func TestPlanSokobanMoveDistinguishesWalkPushAndBlockedPush(t *testing.T) {
	walls := map[Cell]bool{{X: 4, Y: 1}: true}
	boxes := []Cell{{X: 2, Y: 1}, {X: 3, Y: 1}}
	tests := []struct {
		name      string
		player    Cell
		direction Cell
		allowed   bool
		boxIndex  int
	}{
		{name: "walk", player: Cell{1, 2}, direction: Cell{1, 0}, allowed: true, boxIndex: -1},
		{name: "push", player: Cell{1, 1}, direction: Cell{1, 0}, allowed: false, boxIndex: -1},
		{name: "blocked by wall", player: Cell{3, 1}, direction: Cell{1, 0}, allowed: false, boxIndex: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PlanSokobanMove(tt.player, tt.direction, walls, boxes)

			if got.Allowed != tt.allowed || got.BoxIndex != tt.boxIndex {
				t.Fatalf("PlanSokobanMove() = %+v", got)
			}
		})
	}

	push := PlanSokobanMove(Cell{1, 3}, Cell{1, 0}, nil, []Cell{{2, 3}})
	if !push.Allowed || push.BoxIndex != 0 || push.BoxTo != (Cell{3, 3}) {
		t.Fatalf("clear push = %+v", push)
	}
}

func TestGroundEdgesFireExactlyOncePerTransition(t *testing.T) {
	states := []bool{false, false, true, true, false}
	var landed, tookOff int

	for i := 1; i < len(states); i++ {
		edges := DetectGroundEdges(states[i-1], states[i])
		if edges.Landed {
			landed++
		}
		if edges.TookOff {
			tookOff++
		}
	}

	if landed != 1 || tookOff != 1 {
		t.Fatalf("landed=%d tookOff=%d, want one each", landed, tookOff)
	}
}

func TestReactionTimelineAppliesVisualPhasesAfterDamage(t *testing.T) {
	reaction := NewReaction(2, 2, 2)
	var got []ReactionPhase

	for reaction.Phase() != ReactionDone {
		got = append(got, reaction.Phase())
		reaction.Advance()
	}

	want := []ReactionPhase{
		ReactionHitStop, ReactionHitStop,
		ReactionFlash, ReactionFlash,
		ReactionRecover, ReactionRecover,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("phases = %v, want %v", got, want)
	}
	if offset := reaction.Offset(8); math.Abs(offset) > 0.0001 {
		t.Fatalf("finished offset = %.2f, want 0", offset)
	}
}

func TestBombScriptReplaysDeterministically(t *testing.T) {
	run := func() [][]Cue {
		script := BombScript(1)
		var frames [][]Cue
		for !script.Done() {
			frames = append(frames, script.Current())
			script.Advance()
		}
		return frames
	}

	first := run()
	second := run()

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("same script produced different cues:\n%v\n%v", first, second)
	}
	if len(first[0]) != 2 || first[1][0].Kind != CueShockwave || first[12][0].Kind != CueConfetti {
		t.Fatalf("unexpected bomb cue schedule: %v", first)
	}
}
