// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

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

	humanoidSprite, err := ui.LoadSprite("assets/humanoid.png")
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	w, err := pixelgl.NewWindow(wcfg)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}
	defer w.Destroy()

	for !w.Closed() {
		w.Update()
		if w.JustReleased(pixelgl.KeyF) {
			if w.Monitor() == nil {
				w.SetMonitor(pixelgl.PrimaryMonitor())
			} else {
				w.SetMonitor(nil)
			}
		}
		w.Clear(Black)
		camMatrix := cam.Matrix(w.Bounds())
		w.Canvas().SetMatrix(camMatrix)

		grid.Cell(0, 0).Draw(w)
		grid.Cell(0, 1).Draw(w)
		grid.Cell(1, 0).Draw(w)
		grid.Cell(1, 1).Draw(w)
		humanoidSprite.Move(grid, 1, 1).Draw(w)
		humanoidSprite.Move(grid, 0, 1).Draw(w)
		humanoidSprite.Move(grid, 1, 0).Draw(w)
		humanoidSprite.Move(grid, 0, 0).Draw(w)
	}
}
