package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func LoadTileSprite(path string, height int) TileSprite {
	flat, _, err := ebitenutil.NewImageFromFile(filepath.Join(path, "flat.png"))
	if err != nil {
		log.Fatal(err)
	}
	west, _, err := ebitenutil.NewImageFromFile(filepath.Join(path, "west.png"))
	if err != nil {
		log.Fatal(err)
	}
	south, _, err := ebitenutil.NewImageFromFile(filepath.Join(path, "south.png"))
	if err != nil {
		log.Fatal(err)
	}
	return TileSprite{
		flat:   flat,
		west:   west,
		south:  south,
		height: height,
	}
}

// assuming all sprites are PNG
func LoadAnimatedSprite(path string, name string, frames int) Animation {

	sprites := []*ebiten.Image{}

	for i := 0; i < frames; i++ {
		img, _, err := ebitenutil.NewImageFromFile(filepath.Join(path, fmt.Sprintf("%s%d.png", name, i)))
		if err != nil {
			log.Fatal(err)
		}
		sprites = append(sprites, img)
	}

	return Animation{
		frame:   0,
		sprites: sprites,
	}
}
