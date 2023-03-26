// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ebitenginerun

import (
	"errors"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/szabba/tob-cob/run"
)

func Game(game run.Game, config run.Config) error {
	ebiten.SetWindowSize(config.Size())
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetFullscreen(!config.Windowed())
	ebiten.SetWindowTitle(config.Title())
	ebiten.SetTPS(_TicksPerSecond)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetRunnableOnUnfocused(true)
	if !config.ShowCursor() {
		ebiten.SetCursorMode(ebiten.CursorModeHidden)
	}

	err := ebiten.RunGame(&_EbitenGame{game: game})
	if !errors.Is(err, errQuit) {
		return err
	}
	return nil
}

const (
	_Dt             = time.Second / _TicksPerSecond
	_TicksPerSecond = 60
)

var (
	errNilGame = errors.New("nil game")
	errQuit    = errors.New("game quit intentionally")
)

type _EbitenGame struct {
	game  run.Game
	dst   _DrawTarget
	inSrc _InputSource
}

var _ ebiten.Game = &_EbitenGame{}

func (e *_EbitenGame) Update() error {
	if e == nil || e.game == nil {
		return errNilGame
	}

	next, err := e.game.Update(e.inSrc, _Dt)
	if err != nil {
		return err
	}
	if next == nil {
		return errQuit
	}

	e.game = next
	return nil
}

func (e *_EbitenGame) Draw(screen *ebiten.Image) {
	if e == nil || e.game == nil {
		return
	}

	e.dst.dst = screen
	e.game.Draw(&e.dst, e.inSrc)
}

func (e *_EbitenGame) Layout(
	outsideWidth, outsideHeight int,

) (screenWidth, screenHeight int) {

	e.dst._SetBounds(outsideWidth, outsideHeight)
	e.inSrc._SetBounds(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}
