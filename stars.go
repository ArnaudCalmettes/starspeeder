package main

import (
	"image/color"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mlange-42/ark/ecs"
)

func NewStars() System {
	return NewSystem(
		new(StarPopulator),
		new(StarMover),
		new(StarLighter),
		new(StarResetter),
		new(StarDrawer),
	)
}

// Components

// From holds a star segment's source point.
type From struct {
	X, Y float32
}

// To holds a star segment's destination point.
type To struct {
	X, Y float32
}

// Brightness holds a star's brightness.
type Brightness struct {
	V float32
}

// Batch relation helps creating and deleting batches of stars.
type Batch struct {
	ecs.RelationMarker
}

// The StarPopulator system is responsible for matching the number of stars in the simulation
// with settings.StarsCount.
// It typically creates or deletes stars by batches whose size grows or shrinks in powers of 2.
type StarPopulator struct {
	mapper  *ecs.Map4[From, To, Brightness, Batch]
	filter  *ecs.Filter4[From, To, Brightness, Batch]
	batches []starBatch
}

type starBatch struct {
	ID   ecs.Entity
	Size int
}

func (s *StarPopulator) Initialize(w *ecs.World) {
	s.mapper = s.mapper.New(w)
	s.filter = s.filter.New(w)
}

func (s *StarPopulator) Update(w *ecs.World) {
	settings := ecs.GetResource[Settings](w)

	// Reuse as many batches as possible.
	starsCount := 0
	cutIndex := -1
	for i, batch := range s.batches {
		if starsCount+batch.Size > settings.StarsCount {
			cutIndex = i
			break
		}
		starsCount += batch.Size
	}

	// Remove batches with higher indices if we have too many stars.
	if cutIndex != -1 {
		for _, batch := range s.batches[cutIndex:] {
			s.mapper.RemoveBatch(s.filter.Batch(ecs.Rel[Batch](batch.ID)), nil)
			w.RemoveEntity(batch.ID)
		}
		s.batches = s.batches[:cutIndex]
	}

	// Add a new batch of stars to match the desired count.
	if starsCount < settings.StarsCount {
		id := w.NewEntity()
		size := settings.StarsCount - starsCount
		init := func(_ ecs.Entity, from *From, to *To, br *Brightness, _ *Batch) {
			initStar(from, to, br, settings)
		}
		s.mapper.NewBatchFn(size, init, ecs.Rel[Batch](id))
		s.batches = append(s.batches, starBatch{id, size})
	}
}

// initStar initializes a star's components.
func initStar(from *From, to *To, br *Brightness, set *Settings) {
	to.X = rand.Float32() * set.ScreenWidth * set.Scale
	to.Y = rand.Float32() * set.ScreenHeight * set.Scale
	from.X = to.X
	from.Y = to.Y
	br.V = rand.Float32() * 0xff
}

// The StarMover moves the stars by pulling them away from the mouse cursor.
type StarMover struct {
	filter *ecs.Filter2[From, To]
}

func (s *StarMover) Initialize(w *ecs.World) {
	s.filter = s.filter.New(w)
}

func (s *StarMover) Update(w *ecs.World) {
	set := ecs.GetResource[Settings](w)
	mouseX, mouseY := ebiten.CursorPosition()
	x, y := float32(mouseX)*set.Scale, float32(mouseY)*set.Scale

	query := s.filter.Query()
	for query.Next() {
		from, to := query.Get()
		from.X = to.X
		from.Y = to.Y
		to.X += (to.X - x) / set.Speed
		to.Y += (to.Y - y) / set.Speed
	}
}

// The StarResetter resets stars when they go out of screen.
type StarResetter struct {
	filter *ecs.Filter3[From, To, Brightness]
}

func (s *StarResetter) Initialize(w *ecs.World) {
	s.filter = s.filter.New(w)
}

func (s *StarResetter) Update(w *ecs.World) {
	set := ecs.GetResource[Settings](w)

	query := s.filter.Query()
	for query.Next() {
		from, to, br := query.Get()

		if from.X < 0 || set.ScreenWidth*set.Scale < from.X ||
			from.Y < 0 || set.ScreenHeight*set.Scale < from.Y {
			// reuse the star by resetting it
			initStar(from, to, br, set)
		}
	}
}

// The StarLighter makes the stars brighter and brighter.
type StarLighter struct {
	filter *ecs.Filter1[Brightness]
}

func (s *StarLighter) Initialize(w *ecs.World) {
	s.filter = s.filter.New(w)
}

func (s *StarLighter) Update(w *ecs.World) {
	query := s.filter.Query()
	for query.Next() {
		br := query.Get()
		br.V = min(0xff, br.V+1)
	}
}

// The StarDrawer renders the stars on screen.
type StarDrawer struct {
	filter *ecs.Filter3[From, To, Brightness]
}

func (s *StarDrawer) Initialize(w *ecs.World) {
	s.filter = s.filter.New(w)
}

func (s *StarDrawer) Draw(w *ecs.World, screen *ebiten.Image) {
	scale := ecs.GetResource[Settings](w).Scale
	query := s.filter.Query()
	for query.Next() {
		from, to, br := query.Get()
		c := color.RGBA{
			R: uint8(0xbb * br.V / 0xff),
			G: uint8(0xdd * br.V / 0xff),
			B: uint8(0xff * br.V / 0xff),
			A: 0xff,
		}
		vector.StrokeLine(screen, from.X/scale, from.Y/scale, to.X/scale, to.Y/scale, 1, c, false)
	}
}
