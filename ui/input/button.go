// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package input

type Button struct{ id int }

func (btn Button) Zero() bool { return btn.id == 0 }

func btn() Button {
	if lastID+1 < lastID {
		panic("too many buttons defined")
	}

	lastID++
	return Button{lastID}
}

var lastID int

func MouseButtonLeft() Button { return mbLeft }

func KeyF() Button { return keyF }

func KeyLeft() Button  { return keyLeft }
func KeyUp() Button    { return keyUp }
func KeyRight() Button { return keyRight }
func KeyDown() Button  { return keyDown }

var (
	mbLeft = btn()

	keyF = btn()

	keyLeft  = btn()
	keyUp    = btn()
	keyRight = btn()
	keyDown  = btn()
)
