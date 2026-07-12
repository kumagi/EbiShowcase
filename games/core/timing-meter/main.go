package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 480
	screenHeight = 720
	centerX      = screenWidth / 2
	barLeft      = 45
	barRight     = screenWidth - 45
)

type game struct {
	markerX float64 // 動く線の位置
	speed   float64 // 1フレームで動く量（マイナスなら左へ）
	score   int
	round   int
	stopped bool // 止めた直後か
}

// スペース・クリック・タッチのどれかが「いま押された」か
func justPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return true
	}
	if len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		return true
	}
	return false
}

// 中心からの距離で得点を決める
func pointsForDistance(distance float64) (points int, label string) {
	switch {
	case distance <= 8:
		return 100, "PERFECT +100"
	case distance <= 28:
		return 50, "GREAT +50"
	case distance <= 55:
		return 10, "GOOD +10"
	default:
		return 0, "MISS"
	}
}

// --- ここから Update ---
func (g *game) Update() error {
	if justPressed() {
		if g.stopped {
			// 次のラウンドを始める
			g.stopped = false
			g.round++
			g.speed = 3.2 + float64(g.round)*0.18
			// 奇数ラウンドは左向きからスタート
			if g.round%2 == 1 {
				g.speed = -g.speed
			}
		} else {
			// 動いている線を止めて採点する
			g.stopped = true
			distance := math.Abs(g.markerX - centerX)
			points, _ := pointsForDistance(distance)
			g.score += points
		}
	}

	if g.stopped {
		return nil
	}

	// 線を動かす
	g.markerX += g.speed

	// 端に当たったら向きを変える
	if g.markerX < barLeft || g.markerX > barRight {
		g.speed = -g.speed
		g.markerX = math.Max(barLeft, math.Min(barRight, g.markerX))
	}
	return nil
}

// --- ここから Draw ---
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{7, 20, 38, 255})

	// メーターの帯（外側 → 内側ほど高得点）
	vector.DrawFilledRect(screen, 35, 310, screenWidth-70, 80, color.RGBA{16, 44, 69, 255}, false)
	vector.DrawFilledRect(screen, centerX-55, 310, 110, 80, color.RGBA{255, 105, 79, 255}, false)
	vector.DrawFilledRect(screen, centerX-28, 310, 56, 80, color.RGBA{255, 205, 69, 255}, false)
	vector.DrawFilledRect(screen, centerX-8, 310, 16, 80, color.RGBA{45, 226, 194, 255}, false)

	// 動く白い線
	vector.DrawFilledRect(screen, float32(g.markerX-4), 285, 8, 130, color.White, false)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE %04d   ROUND %02d", g.score, g.round+1), 145, 35)

	if g.stopped {
		distance := math.Abs(g.markerX - centerX)
		_, label := pointsForDistance(distance)
		ebitenutil.DebugPrintAt(screen, label+"\n\nTAP FOR NEXT ROUND", 165, 470)
	} else {
		ebitenutil.DebugPrintAt(screen, "TAP / SPACE TO STOP", 165, 470)
	}
}

func (g *game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Stop the Meter — Ebitengine")
	g := &game{markerX: barLeft, speed: 3.2}
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
