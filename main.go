// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	_ "image/png"
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
)

func main() {
	pixelgl.Run(run)
}

var (
	Black = pixel.RGB(0, 0, 0)
	White = pixel.RGB(1, 1, 1)
	Red   = pixel.RGB(1, 0, 0)
)

func run() {
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

	wcfg := pixelgl.WindowConfig{
		Title:  "Tears of Butterflies: Colors of Blood",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	}

	grid := ui.Grid{
		CellWidth:  50,
		CellHeight: 30,
		Dx:         5,
		Dy:         5,
	}

	cam := ui.NewCamera(pixel.V(50, 0))
	camCont := ui.NewCamController(&cam)

	humanoidSprite, err := ui.LoadSprite(filepath.Join(execDir, "assets/humanoid.png"), ui.AnchorSouth())
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}
	cursorSprite, err := ui.LoadSprite(filepath.Join(execDir, "assets/cursor.png"), ui.AnchorNorthWest())
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	space := game.NewSpace()
	for x := -10; x <= 10; x++ {
		for y := -10; y <= 10; y++ {
			space.At(game.P(y, x)).Create()
		}
	}

	outline := ui.GridOutline{Space: space, Grid: grid}

	placements := make([]game.HeadedPlacement, 2)
	placements[0].Place(space.At(game.P(1, 1)))
	actions := []actions.Action{actions.NoAction(), actions.NoAction()}

	w, err := pixelgl.NewWindow(wcfg)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}
	defer w.Destroy()
	w.SetCursorVisible(false)
	w.SetSmooth(true)

	const dt = time.Second / 60

	for !w.Closed() {
		w.Update()
		if w.JustReleased(pixelgl.KeyF) {
			if w.Monitor() == nil {
				w.SetMonitor(pixelgl.PrimaryMonitor())
			} else {
				w.SetMonitor(nil)
			}
		}

		if w.JustPressed(pixelgl.MouseButtonLeft) {
			src := space.At(game.P(3, 0))
			dst := space.At(game.P(1, 1))
			placements[0].Place(src)
			path, _ := game.NewPathFinder(space).FindPath(src, dst)
			actions[0] = placements[0].FollowPath(path, time.Second/2)
		}

		camCont.Process(w)

		w.Clear(Black)

		w.SetMatrix(cam.Matrix(w.Bounds()))
		outline.Draw(w)
		for _, placement := range placements {
			humanoidSprite.Transform(placementTransform(grid, placement)).Draw(w)
		}

		w.SetMatrix(pixel.IM)
		cursorSprite.Transform(pixel.IM.Moved(w.MousePosition())).Draw(w.Canvas())

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

func placementTransform(grid ui.Grid, placement game.HeadedPlacement) pixel.Matrix {
	src := placement.AtPoint()
	mat := grid.Matrix(src.Column, src.Row)
	if placement.Headed() {
		dst := placement.Heading()
		dstMatrix := grid.Matrix(dst.Column, dst.Row)
		prog := placement.Progress()
		for i := range mat {
			mat[i] = dstMatrix[i]*prog + mat[i]*(1-prog)
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
