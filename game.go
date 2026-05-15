package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/ark/ecs"
)

type Game struct {
	World *ecs.World
	System
}

type Settings struct {
	ScreenWidth  float32
	ScreenHeight float32
	Scale        float32
	StarsCount   int
	Speed        float32
}

// NewGame creates a new game that ebiten can run.
// This instanciates a new ecs.World, populates the game's settings
// and register systems to update and draw the game.
func NewGame(settings *Settings, systems ...any) *Game {
	g := Game{
		World:  ecs.NewWorld(),
		System: NewSystem(systems...),
	}
	ecs.AddResource(g.World, settings)
	g.System.Initialize(g.World)
	return &g
}

// Update runs one tick of the simulation.
func (g *Game) Update() error {
	g.System.Update(g.World)
	return nil
}

// Draw draw the game.
func (g *Game) Draw(screen *ebiten.Image) {
	g.System.Draw(g.World, screen)
}

// Layout returns the game's layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	set := ecs.GetResource[Settings](g.World)
	return int(set.ScreenWidth), int(set.ScreenHeight)
}
