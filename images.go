package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func LoadTileSprite(path string, height int) TileSprite {
	flat, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("%s/flat.png", path))
	if err != nil {
		log.Fatal(err)
	}
	west, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("%s/west.png", path))
	if err != nil {
		log.Fatal(err)
	}
	south, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("%s/south.png", path))
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
		img, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("%s/%s%d.png", path, name, i))
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
