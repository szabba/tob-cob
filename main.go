// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	_ "image/png"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/szabba/tob-cob/game"
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

	humanoidSprite, err := ui.LoadSprite("assets/humanoid.png", ui.AnchorSouth())
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}
	cursorSprite, err := ui.LoadSprite("assets/cursor.png", ui.AnchorNorthWest())
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	space := game.NewSpace()
	space.At(game.P(0, 0)).Create()
	space.At(game.P(0, 1)).Create()
	space.At(game.P(1, 0)).Create()
	space.At(game.P(1, 1)).Create()
	placement := game.HeadedPlacement{}
	placement.Place(space.At(game.P(1, 1)))
	action := game.NoAction()

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

		if w.JustReleased(pixelgl.KeySpace) {
			placement.Place(space.At(game.P(1, 1)))
			action = game.Sequence(
				placement.MoveTo(space.At(game.P(0, 0)), 3*time.Second),
				placement.MoveTo(space.At(game.P(1, 0)), 2*time.Second),
				placement.MoveTo(space.At(game.P(0, 1)), 3*time.Second),
				placement.MoveTo(space.At(game.P(1, 1)), 2*time.Second),
			)
		}

		camCont.Process(w)

		w.Clear(Black)

		w.Canvas().SetMatrix(cam.Matrix(w.Bounds()))
		grid.Cell(0, 0).Draw(w)
		grid.Cell(0, 1).Draw(w)
		grid.Cell(1, 0).Draw(w)
		grid.Cell(1, 1).Draw(w)

		humanoidSprite.Transform(placementTransform(grid, placement)).Draw(w)

		w.SetMatrix(pixel.IM)
		cursorSprite.Transform(pixel.IM.Moved(w.MousePosition())).Draw(w.Canvas())

		action.Run(dt)
		time.Sleep(dt)
	}
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
