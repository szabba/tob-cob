// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	_ "image/png"
	"math"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/faiface/pixel"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/szabba/tob-cob/game"
	"github.com/szabba/tob-cob/game/actions"
	"github.com/szabba/tob-cob/run"
	"github.com/szabba/tob-cob/run/pixelglrun"
	"github.com/szabba/tob-cob/ui"
	"github.com/szabba/tob-cob/ui/draw"
	"github.com/szabba/tob-cob/ui/geometry"
	"github.com/szabba/tob-cob/ui/input"
)

func main() {
	log.Logger = log.Logger.Level(zerolog.InfoLevel)
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		log.Warn().Msg("build info unavailable")
	}
	log.Info().Interface("build-info", buildInfo).Msg("starting application")

	err := pixelglrun.Game(
		&_Load{},
		run.DefaultConfig().
			WithTitle("Tears of Butterflies: Colors of Blood"))

	if err != nil {
		log.Error().Err(err).Msg("fatal error")
		os.Exit(1)
	}
}

var (
	Black = pixel.RGB(0, 0, 0)
	White = pixel.RGB(1, 1, 1)
	Red   = pixel.RGB(1, 0, 0)
	Gray  = pixel.RGB(.5, .5, .5)
)

type _Load struct {
	err error

	humanoidSprite, cursorSprite *ui.Sprite
}

func (l *_Load) Draw(dst draw.Target, _ input.Source) {
	dst.Clear(Black)

	execDir, err := execDir()
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	l.humanoidSprite, l.err = ui.LoadSprite(
		filepath.Join(execDir, "assets/humanoid.png"),
		dst,
		ui.AnchorSouth())

	if l.err != nil {
		return
	}

	l.cursorSprite, l.err = ui.LoadSprite(
		filepath.Join(execDir, "assets/cursor.png"),
		dst,
		ui.AnchorNorthWest())
}

func (l *_Load) Update(inSrc input.Source, dt time.Duration) (run.Game, error) {

	if l.err != nil {
		return nil, l.err
	}

	game := newGame(l.humanoidSprite, l.cursorSprite)
	return game, nil
}

type _Game struct {
	space      *game.Space
	placements []game.HeadedPlacement
	actions    []actions.Action

	grid    ui.Grid
	outline ui.GridOutline
	cam     ui.Camera
	camCont *ui.CameraController

	cursorSprite   *ui.Sprite
	humanoidSprite *ui.Sprite
}

func newGame(humanoidSprite, cursorSprite *ui.Sprite) *_Game {
	g := new(_Game)

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
		Space:   g.space,
		Grid:    g.grid,
		Color:   Gray,
		Margins: ui.Margins{X: 2.5, Y: 2.5},
	}

	g.cam = ui.NewCamera(geometry.V(0, 0))
	g.camCont = ui.NewCamController(&g.cam)

	g.humanoidSprite = humanoidSprite
	g.cursorSprite = cursorSprite

	return g
}

func (g *_Game) Draw(dst draw.Target, inSrc input.Source) {
	dst.Clear(Black)

	spriteGroup := ui.OrderedSpriteGroup{}

	dst.SetMatrix(g.cam.Matrix(inSrc.Bounds()))
	g.outline.Draw(dst)

	for _, placement := range g.placements {
		matrix := placementTransform(g.outline, placement)
		sprite := g.humanoidSprite.Transform(matrix)
		spriteGroup.Add(sprite)
	}
	spriteGroup.Draw()

	dst.SetMatrix(geometry.Identity())
	g.cursorSprite.
		Transform(geometry.Translation(inSrc.MousePosition())).
		Draw()
}

func (g *_Game) Update(inSrc input.Source, dt time.Duration) (run.Game, error) {
	if inSrc.JustPressed(input.MouseButtonLeft()) {
		mouseAt := inSrc.MousePosition()
		gridPos := g.grid.UnderCursor(inSrc, g.cam)
		log.Info().
			Float64("screen.mouse.x", mouseAt.X).
			Float64("screen.mouse.y", mouseAt.Y).
			Int("grid.mouse.x", gridPos.Column).
			Int("grid.mouse.y", gridPos.Row).
			Msg("clicked")
	}

	if inSrc.JustPressed(input.MouseButtonLeft()) && !g.placements[0].Headed() {
		src := g.space.At(g.placements[0].AtPoint())
		dst := g.space.At(g.grid.UnderCursor(inSrc, g.cam))
		g.placements[0].Place(src)
		path, _ := game.NewPathFinder(g.space).FindPath(src, dst)
		log.Info().Str("path", fmt.Sprintf("%#v", path)).Msg("found path")
		g.actions[0] = g.placements[0].FollowPath(path, time.Second/8)
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
