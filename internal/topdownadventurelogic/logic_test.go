package topdownadventurelogic

import "testing"

func TestAttackBoxFacesRightAndHitsOnlyFront(t *testing.T) {
	b := AttackBox(Vec{50, 50}, Vec{1, 0}, 30, 20)
	if !b.Intersects(Rect{65, 45, 10, 10}) {
		t.Fatal("front target should be hit")
	}
	if b.Intersects(Rect{35, 45, 10, 10}) {
		t.Fatal("target behind hero must not be hit")
	}
}
func TestDiagonalFacingUsesDominantAxis(t *testing.T) {
	b := AttackBox(Vec{40, 40}, Vec{.2, -1}, 24, 16)
	if !b.Intersects(Rect{36, 18, 8, 8}) {
		t.Fatal("up-facing attack missing")
	}
}
func TestHurtStartsInvulnerabilityAndKnockback(t *testing.T) {
	f := Fighter{Pos: Vec{10, 0}, HP: 3}
	if !f.Hurt(1, Vec{0, 0}, 5) {
		t.Fatal("first hit rejected")
	}
	if f.HP != 2 || f.Velocity.X <= 0 {
		t.Fatal("damage or knockback missing")
	}
	if f.Hurt(1, Vec{0, 0}, 5) {
		t.Fatal("hit during invulnerability accepted")
	}
	for i := 0; i < 5; i++ {
		f.Tick()
	}
	if !f.Hurt(1, Vec{0, 0}, 5) {
		t.Fatal("hit after recovery rejected")
	}
}
func TestRoomClearsOnlyAfterLastEnemy(t *testing.T) {
	r := Room{Enemies: 2}
	r.Enter()
	r.EnemyDefeated()
	if r.Phase != RoomFight {
		t.Fatal("room cleared early")
	}
	r.EnemyDefeated()
	if r.Phase != RoomCleared {
		t.Fatal("room did not clear")
	}
}
func TestBossPhases(t *testing.T) {
	cases := []struct {
		hp   int
		want BossPhase
	}{{9, BossGuard}, {5, BossDash}, {2, BossStorm}, {0, BossDefeated}}
	for _, c := range cases {
		if got := PhaseForHP(c.hp, 9); got != c.want {
			t.Fatalf("hp %d: got %v", c.hp, got)
		}
	}
}

func TestDungeonRouteNeedsDataKeyAndBoss(t *testing.T) {
	rooms := []RoomSpec{
		{Name: "Key", Goal: "Open the door", NeedsKey: true},
		{Name: "Fight", Goal: "Clear crawlers", Enemies: 4},
		{Name: "Tools", Goal: "Light seals", NeedsTools: true},
		{Name: "Guardian", Goal: "Defeat boss", Enemies: 1, Boss: true},
	}
	if !ValidDungeonRoute(rooms) {
		t.Fatal("complete four-room route rejected")
	}
	rooms[0].NeedsKey = false
	if ValidDungeonRoute(rooms) {
		t.Fatal("route without a key room accepted")
	}
}

func TestRunGradeRewardsHealthScoreAndTime(t *testing.T) {
	if got := RunGrade(1900, 5, 1200); got != "S" {
		t.Fatalf("S run = %q", got)
	}
	if got := RunGrade(1400, 3, 2200); got != "A" {
		t.Fatalf("A run = %q", got)
	}
	if got := RunGrade(700, 1, 4000); got != "B" {
		t.Fatalf("survival run = %q", got)
	}
	if got := RunGrade(0, 0, 4000); got != "C" {
		t.Fatalf("empty run = %q", got)
	}
}
