package main

import (
	"fmt"
	"image"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mlange-42/ark/ecs"
)

func NewUI() System {
	return NewSystem(new(UI))
}

type UI struct{}

func (s *UI) Initialize(w *ecs.World) {
	ecs.AddResource(w, new(debugui.DebugUI))
}

func (s *UI) Update(w *ecs.World) {
	settings := ecs.GetResource[Settings](w)
	ecs.GetResource[debugui.DebugUI](w).Update(func(ctx *debugui.Context) error {
		ctx.Window("Settings", image.Rect(10, 150, 125, 225), func(layout debugui.ContainerLayout) {
			ctx.SetGridLayout([]int{-2, -1, -1}, nil)

			ctx.Text("Speed")
			ctx.Button("-").On(func() {
				settings.Speed *= 2
			})
			ctx.Button("+").On(func() {
				settings.Speed = max(1, settings.Speed/2)
			})

			ctx.Text("Stars")
			ctx.Button("-").On(func() {
				settings.StarsCount = max(1, settings.StarsCount/2)
			})
			ctx.Button("+").On(func() {
				settings.StarsCount = min(1024*16, settings.StarsCount*2)
			})
		})
		return nil
	})
}

func (s *UI) Draw(w *ecs.World, screen *ebiten.Image) {
	settings := ecs.GetResource[Settings](w)
	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("Speed: C/%.0f\nStars: %d\nFPS: %.1f\nTPS: %.1f",
			settings.Speed, settings.StarsCount,
			ebiten.ActualFPS(), ebiten.ActualTPS(),
		))
	ecs.GetResource[debugui.DebugUI](w).Draw(screen)
}
