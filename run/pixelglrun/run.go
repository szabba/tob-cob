// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pixelglrun

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/exp/slog"

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

	start := time.Now()
	for !w.Closed() {
		game.Draw(dst, inSrc)

		if slog.Default().Enabled(nil, slog.LevelDebug) {
			typName := fmt.Sprintf("%T", game)
			slog.Debug("running game", slog.String("game-type", typName))
		}

		w.Update()

		game, err = game.Update(inSrc, dt)
		if err != nil {
			return err
		}
		if game == nil {
			return nil
		}

		end := time.Now()
		sleepRestOfFrame(start, end)
		start = end
	}
	return nil
}

func sleepRestOfFrame(start, end time.Time) {
	passed := end.Sub(start)
	if dt <= passed {
		return
	}

	left := dt - passed

	slog.Debug(
		"sleeping until next frame",
		slog.Duration("time-left", left))

	time.Sleep(left)
}

const dt = time.Second / 60
