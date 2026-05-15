package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	s := &Settings{
		ScreenWidth:  1280,
		ScreenHeight: 720,
		Scale:        16,
		StarsCount:   1024,
		Speed:        32,
	}
	game := NewGame(s, NewStars(), NewUI())
	ebiten.SetWindowSize(int(s.ScreenWidth), int(s.ScreenHeight))
	ebiten.SetWindowTitle("Star speeder")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
