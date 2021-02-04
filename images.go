package main

import (
	"fmt"
	"image"
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var icons []*ebiten.Image

const IconsUranium = 0
const IconsOil = 1
const IconsFood = 2
const IconsMedicine = 3
const IconsWood = 4
const IconsStone = 5

func LoadTileSprite(path string) TileSprite {
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
	westMid, _, err := ebitenutil.NewImageFromFile(filepath.Join(path, "west-mid.png"))
	if err != nil {
		log.Fatal(err)
	}
	southMid, _, err := ebitenutil.NewImageFromFile(filepath.Join(path, "south-mid.png"))
	if err != nil {
		log.Fatal(err)
	}
	return TileSprite{
		flat:     flat,
		west:     west,
		westMid:  westMid,
		south:    south,
		southMid: southMid,
	}
}

func LoadIcons() {

	iconSize := 16

	img, _, err := ebitenutil.NewImageFromFile(filepath.Join("img", "icons", "resources.png"))
	if err != nil {
		log.Fatal(err)
	}

	icons = []*ebiten.Image{}
	count := img.Bounds().Dx() / iconSize
	for i := 0; i < count; i++ {
		min := image.Point{
			X: i * iconSize,
			Y: 0,
		}
		max := image.Point{
			X: i*iconSize + 15,
			Y: iconSize,
		}
		icons = append(icons, ebiten.NewImageFromImage(img.SubImage(image.Rectangle{Min: min, Max: max})))
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
