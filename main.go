// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"fmt"
	"image/color"
	_ "image/png"
	"math"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"golang.org/x/exp/slog"

	"github.com/szabba/tob-cob/game"
	"github.com/szabba/tob-cob/game/actions"
	"github.com/szabba/tob-cob/run"
	"github.com/szabba/tob-cob/run/ebitenginerun"
	"github.com/szabba/tob-cob/ui"
	"github.com/szabba/tob-cob/ui/draw"
	"github.com/szabba/tob-cob/ui/geometry"
	"github.com/szabba/tob-cob/ui/input"
)

var _Black = color.Gray{Y: 0}

func main() {
	configLogger()

	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		slog.Warn("build info unavailable")
	}
	slog.Info(
		"starting application",
		slog.Any("build-info", buildInfo))

	config := run.DefaultConfig().
		WithTitle("Tears of Butterflies: Colors of Blood")

	err := ebitenginerun.Game(&_Load{}, config)

	if err != nil {
		slog.Error(
			"fatal error",
			slog.String("err", err.Error()))
		os.Exit(1)
	}
}

func configLogger() {
	envLevel := os.Getenv("LOG_LEVEL")

	lvl := slog.LevelInfo
	levelErr := lvl.UnmarshalText([]byte(envLevel))

	opts := slog.HandlerOptions{
		Level: lvl,
	}

	logFmt := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_FMT")))

	var handler slog.Handler
	if logFmt == "text" {
		handler = opts.NewTextHandler(os.Stderr)
	} else {
		handler = opts.NewJSONHandler(os.Stderr)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	if levelErr != nil {
		slog.Warn(
			"improper LOG_LEVEL",
			slog.Group("env", slog.String("LOG_LEVEL", envLevel)),
			slog.String("err", levelErr.Error()),
		)
	}

	slog.Log(
		context.Background(),
		lvl,
		"set log level")
}

type _Load struct {
	err error

	loadingAttempted bool
	pathPrefix       string

	humanoid, cursor, tile *ui.Sprite
}

func (l *_Load) Draw(dst draw.Target, _ input.Source) {
	dst.Clear(_Black)

	if l.loadingAttempted {
		return
	}
	l.loadingAttempted = true

	l.findPathPrefix()

	l.humanoid = l.load("assets/humanoid.png", dst, ui.AnchorSouth())

	l.cursor = l.load("assets/cursor.png", dst, ui.AnchorNorthWest())

	l.tile = l.load("assets/tile.png", dst, ui.AnchorCenter())
}

func (l *_Load) findPathPrefix() {
	pathPrefix, err := execDir()
	if err != nil {
		l.err = fmt.Errorf("cannot find asset path prefix: %w", err)
		return
	}
	l.pathPrefix = pathPrefix
}

func (l *_Load) load(path string, dst draw.Target, anchor ui.Anchor) *ui.Sprite {
	if l.err != nil {
		return nil
	}

	fpath := filepath.Join(l.pathPrefix, path)

	s, err := ui.LoadSprite(fpath, dst, anchor)
	if err != nil {
		l.err = fmt.Errorf("failed to load asser %q: %w", path, err)
		return nil
	}

	return s
}

func (l *_Load) Update(inSrc input.Source, dt time.Duration) (run.Game, error) {
	if !l.loadingAttempted {
		return l, nil
	}

	if l.err != nil {
		return nil, l.err
	}

	game := newGame(l.humanoid, l.cursor, l.tile)
	return game, nil
}

type _Game struct {
	cursor   *ui.Sprite
	humanoid *ui.Sprite

	space      *game.Space
	placements []game.HeadedPlacement
	actions    []actions.Action

	grid    ui.Grid
	outline ui.GridOutline
	cam     ui.Camera
	camCont *ui.CameraController
}

func newGame(humanoid, cursor, tile *ui.Sprite) *_Game {
	g := new(_Game)

	g.humanoid = humanoid
	g.cursor = cursor

	g.space = game.NewSpace()
	for x := -10; x <= 10; x++ {
		for y := -10; y <= 10; y++ {
			g.space.At(game.P(y, x)).Create()
		}
	}

	g.placements = make([]game.HeadedPlacement, 2)
	g.placements[0].Place(g.space.At(game.P(1, 1)))
	g.placements[1].Place(g.space.At(game.P(0, 0)))
	g.actions = []actions.Action{actions.NoAction(), actions.NoAction()}

	g.grid = ui.Grid{
		CellWidth:  30,
		CellHeight: 30,
	}
	g.outline = ui.GridOutline{
		Sprite:  tile,
		Space:   g.space,
		Grid:    g.grid,
		Margins: ui.Margins{X: 2.5, Y: 2.5},
	}

	g.cam = ui.NewCamera(geometry.V(0, 0))
	g.camCont = ui.NewCamController(&g.cam)

	return g
}

func (g *_Game) Draw(dst draw.Target, inSrc input.Source) {
	dst.Clear(_Black)

	spriteGroup := ui.OrderedSpriteGroup{}

	camMatrix := g.cam.Matrix(inSrc.Bounds())
	dst.SetMatrix(camMatrix)
	g.outline.Draw(dst)

	for _, placement := range g.placements {
		matrix := placementTransform(g.outline, placement)

		sprite := g.humanoid.Transform(matrix)
		spriteGroup.Add(sprite)
	}
	spriteGroup.Draw()

	dst.SetMatrix(geometry.Identity())
	cursorM := geometry.Translation(inSrc.MousePosition())
	g.cursor.
		Transform(cursorM).
		Draw()
}

func (g *_Game) Update(inSrc input.Source, dt time.Duration) (run.Game, error) {
	if inSrc.JustPressed(input.MouseButtonLeft()) {
		mouseAt := inSrc.MousePosition()
		gridPos := g.grid.UnderCursor(inSrc, g.cam)

		slog.Debug("clicked",

			slog.Group("screen",
				slog.Float64("mouse-x", mouseAt.X),
				slog.Float64("mouse-y", mouseAt.Y)),

			slog.Group("grid",
				slog.Int("mouse-x", gridPos.Column),
				slog.Int("mouse-y", gridPos.Row)),
		)
	}

	if inSrc.JustPressed(input.MouseButtonLeft()) && !g.placements[0].Headed() {
		src := g.space.At(g.placements[0].AtPoint())
		dst := g.space.At(g.grid.UnderCursor(inSrc, g.cam))
		g.placements[0].Place(src)
		path, _ := game.NewPathFinder(g.space).FindPath(src, dst)

		if slog.Default().Enabled(nil, slog.LevelDebug) {

			asStr := fmt.Sprintf("%#v", path)
			slog.Debug(
				"found path",
				slog.String("path", asStr))
		}

		g.actions[0] = g.placements[0].FollowPath(path, time.Second/4)
	}

	g.camCont.Process(inSrc)

	for i, action := range g.actions {
		g.actions[i] = runFor(action, dt)
	}

	return g, nil
}

func runFor(action actions.Action, dt time.Duration) actions.Action {
	if action == nil {
		return actions.NoAction()
	}

	status := action.Run(dt)

	if status.Done() || status.Interrupted() {
		return actions.NoAction()
	}
	return action
}

func placementTransform(outline ui.GridOutline, placement game.HeadedPlacement) geometry.Mat {
	src := placement.AtPoint()
	grid := outline.Grid
	bottom := geometry.V(0, -grid.CellHeight/2+math.Abs(outline.Margins.Y))
	mat := grid.Matrix(src.Column, src.Row).Compose(geometry.Translation((bottom)))
	if placement.Headed() {
		dst := placement.Heading()
		dstMatrix := grid.Matrix(dst.Column, dst.Row).Compose(geometry.Translation(bottom))
		prog := placement.Progress()
		// TODO: factor out mixing function?
		for i := range [...]int{0, 1} {
			for j := range [...]int{0, 1, 2} {
				mat[i][j] = dstMatrix[i][j]*prog + mat[i][j]*(1-prog)
			}
		}
	}
	return mat
}

func execDir() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("cannot get the current executable's directory: %w", err)
	}
	dir := filepath.Dir(path)
	return dir, nil
}
