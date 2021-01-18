package main

// GrassWorldTiles is an 8x8 grid of grass tiles
func GrassWorldTiles() [][]Tile {
	return [][]Tile{
		{CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass)},
		{CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass)},
		{CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass)},
		{CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass)},
		{CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass)},
		{CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass)},
		{CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass)},
		{CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass)},
	}
}

// IslandWorldTiles is Kailynn's island
func IslandWorldTiles() [][]Tile {
	return [][]Tile{
		{CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater)},
		{CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater)},
		{CreateTile(TWater), CreateTile(TWater), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TWater), CreateTile(TWater)},
		{CreateTile(TWater), CreateTile(TWater), CreateTile(TGrass), CreateTile(TWater), CreateTile(TGrass), CreateTile(TWater), CreateTile(TGrass), CreateTile(TWater)},
		{CreateTile(TWater), CreateTile(TWater), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TWater)},
		{CreateTile(TWater), CreateTile(TWater), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TWater), CreateTile(TGrass), CreateTile(TWater), CreateTile(TWater)},
		{CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater)},
		{CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater)},
	}
}
