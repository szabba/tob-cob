// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui_test

import (
	"testing"

	"github.com/szabba/assert"

	"github.com/szabba/tob-cob/ui"
)

func TestOrderedSpriteGroupDoesNotPanicWhenDrawingNothing(t *testing.T) {
	// given
	group := ui.OrderedSpriteGroup{}

	// when
	oops := catchPanic(func() {
		group.Draw(nil)
	})

	// then
	assert.That(oops == nil, t.Errorf, "unexpected panic: %v", oops)
}

func catchPanic(f func()) (p interface{}) {
	defer func() { p = recover() }()
	f()
	return p
}
