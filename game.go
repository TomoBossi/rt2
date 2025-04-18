package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type fps struct {
	since              time.Time
	count              int
	average            float64
	averageRefreshRate float64
	cap                int
}

type game struct {
	img        *ebiten.Image
	camera     *camera
	world      *world
	fullscreen bool
	fps        fps
}

type gameParams struct {
	camera *camera
	world  *world
	fpsCap int
}

func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		g.fullscreen = !g.fullscreen
		ebiten.SetFullscreen(g.fullscreen)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errors.New("esc")
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		err := clipboard.WriteAll(fmt.Sprintf("lookFrom: %s,\nlookAt: %s,", g.camera.center, g.camera.center.subtract(g.camera.w)))
		if err != nil {
			return errors.New("copy")
		}
	}

	movement := vec3{0, 0, 0}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		movement = movement.add(vec3{0, 0, -0.2})
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		movement = movement.add(vec3{0, 0, 0.1})
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		movement = movement.add(vec3{-0.1, 0, 0})
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		movement = movement.add(vec3{0.1, 0, 0})
	}
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		movement = movement.add(vec3{0, -0.1, 0})
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		movement = movement.add(vec3{0, 0.1, 0})
	}

	g.updateFps()
	g.camera.translate(movement)
	g.camera.render(g.world)
	g.img.WritePixels(g.camera.pixels)

	return nil
}

func (g *game) updateFps() {
	g.fps.count++
	elapsed := time.Since(g.fps.since).Seconds()
	if elapsed >= 1/g.fps.averageRefreshRate {
		g.fps.average = float64(g.fps.count) / elapsed
		g.fps.count = 0
		g.fps.since = time.Now()
	}
}

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
	game := &game{
		img:    ebiten.NewImage(params.camera.imgWidth, params.camera.imgHeight),
		camera: params.camera,
		world:  params.world,
		fps:    fps{since: time.Now(), cap: params.fpsCap, averageRefreshRate: 1},
	}

	ebiten.SetWindowTitle("(RT)²")
	ebiten.SetWindowSize(800, int(800/params.camera.aspectRatio))
	ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(game.fps.cap)
	ebiten.RunGame(game)
}
