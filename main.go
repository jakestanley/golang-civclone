package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"

	_ "image/png"

	"github.com/hajimehoshi/ebiten"
)

type Game struct{}

type Citizen struct {
	age int
}

type Animated struct {
	frame   int
	sprites []*ebiten.Image
}

func (a *Animated) Animate() {
	a.frame++
	if a.frame == len(a.sprites) {
		a.frame = 0
	}
}

type TileSprite struct {
	flat   *ebiten.Image
	south  *ebiten.Image
	west   *ebiten.Image
	height int
}

type Tile struct {
	category int
	building int
	selected bool
	op       *ebiten.DrawImageOptions
}

type World struct {
	// TODO Tile struct
	tiles   [][]Tile
	things  [][]Thing
	xOffset int
	yOffset int
}

type Thing struct {
	animated *Animated
	nothing  bool
}

type Civilization struct {
	citizens []Citizen
	max      int
}

const (
	TWater     = 0
	TGrass     = 1
	TileWidth  = 64
	TileHeight = 32
)

var (
	lastFrame    int = 0
	ticks        int = 0
	world        World
	civilization Civilization
	grass        TileSprite
	water        TileSprite
	village      Animated
	house        Animated
	btnEndTurn   *ebiten.Image
	// ctx and cty are the coordinate of the tile that the cursor is on
	ctx int = 0
	cty int = 0
	mtx int = -1
	mty int = -1
	// nothing is used to initialise the 2D things array
	nothing Thing = Thing{
		nothing: true,
	}
)

func UpdateDrawLocations() {

	// north will be top left
	xOffset := 0
	yOffset := 3*32 + 16
	mx, my := ebiten.CursorPosition()
	mxf, myf := float64(mx), float64(my)
	mouseFound := false

	for x := 0; x < len(world.tiles); x++ {
		for y := 0; y < len(world.tiles[x]); y++ {
			// use tile width vars or consts
			op := &ebiten.DrawImageOptions{}
			op.ColorM.Scale(1, 1, 1, 1)

			// tx and ty are where the tile will be drawn from
			tx := float64(xOffset) + float64(y*32) + float64(x*32)
			ty := float64(yOffset) - float64(16*y) + float64(x*16)

			op.GeoM.Translate(tx, ty)
			world.tiles[x][y].op = op

			if world.tiles[x][y].category == TGrass {
				op.GeoM.Translate(0, -float64(grass.height))
			}

			if !mouseFound {
				// this matches a box in the centre of the sprite. needs to actually fit the iso
				// if you treat what the player sees as a rectangle, it won't work correctly
				if (tx+16 < mxf) && (mxf < tx+48) && (ty+8 < myf) && (myf < ty+24) {
					world.tiles[x][y].selected = true
					mtx, mty = x, y
					mouseFound = true
				} else {
					world.tiles[x][y].selected = false
				}
			} else {
				world.tiles[x][y].selected = false
			}
		}
	}
}

func UpdateInputs() {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

		if world.tiles[mtx][mty].category == TGrass {

			world.things[mtx][mty] = Thing{
				animated: &village,
				nothing:  false,
			}
		}

		// civilization.citizens = append(civilization.citizens, Citizen{
		// 	age: 17,
		// })
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if world.tiles[mtx][mty].category == TGrass {

			world.things[mtx][mty] = Thing{
				animated: &house,
				nothing:  false,
			}
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		world.things[mtx][mty] = nothing
	}

	// move cursor north
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		if ctx > 0 {
			ctx--
		}
	}
	// move cursor south
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if ctx < 7 {
			ctx++
		}
	}
	// move cursor west
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		if cty > 0 {
			cty--
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if cty < 7 {
			cty++
		}
	}
}

func (g *Game) Update() error {

	// this also finds which tile the mouse is on
	UpdateDrawLocations()

	UpdateInputs()

	// TODO delta
	// if the game is running at half speed, the delta should be 2
	// if the game is running at normal speed, the delta should be 1 etc

	if ticks == 0 {
		village.Animate()
		house.Animate()
	}
	// possible to use a float here for proper delta time?
	ticks++

	// will mean we update twice per second
	if ticks > 30 {
		ticks = 0
	}

	return nil
}

func DrawWorld(screen *ebiten.Image, world *World) {

	// don't redraw if map doesn't change between frames

	// north will be top left
	// TODO figure this out before the draw call, as we'll need it in the update anyway. saves a load of maths too
	// draw water layer
	for x := 0; x < len(world.tiles); x++ {
		for y := 0; y < len(world.tiles[x]); y++ {
			if world.tiles[x][y].category == TWater {

				if x == ctx && y == cty {
					world.tiles[x][y].op.ColorM.Scale(0.6, 1, 0.6, 1)
				}
				screen.DrawImage(water.flat, world.tiles[x][y].op)

				// if we're at a map edge, also draw the edge tiles
				// ideally we want to also handle adjacent tiles being on the lower layer
				if y == 0 {
					screen.DrawImage(water.west, world.tiles[x][y].op)
				}
				if x == len(world.tiles)-1 {
					screen.DrawImage(water.south, world.tiles[x][y].op)
				}

				// reset the color scaling in case we changed it
				world.tiles[x][y].op.ColorM.Scale(1, 1, 1, 1)
			}
		}
	}

	// draw grass layer
	for x := 0; x < len(world.tiles); x++ {
		for y := 0; y < len(world.tiles[x]); y++ {
			if world.tiles[x][y].category == TGrass {
				if world.tiles[x][y].selected {
					world.tiles[x][y].op.ColorM.Scale(1, 0.6, 0.6, 1)
				}

				screen.DrawImage(grass.flat, world.tiles[x][y].op)

				// if the west adjacent tile is lower, draw the west side
				// need to check array bounds
				if world.tiles[x][y-1].category < world.tiles[x][y].category {
					screen.DrawImage(grass.west, world.tiles[x][y].op)
				}

				// if the south adjacent tile is lower, draw the south side
				// need to check array bounds
				if world.tiles[x+1][y].category < world.tiles[x][y].category {
					screen.DrawImage(grass.south, world.tiles[x][y].op)
				}

				// reset the color scaling in case we changed it
				world.tiles[x][y].op.ColorM.Scale(1, 1, 1, 1)
			}
		}
	}

	// draw things
	for x := 0; x < len(world.things); x++ {
		for y := 0; y < len(world.things[x]); y++ {
			if !world.things[x][y].nothing {
				screen.DrawImage(world.things[x][y].animated.sprites[world.things[x][y].animated.frame], world.tiles[x][y].op)
			}
		}
	}
	//screen.DrawImage(village.sprites[village.frame], &ebiten.DrawImageOptions{})

	// debugging only
	// mx, my := ebiten.CursorPosition()
	// ebitenutil.DrawRect(screen, float64(mx-32), float64(my-16), 64, 32, color.Opaque)
}

func DrawUi(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 1)
	op.GeoM.Translate(16, 200)
	screen.DrawImage(btnEndTurn, op)
}

func (g *Game) Draw(screen *ebiten.Image) {
	// render
	screen.Fill(color.Black)

	DrawWorld(screen, &world)
	DrawUi(screen)

	// TODO if debug
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Citizens: %d/%d", len(civilization.citizens), civilization.max), 16, 40)
	// TODO don't calculate mouse pos on the draw call. this is for debugging only
	mx, my := ebiten.CursorPosition()
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse pos: %d,%d", mx, my), 16, 60)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %f", ebiten.CurrentTPS()), 16, 80)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse on tile: %d, %d", mtx, mty), 16, 100)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Cursor on tile: %d, %d", ctx, cty), 16, 120)
}

func (g *Game) Layout(outsideWith, outsideHeight int) (screenWidth, screenHeight int) {
	return 64 * 8, 32 * 8
	//return int(math.Floor(float64(outsideWith / 2))), int(math.Floor(float64(outsideHeight / 2)))
}

func CreateTile(category int) Tile {
	return Tile{
		category: category,
		selected: false,
	}
}

func CreateSelectedTile(category int) Tile {
	return Tile{
		category: category,
		selected: true,
	}
}

func CreateWorld() World {

	// TODO make sure cardinal directions are easy to spot
	tiles := [][]Tile{
		{CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater)},
		{CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater)},
		{CreateTile(TWater), CreateTile(TWater), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TWater), CreateTile(TWater)},
		{CreateTile(TWater), CreateTile(TWater), CreateTile(TGrass), CreateTile(TWater), CreateTile(TGrass), CreateTile(TWater), CreateTile(TGrass), CreateTile(TWater)},
		{CreateTile(TWater), CreateTile(TWater), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TWater)},
		{CreateTile(TWater), CreateTile(TWater), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TWater), CreateTile(TGrass), CreateTile(TWater), CreateTile(TWater)},
		{CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TGrass), CreateTile(TGrass), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater)},
		{CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater), CreateTile(TWater)},
	}

	w := World{
		tiles: tiles,
	}

	w.CreateThings()
	return w
}

// CreateThings basically just initialises an empty 2D array
func (w *World) CreateThings() {

	t := [][]Thing{}

	for x := 0; x < len(w.tiles); x++ {
		txa := []Thing{}
		for y := 0; y < len(w.tiles[x]); y++ {
			txa = append(txa, nothing)
		}
		t = append(t, txa)
	}

	w.things = t
}

func CreateCivilization() Civilization {

	citizens := []Citizen{}
	citizens = append(citizens, Citizen{
		age: 17,
	})

	return Civilization{
		citizens: citizens,
		max:      0,
	}
}

func LoadSprite(path string, height int) TileSprite {
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
func LoadAnimatedSprite(path string, name string, frames int) Animated {

	sprites := []*ebiten.Image{}

	for i := 0; i < frames; i++ {
		img, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("%s/%s%d.png", path, name, i))
		if err != nil {
			log.Fatal(err)
		}
		sprites = append(sprites, img)
	}

	return Animated{
		frame:   0,
		sprites: sprites,
	}
}

func LoadSprites() {

	// see loading github.com/rakyll/statik in NewImageFromFile documentation

	grass = LoadSprite("img/tiles/grass", 2)
	water = LoadSprite("img/tiles/water", 0)
	village = LoadAnimatedSprite("img/sprites/buildings", "village", 2)
	house = LoadAnimatedSprite("img/sprites/buildings", "house", 2)

	var err error
	btnEndTurn, _, err = ebitenutil.NewImageFromFile("img/ui/btn_end_turn.png")
	if err != nil {
		log.Fatal(err)
	}
}

func Init() {
	LoadSprites()
	world = CreateWorld()
	civilization = CreateCivilization()
}

func main() {
	fmt.Println("hello world")

	Init()

	// do this with a function. it's to make the screen size fit the map (assuming 8x8) like minesweeper
	ebiten.SetWindowSize(64*8*2, 32*8*2)
	ebiten.SetWindowTitle("title")
	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
