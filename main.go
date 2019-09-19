// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/rs/zerolog/log"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	wcfg := pixelgl.WindowConfig{
		Title:  "Tears of Butterflies: Colors of Blood",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	}
	w, err := pixelgl.NewWindow(wcfg)
	if err != nil {
		log.Error().
			Err(err).
			Msg("cannot open window")
		return
	}
	for !w.Closed() {
		w.Update()
	}
}
