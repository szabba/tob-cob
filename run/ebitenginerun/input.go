// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ebitenginerun

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/szabba/tob-cob/ui/geometry"
	"github.com/szabba/tob-cob/ui/input"
)

type _InputSource struct {
	bounds geometry.Rect
}

var _ input.Source = _InputSource{}

func (i *_InputSource) _SetBounds(width, height int) {
	i.bounds = geometry.R(0, 0, float64(width), float64(height))
}

func (i _InputSource) Bounds() geometry.Rect { return i.bounds }

func (_InputSource) Focused() bool { return ebiten.IsFocused() }

func (i _InputSource) JustReleased(btn input.Button) bool {
	return i.checkButton(
		btn,
		inpututil.IsKeyJustReleased,
		inpututil.IsMouseButtonJustReleased)
}

func (i _InputSource) JustPressed(btn input.Button) bool {
	return i.checkButton(
		btn,
		inpututil.IsKeyJustPressed,
		inpututil.IsMouseButtonJustPressed)
}

func (i _InputSource) Pressed(btn input.Button) bool {
	return i.checkButton(
		btn,
		ebiten.IsKeyPressed,
		ebiten.IsMouseButtonPressed)
}

func (_InputSource) checkButton(
	btn input.Button,
	keyPred func(ebiten.Key) bool,
	mbPred func(ebiten.MouseButton) bool,
) bool {
	if btn.Zero() {
		return false
	}

	if k, ok := keyMap[btn]; ok {
		return keyPred(k)
	}

	if mb, ok := mbMap[btn]; ok {
		return mbPred(mb)
	}

	return false
}

func (i _InputSource) MousePosition() geometry.Vec {
	x, y := ebiten.CursorPosition()
	return geometry.V(
		float64(x),
		i.bounds.H()-float64(y),
	)
}

func (i _InputSource) MouseInsideWindow() bool {
	pos := i.MousePosition()
	return (i.bounds.Min.X <= pos.X && pos.X <= i.bounds.Max.X) &&
		(i.bounds.Min.Y <= pos.Y && pos.Y <= i.bounds.Max.Y)

}

var mbMap = map[input.Button]ebiten.MouseButton{
	input.MouseButtonLeft(): ebiten.MouseButtonLeft,
}

var keyMap = map[input.Button]ebiten.Key{
	input.KeyF(): ebiten.KeyF,

	input.KeyLeft():  ebiten.KeyArrowLeft,
	input.KeyUp():    ebiten.KeyArrowUp,
	input.KeyRight(): ebiten.KeyArrowRight,
	input.KeyDown():  ebiten.KeyArrowDown,
}
