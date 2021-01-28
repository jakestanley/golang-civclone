package main

const (
	// WaterHeight default height of water tiles
	WaterHeight = 0
	// GrassHeight default height of grass tiles
	GrassHeight = 4
)

// GrassWorldTiles is an 8x8 grid of grass tiles
func GrassWorldTiles() [][]Tile {
	return [][]Tile{
		{CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false)},
		{CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false)},
		{CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false)},
		{CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false)},
		{CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false)},
		{CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false)},
		{CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false)},
		{CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false)},
	}
}

// IslandWorldTiles is Kailynn's island
func IslandWorldTiles() [][]Tile {
	return [][]Tile{
		{CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true)},
		{CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true)},
		{CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true)},
		{CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TGrass, GrassHeight, false), CreateTile(TWater, WaterHeight, true), CreateTile(TGrass, GrassHeight, false), CreateTile(TWater, WaterHeight, true), CreateTile(TGrass, GrassHeight, false), CreateTile(TWater, WaterHeight, true)},
		{CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TWater, WaterHeight, true)},
		{CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TWater, WaterHeight, true), CreateTile(TGrass, GrassHeight, false), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true)},
		{CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true)},
		{CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true), CreateTile(TWater, WaterHeight, true)},
	}
}

// HeightTestTiles crazy height difference to stress you out
func HeightTestTiles() [][]Tile {
	return [][]Tile{
		{CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false)},
		{CreateTile(TWater, WaterHeight, true), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false)},
		{CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false)},
		{CreateTile(TWater, WaterHeight, true), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false)},
		{CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false)},
		{CreateTile(TWater, WaterHeight, true), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false)},
		{CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false)},
		{CreateTile(TWater, WaterHeight, true), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false), CreateTile(TGrass, GrassHeight, false)},
	}
}
