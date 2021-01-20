package main

import (
	"io/ioutil"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func LoadFonts() {

	const dpi = 72

	fopsTitle := &opentype.FaceOptions{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingNone,
	}

	fopsDetail := &opentype.FaceOptions{
		Size:    9,
		DPI:     dpi,
		Hinting: font.HintingNone,
	}

	// load title font
	data, err := ioutil.ReadFile("font/alagard.ttf")
	if err != nil {
		log.Fatal(err)
	}

	ttf, err := opentype.Parse(data)
	if err != nil {
		log.Fatal(err)
	}

	fontTitle, err = opentype.NewFace(ttf, fopsTitle)

	if err != nil {
		log.Fatal(err)
	}

	// load detail font
	data, err = ioutil.ReadFile("font/Volter__28Goldfish_29.ttf")
	if err != nil {
		log.Fatal(err)
	}

	ttf, err = opentype.Parse(data)

	if err != nil {
		log.Fatal(err)
	}

	fontDetail, err = opentype.NewFace(ttf, fopsDetail)

	if err != nil {
		log.Fatal(err)
	}

}
