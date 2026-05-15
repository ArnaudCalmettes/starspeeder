package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/ark/ecs"
)

type System interface {
	Initializer
	Updater
	Drawer
}

type Initializer interface {
	Initialize(w *ecs.World)
}

type Updater interface {
	Update(w *ecs.World)
}

type Drawer interface {
	Draw(w *ecs.World, screen *ebiten.Image)
}

type BaseSystem struct {
	initializers []Initializer
	updaters     []Updater
	drawers      []Drawer
}

func NewSystem(subs ...any) System {
	s := new(BaseSystem)
	for _, sub := range subs {
		used := false
		if initializer, ok := sub.(Initializer); ok {
			s.initializers = append(s.initializers, initializer)
			used = true
		}
		if updater, ok := sub.(Updater); ok {
			s.updaters = append(s.updaters, updater)
			used = true
		}
		if drawer, ok := sub.(Drawer); ok {
			s.drawers = append(s.drawers, drawer)
			used = true
		}
		if !used {
			panic(fmt.Errorf("not a subsystem: %v", sub))
		}
	}
	return s
}

func (s *BaseSystem) Initialize(w *ecs.World) {
	for _, initializer := range s.initializers {
		initializer.Initialize(w)
	}
}

func (s *BaseSystem) Update(w *ecs.World) {
	for _, updater := range s.updaters {
		updater.Update(w)
	}
}

func (s *BaseSystem) Draw(w *ecs.World, screen *ebiten.Image) {
	for _, drawer := range s.drawers {
		drawer.Draw(w, screen)
	}
}
