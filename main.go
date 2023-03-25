// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	_ "image/png"
	"math"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/szabba/tob-cob/game"
	"github.com/szabba/tob-cob/game/actions"
	"github.com/szabba/tob-cob/ui"
	"github.com/szabba/tob-cob/ui/draw"
	"github.com/szabba/tob-cob/ui/draw/pixelgldraw"
	"github.com/szabba/tob-cob/ui/geometry"
	"github.com/szabba/tob-cob/ui/input"
	"github.com/szabba/tob-cob/ui/input/pixelglinput"
)

func main() {

	log.Logger = log.Logger.Level(zerolog.InfoLevel)
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		log.Warn().Msg("build info unavailable")
	}
	log.Info().Interface("build-info", buildInfo).Msg("starting application")

	execDir, err := execDir()
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	g := newGame(execDir)

	pixelgl.Run(g.run)
}

var (
	Black = pixel.RGB(0, 0, 0)
	White = pixel.RGB(1, 1, 1)
	Red   = pixel.RGB(1, 0, 0)
	Gray  = pixel.RGB(.5, .5, .5)
)

type _Game struct {
	execDir string

	space      *game.Space
	placements []game.HeadedPlacement
	actions    []actions.Action

	grid    ui.Grid
	outline ui.GridOutline
	cam     ui.Camera
	camCont *ui.CameraController

	cursorSprite   *ui.Sprite
	humanoidSprite *ui.Sprite
}

func newGame(execDir string) *_Game {
	g := new(_Game)
	g.execDir = execDir

	g.space = game.NewSpace()
	for x := -10; x <= 10; x++ {
		for y := -10; y <= 10; y++ {
			g.space.At(game.P(y, x)).Create()
		}
	}

	g.placements = make([]game.HeadedPlacement, 2)
	g.placements[0].Place(g.space.At(game.P(1, 1)))
	g.placements[1].Place(g.space.At(game.P(0, 0)))
	g.actions = []actions.Action{actions.NoAction(), actions.NoAction()}

	g.grid = ui.Grid{
		CellWidth:  30,
		CellHeight: 30,
	}
	g.outline = ui.GridOutline{
		Space:   g.space,
		Grid:    g.grid,
		Color:   Gray,
		Margins: ui.Margins{X: 2.5, Y: 2.5},
	}

	g.cam = ui.NewCamera(geometry.V(0, 0))
	g.camCont = ui.NewCamController(&g.cam)

	return g
}

func (g *_Game) run() {

	wcfg := pixelgl.WindowConfig{
		Title:   "Tears of Butterflies: Colors of Blood",
		Bounds:  pixel.R(0, 0, 800, 600),
		VSync:   true,
		Monitor: pixelgl.PrimaryMonitor(),
	}

	w, err := pixelgl.NewWindow(wcfg)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}
	defer w.Destroy()
	w.SetCursorVisible(false)
	w.SetSmooth(false)
	w.SetVSync(true)

	inSrc := pixelglinput.New(w)
	dst := pixelgldraw.New(w)

	err = g.load(dst)
	if err != nil {
		log.Error().Err(err).Msg("failed to load assets")
		return
	}

	const dt = time.Second / 60

	for !w.Closed() {
		start := time.Now()

		g.Draw(dst, inSrc)

		w.Update()

		g.Update(inSrc, dt)

		end := time.Now()
		passed := end.Sub(start)
		if dt > passed {
			time.Sleep(dt - passed)
		}
	}
}

func (g *_Game) load(dst draw.Target) error {

	var err error

	g.humanoidSprite, err = ui.LoadSprite(
		filepath.Join(g.execDir, "assets/humanoid.png"),
		dst,
		ui.AnchorSouth())

	if err != nil {
		return err
	}

	g.cursorSprite, err = ui.LoadSprite(
		filepath.Join(g.execDir, "assets/cursor.png"),
		dst,
		ui.AnchorNorthWest())

	return err
}

func (g *_Game) Draw(dst draw.Target, inSrc input.Source) {
	dst.Clear(Black)

	spriteGroup := ui.OrderedSpriteGroup{}

	dst.SetMatrix(g.cam.Matrix(inSrc.Bounds()))
	g.outline.Draw(dst)

	for _, placement := range g.placements {
		matrix := placementTransform(g.outline, placement)
		sprite := g.humanoidSprite.Transform(matrix)
		spriteGroup.Add(sprite)
	}
	spriteGroup.Draw()

	dst.SetMatrix(geometry.Identity())
	g.cursorSprite.
		Transform(geometry.Translation(inSrc.MousePosition())).
		Draw()
}

func (g *_Game) Update(inSrc input.Source, dt time.Duration) {
	if inSrc.JustPressed(input.MouseButtonLeft()) {
		mouseAt := inSrc.MousePosition()
		gridPos := g.grid.UnderCursor(inSrc, g.cam)
		log.Info().
			Float64("screen.mouse.x", mouseAt.X).
			Float64("screen.mouse.y", mouseAt.Y).
			Int("grid.mouse.x", gridPos.Column).
			Int("grid.mouse.y", gridPos.Row).
			Msg("clicked")
	}

	if inSrc.JustPressed(input.MouseButtonLeft()) && !g.placements[0].Headed() {
		src := g.space.At(g.placements[0].AtPoint())
		dst := g.space.At(g.grid.UnderCursor(inSrc, g.cam))
		g.placements[0].Place(src)
		path, _ := game.NewPathFinder(g.space).FindPath(src, dst)
		log.Info().Str("path", fmt.Sprintf("%#v", path)).Msg("found path")
		g.actions[0] = g.placements[0].FollowPath(path, time.Second/8)
	}

	g.camCont.Process(inSrc)

	for i, action := range g.actions {
		g.actions[i] = runFor(action, dt)
	}
}

func runFor(action actions.Action, dt time.Duration) actions.Action {
	if action == nil {
		return actions.NoAction()
	}

	status := action.Run(dt)

	if status.Done() || status.Interrupted() {
		return actions.NoAction()
	}
	return action
}

func placementTransform(outline ui.GridOutline, placement game.HeadedPlacement) geometry.Mat {
	src := placement.AtPoint()
	grid := outline.Grid
	bottom := geometry.V(0, -grid.CellHeight/2+math.Abs(outline.Margins.Y))
	mat := grid.Matrix(src.Column, src.Row).Compose(geometry.Translation((bottom)))
	if placement.Headed() {
		dst := placement.Heading()
		dstMatrix := grid.Matrix(dst.Column, dst.Row).Compose(geometry.Translation(bottom))
		prog := placement.Progress()
		// TODO: factor out mixing function?
		for i := range [...]int{0, 1} {
			for j := range [...]int{0, 1, 2} {
				mat[i][j] = dstMatrix[i][j]*prog + mat[i][j]*(1-prog)
			}
		}
	}
	return mat
}

func execDir() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("cannot get the current executable's directory: %w", err)
	}
	dir := filepath.Dir(path)
	return dir, nil
}
