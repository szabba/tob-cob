// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui_test

import (
	"testing"

	"github.com/szabba/assert/v2"
	"github.com/szabba/assert/v2/assertions/theval"

	"github.com/szabba/tob-cob/ui"
	"github.com/szabba/tob-cob/ui/geometry"
)

func TestAnchor(t *testing.T) {
	// given
	kases := map[string]struct {
		Anchor ui.Anchor
		Bounds geometry.Rect
		Point  geometry.Vec
	}{
		"NorthWest": {
			Anchor: ui.AnchorNorthWest(),
			Bounds: geometry.R(-3, -5, 10, 20),
			Point:  geometry.V(-3, -5+20),
		},
		"South": {
			Anchor: ui.AnchorSouth(),
			Bounds: geometry.R(-3, -5, 10, 20),
			Point:  geometry.V(-3+10/2, -5),
		},
		"Center": {
			Anchor: ui.AnchorCenter(),
			Bounds: geometry.R(-3, -5, 10, 20),
			Point:  geometry.V(-3+10/2, -5+20/2),
		},
	}

	for name, tt := range kases {
		t.Run(name, func(t *testing.T) {

			// when
			pt := tt.Anchor.For(tt.Bounds)

			// then

			assert.Using(t.Errorf).
				That(theval.Equal(pt, tt.Point))

		})
	}
}
