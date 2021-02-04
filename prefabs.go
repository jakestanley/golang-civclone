package main

import "math/rand"

const (
	// WaterHeight default height of water tiles
	WaterHeight = 0
	// GrassHeight default height of grass tiles
	GrassHeight = 4
)

func CreateWater() Square {

	square := CreateSquare()

	square.kind = TWater
	square.height = WaterHeight
	square.liquid = true

	return square
}

func CreateGrass() Square {

	square := CreateSquare()

	square.kind = TGrass
	square.height = GrassHeight
	square.liquid = false

	return square
}

func CreateWoods() Square {
	square := CreateSquare()

	square.kind = TGrass
	square.height = GrassHeight
	square.liquid = false
	square.resource = resourcesTypes[rtForest]

	return square
}

// GrassWorldTiles is an 8x8 grid of grass tiles
func GrassWorldTiles() [][]Square {
	return [][]Square{
		{CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass()},
		{CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass()},
		{CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass()},
		{CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass()},
		{CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass()},
		{CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass()},
		{CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass()},
		{CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass()},
	}
}

// IslandWorldTiles is Kailynn's island
func IslandWorldTiles() [][]Square {
	tiles := [][]Square{
		{CreateWater(), CreateWater(), CreateWater(), CreateWater(), CreateWater(), CreateWater(), CreateWater(), CreateWater()},
		{CreateWater(), CreateWater(), CreateWater(), CreateGrass(), CreateGrass(), CreateWater(), CreateWater(), CreateWater()},
		{CreateWater(), CreateWater(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateWater(), CreateWater()},
		{CreateWater(), CreateWater(), CreateGrass(), CreateWoods(), CreateGrass(), CreateWater(), CreateGrass(), CreateWater()},
		{CreateWater(), CreateWater(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateWater()},
		{CreateWater(), CreateWater(), CreateGrass(), CreateGrass(), CreateWater(), CreateGrass(), CreateWater(), CreateWater()},
		{CreateWater(), CreateWater(), CreateWater(), CreateGrass(), CreateGrass(), CreateWater(), CreateWater(), CreateWater()},
		{CreateWater(), CreateWater(), CreateWater(), CreateWater(), CreateWater(), CreateWater(), CreateWater(), CreateWater()},
	}

	// randomise tile heights
	for x := 0; x < len(tiles); x++ {
		for y := 0; y < len(tiles[x]); y++ {
			tile := &tiles[x][y]
			if tile.kind == TGrass {
				tile.height += ((rand.Intn(4) * 2) - 2)
			}
		}
	}

	return tiles
}

// HeightTestTiles crazy height difference to stress you out
func HeightTestTiles() [][]Square {
	tiles := [][]Square{
		{CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass()},
		{CreateWater(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass()},
		{CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass()},
		{CreateWater(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass()},
		{CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass()},
		{CreateWater(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass()},
		{CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass()},
		{CreateWater(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass(), CreateGrass()},
	}

	return tiles
}
