// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package assets

import (
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"io/fs"
	"reflect"
	"time"

	"github.com/szabba/tob-cob/run"
	"github.com/szabba/tob-cob/ui/draw"
	"github.com/szabba/tob-cob/ui/input"
	"golang.org/x/exp/slog"
)

func ErrUnacceptableTarget() error { return errUnacceptableTarget }

var errUnacceptableTarget = errors.New("unacceptable asset-loading target")

func Load[Loaded any](
	filesys fs.FS,
	next func(assets Loaded) run.Game,
) run.Game {
	return &_Load[Loaded]{
		fs:   filesys,
		next: next,
	}
}

type _Load[Loaded any] struct {
	fs   fs.FS
	next func(assets Loaded) run.Game

	started struct {
		already bool
	}
	err error

	tasksLeft []taskFunc
	loaded    Loaded
}

type taskFunc func(dst draw.Target) error

var _ run.Game = new(_Load[struct{}])

func (l *_Load[Loaded]) Draw(dst draw.Target, src input.Source) {
	l.init()

	if l.err != nil || len(l.tasksLeft) == 0 {
		return
	}

	l.err = l.tasksLeft[0](dst)
	l.tasksLeft = l.tasksLeft[1:]
}

func (l *_Load[Loaded]) Update(src input.Source, _ time.Duration) (run.Game, error) {
	l.init()

	if l.err != nil {
		return nil, l.err
	}

	if len(l.tasksLeft) > 0 {
		return l, nil
	}

	next := l.next(l.loaded)
	return next, nil
}

func (l *_Load[Loaded]) init() {
	if l.started.already || l.err != nil {
		return
	}
	l.started.already = true

	tasks, err := l.computeTasks()
	if err != nil {
		l.err = err
		return
	}
	l.tasksLeft = tasks
}

func (l *_Load[Loaded]) computeTasks() ([]taskFunc, error) {
	typ := reflect.TypeOf(l.loaded)
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%T: not a struct: %w", l.loaded, ErrUnacceptableTarget())
	}

	fields := reflect.VisibleFields(typ)
	tasks := make([]taskFunc, 0, len(fields))

	for _, f := range fields {

		if f.PkgPath != "" {
			slog.Warn("target has unexported field",
				slog.Any("target-type", reflect.TypeOf(l.loaded)),
				slog.Any("field", f),
			)
			continue
		}

		if f.Tag.Get(tagKey) == "" {
			slog.Warn("target has field with missing asset struct tag",
				slog.Any("target-type", reflect.TypeOf(l.loaded)),
				slog.Any("field", f),
				slog.String("missing-key", tagKey),
			)
			continue
		}

		if f.Type == imgType {
			tasks = append(tasks, l.loadImage(f))
		}
	}

	return tasks, nil
}

func (l *_Load[Loaded]) loadImage(typField reflect.StructField) taskFunc {
	return func(dst draw.Target) error {

		fname := typField.Tag.Get(tagKey)

		file, err := l.fs.Open(fname)
		if err != nil {
			return fmt.Errorf("cannot open image %q: %w", fname, err)
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return fmt.Errorf("cannot decode image %q: %w", fname, err)
		}

		imported := dst.Import(img)

		valField := reflect.ValueOf(&l.loaded).Elem().FieldByIndex(typField.Index)
		if !valField.CanSet() {
			return fmt.Errorf("cannot set %s", typField.Name)
		}

		valField.Set(reflect.ValueOf(imported))
		return nil
	}
}

var imgType = func() reflect.Type {
	return reflect.TypeOf(new(draw.Image)).Elem()
}()

const tagKey = "asset"
