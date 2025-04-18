package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type fps struct {
	since   time.Time
	count   int
	average float64
	cap     int
}

type game struct {
	img        *ebiten.Image
	camera     camera
	world      world
	fullscreen bool
	fps        fps
}

type gameParams struct {
	camera camera
	world  world
	fpsCap int
}

func (g *game) Update() error {
	// Calculate actual FPS
	g.fps.count++
	elapsed := time.Since(g.fps.since).Seconds()
	if elapsed >= 1.0 {
		g.fps.average = float64(g.fps.count) / elapsed
		g.fps.count = 0
		g.fps.since = time.Now()
	}

	// Check for F11 key press to toggle fullscreen
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		g.fullscreen = !g.fullscreen
		ebiten.SetFullscreen(g.fullscreen)
	}

	// Close game on Esc key press
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errors.New("esc")
	}

	// Update pixel data every frame
	g.camera.render(g.world)
	g.img.WritePixels(g.camera.pixels)

	return nil
}

// Draw draws the game screen
func (g *game) Draw(screen *ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	bounds := screen.Bounds()
	scaleX := float64(bounds.Dx()) / float64(g.camera.imgWidth)
	scaleY := float64(bounds.Dy()) / float64(g.camera.imgHeight)
	opt.GeoM.Scale(scaleX, scaleY)

	screen.DrawImage(g.img, opt)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.2f", g.fps.average), 10, 10)
}

func (g *game) Layout(width, height int) (int, int) {
	return width, height
}

func gameInit(params gameParams) {
	// Create the game object
	game := &game{
		camera: params.camera,
		world:  params.world,
		fps:    fps{since: time.Now(), cap: params.fpsCap},
	}

	game.img = ebiten.NewImage(game.camera.imgWidth, game.camera.imgHeight)

	ebiten.SetWindowTitle("RT2")
	ebiten.SetWindowSize(800, 480)
	ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(game.fps.cap)
	ebiten.RunGame(game)
}
