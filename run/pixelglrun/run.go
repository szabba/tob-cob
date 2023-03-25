// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pixelglrun

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/rs/zerolog/log"

	"github.com/szabba/tob-cob/run"
	"github.com/szabba/tob-cob/ui/draw/pixelgldraw"
	"github.com/szabba/tob-cob/ui/input/pixelglinput"
)

func Game(game run.Game, config run.Config) error {
	var err error

	pixelgl.Run(func() { err = runFallibly(game, config) })

	return err
}

func runFallibly(game run.Game, config run.Config) error {

	width, height := config.Size()

	wcfg := pixelgl.WindowConfig{
		Title:   config.Title(),
		Bounds:  pixel.R(0, 0, float64(width), float64(height)),
		VSync:   true,
		Monitor: pixelgl.PrimaryMonitor(),
	}

	w, err := pixelgl.NewWindow(wcfg)
	if err != nil {
		return err
	}
	defer w.Destroy()
	w.SetCursorVisible(false)
	w.SetSmooth(false)
	w.SetVSync(true)

	inSrc := pixelglinput.New(w)
	dst := pixelgldraw.New(w)

	const dt = time.Second / 60

	for !w.Closed() {
		log.Debug().Str("game.type", fmt.Sprintf("%T", game)).Msg("running game")
		start := time.Now()

		game.Draw(dst, inSrc)

		w.Update()

		game, err = game.Update(inSrc, dt)
		if err != nil {
			return err
		}
		if game == nil {
			return nil
		}

		end := time.Now()
		passed := end.Sub(start)
		if dt > passed {
			left := dt - passed
			log.Debug().Dur("time.left", left).Msg("sleeping until next frame")
			time.Sleep(left)
		}
	}
	return nil
}
