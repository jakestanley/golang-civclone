package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const dpi = 72

func LoadTtfFont(path string, size int) font.Face {

	fops := &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     dpi,
		Hinting: font.HintingNone,
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	ttf, err := opentype.Parse(data)
	if err != nil {
		log.Fatal(err)
	}

	font, err := opentype.NewFace(ttf, fops)
	if err != nil {
		log.Fatal(err)
	}

	return font
}

func LoadFonts() {

	fontTitle = LoadTtfFont(filepath.Join("font", "alagard.ttf"), 16)
	fontDetail = LoadTtfFont(filepath.Join("font", "Volter__28Goldfish_29.ttf"), 9)
	fontSmall = LoadTtfFont(filepath.Join("font", "small_pixel.ttf"), 8)

}
