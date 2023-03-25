// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package run

import (
	"os"
	"path/filepath"
	"time"

	"github.com/szabba/tob-cob/ui/draw"
	"github.com/szabba/tob-cob/ui/input"
)

type Game interface {
	Draw(tgt draw.Target, src input.Source)
	Update(src input.Source, dt time.Duration) (Game, error)
}

type Loader interface {
	Load(tgt draw.Target) (Game, error)
}

type Config struct {
	title string

	width, height int
}

func DefaultConfig() Config { return Config{} }

func (c Config) WithTitle(title string) Config {
	c.title = title
	return c
}

func (c Config) Title() string {
	if c.title == "" {

		path, err := os.Executable()
		if err != nil {
			return ""
		}

		return filepath.Base(path)
	}

	return c.title
}

func (c Config) WithSize(width, height int) Config {
	c.width = width
	c.height = height
	return c
}

func (c Config) Size() (width, height int) {
	width, height = c.width, c.height

	if width == 0 {
		width = 800
	}

	if height == 0 {
		height = 600
	}

	return width, height
}
