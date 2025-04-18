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

type game struct {
	img        *ebiten.Image
	camera     *camera
	world      *world
	fullscreen bool
	mouse      mouse
	fps        fps
}

type gameParams struct {
	camera     *camera
	world      *world
	fpsCap     int
	fullscreen bool
}

type fps struct {
	since              time.Time
	count              int
	average            float64
	averageRefreshRate float64
	cap                int
}

type mouse struct {
	x, y int
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
		err := clipboard.WriteAll(fmt.Sprintf("verticalFov: %.2f,\nlookFrom: %s,\nlookAt: %s,", g.camera.verticalFov, g.camera.center, g.camera.center.subtract(g.camera.w)))
		if err != nil {
			return errors.New("copy")
		}
	}

	movement := vec3{0, 0, 0}
	movementScale := 1.0
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		movementScale = 0.2
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		movement = movement.add(vec3{0, 0, -0.1})
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		movement = movement.add(vec3{0, 0, 0.05})
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		movement = movement.add(vec3{-0.05, 0, 0})
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		movement = movement.add(vec3{0.05, 0, 0})
	}
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		movement = movement.add(vec3{0, -0.05, 0})
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		movement = movement.add(vec3{0, 0.05, 0})
	}
	movement = movement.scale(movementScale)

	mx, my := ebiten.CursorPosition()
	pitch, yaw := -float64(mx-g.mouse.x)*0.002, -float64(my-g.mouse.y)*0.002

	fov := 0.0
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		fov += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		fov -= 1
	}

	g.updateFps()
	g.camera.update(movement, fov, pitch, yaw)
	g.camera.render(g.world)
	g.img.WritePixels(g.camera.pixels)
	g.mouse.x, g.mouse.y = mx, my
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
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FOV:  %.2f\nFROM: %s\nAT:   %s", g.camera.verticalFov, g.camera.center, g.camera.center.subtract(g.camera.w)), 10, bounds.Dy()-60)
}

func (g *game) Layout(width, height int) (int, int) {
	return width, height
}

func gameInit(params gameParams) {
	windowWidth := 800
	windowHeight := int(float64(windowWidth) / params.camera.aspectRatio)

	game := &game{
		img:        ebiten.NewImage(params.camera.imgWidth, params.camera.imgHeight),
		camera:     params.camera,
		world:      params.world,
		fullscreen: params.fullscreen,
		fps:        fps{since: time.Now(), cap: params.fpsCap, averageRefreshRate: 1},
	}

	ebiten.SetWindowTitle("(RT)²")
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetFullscreen(params.fullscreen)
	ebiten.SetTPS(game.fps.cap)
	ebiten.RunGame(game)
	mx, my := ebiten.CursorPosition()
	game.mouse.x, game.mouse.y = mx, my
}
