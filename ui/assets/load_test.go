// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package assets_test

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
	"testing"
	"time"

	"github.com/szabba/assert/v2"
	"github.com/szabba/assert/v2/assertions/theerr"
	"github.com/szabba/assert/v2/assertions/theval"

	"github.com/szabba/tob-cob/run"
	"github.com/szabba/tob-cob/ui/assets"
	"github.com/szabba/tob-cob/ui/draw"
	"github.com/szabba/tob-cob/ui/draw/testdraw"
	"github.com/szabba/tob-cob/ui/input"
	"github.com/szabba/tob-cob/ui/input/testinput"
)

//go:embed test-assets
var testAssets embed.FS

var assetFS, _ = fs.Sub(testAssets, "test-assets")

const dt = time.Second / 60

func TestLoadNotStruct(t *testing.T) {
	// given
	src := testinput.Source{}

	game := assets.Load[int](assetFS, nil)

	// when
	next, err := game.Update(src, dt)

	assert.Using(t.Errorf).
		That(theval.Equal(next, nil)).
		That(theerr.Is(err, assets.ErrUnacceptableTarget())).
		That(
			strings.Contains(errMsg(err), "int"),
			"error message %q does not contain %q",
			errMsg(err), "int")
}

func TestLoadEmptyStruct(t *testing.T) {
	// given
	src := testinput.Source{}

	dummy := DummyGame{}
	game := assets.Load(assetFS, func(_ struct{}) run.Game { return dummy })

	// when
	next, err := game.Update(src, dt)

	assert.Using(t.Errorf).
		That(theval.Equal[run.Game](next, dummy)).
		That(theerr.IsNil(err))

}

func TestLoadStructWithNonLoadableFields(t *testing.T) {
	// given
	type NonLoadableAssets struct {
		X bool       // non-loadable type
		y draw.Image // not exported
		Z draw.Image // no struct tag
	}
	src := testinput.Source{}

	dummy := DummyGame{}
	game := assets.Load(assetFS, func(_ NonLoadableAssets) run.Game { return dummy })

	// when
	next, err := game.Update(src, dt)

	assert.Using(t.Errorf).
		That(theval.Equal[run.Game](next, dummy)).
		That(theerr.IsNil(err))
}

func TestLoadStructWithSingleImageFieldMissingAsset(t *testing.T) {
	// given
	type MissingOnlyAsset struct {
		Missing draw.Image `asset:"missing.png"`
	}
	src := testinput.Source{}

	loaded := false
	game := assets.Load(assetFS, func(assets MissingOnlyAsset) run.Game {
		loaded = assets.Missing != nil
		return nil
	})

	dst := &testdraw.Target{}

	// when
	game.Draw(dst, nil)
	_, err := game.Update(src, dt)

	assert.Using(t.Errorf).
		That(!loaded, "tile was loaded").
		That(err != nil, "error is nil")
}

func TestLoadStructWithSingleImageField(t *testing.T) {
	// given
	type OnlyTile struct {
		Tile draw.Image `asset:"tile.png"`
	}
	src := testinput.Source{}

	loaded := false
	game := assets.Load(assetFS, func(assets OnlyTile) run.Game {
		loaded = assets.Tile != nil
		return nil
	})

	dst := &testdraw.Target{}

	// when
	game.Draw(dst, nil)
	_, err := game.Update(src, dt)

	assert.Using(t.Errorf).
		That(loaded, "tile was not loaded").
		That(theerr.IsNil(err))
}

type DummyGame struct{}

func (DummyGame) Draw(_ draw.Target, _ input.Source) {}
func (DummyGame) Update(_ input.Source, _ time.Duration) (run.Game, error) {
	return DummyGame{}, nil
}

func errMsg(err error) string {
	return fmt.Sprint(err)
}
