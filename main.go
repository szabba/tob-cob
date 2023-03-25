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
	"github.com/szabba/tob-cob/ui/draw/pixelgldraw"
	"github.com/szabba/tob-cob/ui/geometry"
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
}

func newGame(execDir string) *_Game {
	g := new(_Game)
	g.execDir = execDir
	return g
}

func (g *_Game) run() {

	wcfg := pixelgl.WindowConfig{
		Title:  "Tears of Butterflies: Colors of Blood",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	}

	grid := ui.Grid{
		CellWidth:  30,
		CellHeight: 30,
	}

	cam := ui.NewCamera(geometry.V(0, 0))
	camCont := ui.NewCamController(&cam)

	space := game.NewSpace()
	for x := -10; x <= 10; x++ {
		for y := -10; y <= 10; y++ {
			space.At(game.P(y, x)).Create()
		}
	}

	outline := ui.GridOutline{
		Space:   space,
		Grid:    grid,
		Color:   Gray,
		Margins: ui.Margins{X: 2.5, Y: 2.5},
	}

	placements := make([]game.HeadedPlacement, 2)
	placements[0].Place(space.At(game.P(1, 1)))
	placements[1].Place(space.At(game.P(0, 0)))
	actions := []actions.Action{actions.NoAction(), actions.NoAction()}

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

	humanoidSprite, err := ui.LoadSprite(
		filepath.Join(g.execDir, "assets/humanoid.png"),
		dst,
		ui.AnchorSouth())

	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	cursorSprite, err := ui.LoadSprite(
		filepath.Join(g.execDir, "assets/cursor.png"),
		dst,
		ui.AnchorNorthWest())

	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	const dt = time.Second / 60

	spriteGroup := ui.OrderedSpriteGroup{}

	for !w.Closed() {

		// Draw
		dst.Clear(Black)

		dst.SetMatrix(cam.Matrix(inSrc.Bounds()))
		outline.Draw(dst)

		for _, placement := range placements {
			matrix := placementTransform(outline, placement)
			sprite := humanoidSprite.Transform(matrix)
			spriteGroup.Add(sprite)
		}
		spriteGroup.Draw()

		dst.SetMatrix(geometry.Identity())
		cursorSprite.
			Transform(geometry.Translation(inSrc.MousePosition())).
			Draw()

		// Update
		w.Update()
		if w.JustReleased(pixelgl.KeyF) {
			if w.Monitor() == nil {
				w.SetMonitor(pixelgl.PrimaryMonitor())
			} else {
				w.SetMonitor(nil)
			}
		}

		if w.JustPressed(pixelgl.MouseButtonLeft) {
			mouseAt := w.MousePosition()
			gridPos := grid.UnderCursor(inSrc, cam)
			log.Info().
				Float64("screen.mouse.x", mouseAt.X).
				Float64("screen.mouse.y", mouseAt.Y).
				Int("grid.mouse.x", gridPos.Column).
				Int("grid.mouse.y", gridPos.Row).
				Msg("clicked")
		}

		if w.JustPressed(pixelgl.MouseButtonLeft) && !placements[0].Headed() {
			src := space.At(placements[0].AtPoint())
			dst := space.At(grid.UnderCursor(inSrc, cam))
			placements[0].Place(src)
			path, _ := game.NewPathFinder(space).FindPath(src, dst)
			log.Info().Str("path", fmt.Sprintf("%#v", path)).Msg("found path")
			actions[0] = placements[0].FollowPath(path, time.Second/8)
		}

		camCont.Process(inSrc)

		for i, action := range actions {
			actions[i] = runFor(action, dt)
		}
		time.Sleep(dt)
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
