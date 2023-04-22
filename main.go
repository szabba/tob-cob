// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"fmt"
	"image/color"
	_ "image/png"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"golang.org/x/exp/slog"

	"github.com/szabba/tob-cob/game/actions"
	"github.com/szabba/tob-cob/game/grid"

	"github.com/szabba/tob-cob/run"
	"github.com/szabba/tob-cob/run/ebitenginerun"

	"github.com/szabba/tob-cob/ui"
	"github.com/szabba/tob-cob/ui/assets"
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

	err := mainFallible()

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

func mainFallible() error {

	config := run.DefaultConfig().
		WithTitle("Tears of Butterflies: Colors of Blood")

	execDir, err := execDir()
	if err != nil {
		return err
	}

	assetFs, _ := fs.Sub(os.DirFS(execDir), "assets")

	load := assets.Load(assetFs, func(loaded _Assets) run.Game {
		return newGame(loaded)
	})

	return ebitenginerun.Game(load, config)
}

type _Game struct {
	cursor   ui.Sprite
	humanoid ui.Sprite

	space      *grid.Space
	placements []grid.HeadedPlacement
	actions    []actions.Action

	grid    ui.GridDimensions
	outline ui.GridOutline
	cam     ui.Camera
	camCont *ui.CameraController
}

type _Assets struct {
	Cursor   draw.Image `asset:"cursor.png"`
	Humanoid draw.Image `asset:"humanoid.png"`
	Tile     draw.Image `asset:"tile.png"`
}

func newGame(loaded _Assets) *_Game {
	g := new(_Game)

	g.humanoid = ui.NewSprite(loaded.Humanoid, ui.AnchorSouth())
	g.cursor = ui.NewSprite(loaded.Cursor, ui.AnchorNorthWest())

	g.space = grid.NewSpace()
	for x := -10; x <= 10; x++ {
		for y := -10; y <= 10; y++ {
			g.space.At(grid.P(y, x)).Create()
		}
	}

	g.placements = make([]grid.HeadedPlacement, 2)
	g.placements[0].Place(g.space.At(grid.P(1, 1)))
	g.placements[1].Place(g.space.At(grid.P(0, 0)))
	g.actions = []actions.Action{actions.NoAction(), actions.NoAction()}

	g.grid = ui.GridDimensions{
		CellWidth:  30,
		CellHeight: 30,
	}
	g.outline = ui.GridOutline{
		Sprite:  ui.NewSprite(loaded.Tile, ui.AnchorCenter()),
		Space:   g.space,
		Dims:    g.grid,
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
		path, _ := grid.NewPathFinder(g.space).FindPath(src, dst)

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

func placementTransform(outline ui.GridOutline, placement grid.HeadedPlacement) geometry.Mat {
	src := placement.AtPoint()
	grid := outline.Dims
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
