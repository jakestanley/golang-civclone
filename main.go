package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"

	_ "image/png"

	"github.com/hajimehoshi/ebiten"
)

type Game struct{}

type Citizen struct {
	age int
}

type MessageQueue struct {
	queue []string
	max   int
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

type UiSprite struct {
	left *ebiten.Image
	// middle MUST be one pixel wide for text fit scaling
	middle *ebiten.Image
	right  *ebiten.Image
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

type Settlement struct {
	citizens []*Citizen
}

type World struct {
	// TODO Tile struct
	tiles   [][]Tile
	things  [][]Thing
	xOffset int
	yOffset int
}

// TODO consider polymorphism for buildings
type Thing struct {
	progress  float64
	completed bool
	animated  *Animated
	nothing   bool
}

type Button struct {
	redraw  bool
	content string
	x, y    int
	hover   bool // should default to false
	img     *UiSprite
	bounds  image.Rectangle
}

type Civilization struct {
	citizens []Citizen
	max      int
}

const (
	// WWidth default window width
	WWidth = 1024
	// WHeight default window height
	WHeight = 600
	// TWater water tile type index
	TWater = 0
	// TGrass grass tile type index
	TGrass = 1
	// TileWidth width of tiles in pixels (unscaled)
	TileWidth = 64
	// TileHeight height of tiles in pixels (unscaled)
	TileHeight = 32
	// SettlementCapacityVillage max citizens a settlement can contain
	SettlementCapacityVillage = 10
	// BtnEndTurn is the button map key for ending a turn
	BtnEndTurn       = "END_TURN"
	BtnShowBuildings = "SHOW_BUILDINGS"
)

var (
	epochs []string = []string{
		"Neolithic Age",
		"Roman Age",
		"Classical Age",
		"Age of Steam",
		"Modern Age",
		"Transhuman Age",
		"Apocalyptic Age",
	}
	sHeight      int
	sWidth       int
	lastFrame    int = 0
	ticks        int = 0
	year         int = 1
	epoch        int = 0
	world        World
	civilization Civilization
	fontTitle    font.Face
	fontDetail   font.Face
	grass        TileSprite
	water        TileSprite
	village      Animated
	house        Animated
	north        *ebiten.Image
	btn          UiSprite
	// ctx and cty are the coordinate of the tile that the cursor is on
	ctx                 int  = 0
	cty                 int  = 0
	mx                  int  = 0
	my                  int  = 0
	mtx                 int  = -1
	mty                 int  = -1
	validMouseSelection bool = false
	// nothing is used to initialise the 2D things array
	nothing Thing = Thing{
		nothing: true,
	}

	Messages MessageQueue

	// AllButtons is the master list of buttons. Used by renderer and mouse picker
	AllButtons []*Button

	// SButtons "static buttons", a map of buttons usually defined once and part of the main UI
	SButtons map[string]*Button

	// AButtons "anonymous buttons", array of buttons usually created on the fly and handled differently
	AButtons []*Button
)

// TODO this should return some kind of tile build status object, e.g has building, can build on, etc
func IsTileSelectionValid() bool {

	// TODO check this works for a non-square world
	return world.things[mtx][mty].completed ||
		(mtx > 0 && world.things[mtx-1][mty].completed) ||
		(mty > 0 && world.things[mtx][mty-1].completed) ||
		(mtx < len(world.things)-1 && world.things[mtx+1][mty].completed) ||
		(mtx < len(world.things) && mty < len(world.things[mtx])-1 && world.things[mtx][mty+1].completed)
}

// ResetFrameState is a handy function that will reset any variables that should not persist between updates
func ResetFrameState() {
	validMouseSelection = false
	mtx = -1
	mty = -1
}

func CreateMessages() MessageQueue {
	q := make([]string, 0)
	return MessageQueue{
		queue: q,
		max:   3,
	}
}

func (m *MessageQueue) AddMessage(message string) {
	queue := append(m.queue, message)
	if len(queue) > m.max {
		queue = queue[1:]
	}
	m.queue = queue
}

func (m *MessageQueue) DrawMessages(screen *ebiten.Image) {

	x := 300
	y := 0

	for i := 0; i < len(Messages.queue); i++ {

		// TODO calculate longest string, justify each message and translate
		// 	the canvas accordingly (instead of using static coordinates)
		canvas := ebiten.NewImage(600, 100)
		text.Draw(canvas, Messages.queue[len(Messages.queue)-1-i], fontDetail, 20, 20, color.White)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y+(i*20)))
		alpha := 1 - (0.3 * float64(i))
		op.ColorM.Scale(1, 1, 1, alpha)
		screen.DrawImage(canvas, op)
	}
}

func UpdateDrawLocations() {

	// north will be top left
	xOffset := 0
	yOffset := 120

	// mouse position must have been updated by now
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
				// can use rect and Point::in for this I think
				if (tx+16 < mxf) && (mxf < tx+48) && (ty+8 < myf) && (myf < ty+24) {

					world.tiles[x][y].selected = true
					mtx, mty = x, y
					mouseFound = true
					validMouseSelection = IsTileSelectionValid()
				} else {
					world.tiles[x][y].selected = false
				}
			} else {
				world.tiles[x][y].selected = false
			}
		}
	}
}

// UpdateInputs calls appropriate functions when inputs detected
func UpdateInputs() {

	// TODO only on mouse release, so a user can cancel by moving the cursor before they release click
	//  see IsMouseButtonJustReleased
	for i := 0; i < len(AllButtons); i++ {
		button := AllButtons[i]

		point := image.Point{
			X: mx,
			Y: my,
		}

		// TODO remove if statement break if true (but not right now as still under development)
		// button.hover = point.In(button.bounds)
		if point.In(button.bounds) {
			button.hover = true
		} else {
			button.hover = false
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

		if validMouseSelection && world.tiles[mtx][mty].category == TGrass {

			if world.things[mtx][mty].nothing {
				world.things[mtx][mty] = Thing{
					animated: &village,
					nothing:  false,
				}
			}
		}

		// civilization.citizens = append(civilization.citizens, Citizen{
		// 	age: 17,
		// })
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {

		if validMouseSelection && world.tiles[mtx][mty].category == TGrass {

			if world.things[mtx][mty].nothing {
				world.things[mtx][mty] = Thing{
					progress:  0,
					completed: false,
					animated:  &house,
					nothing:   false,
				}
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

func HandleTurnEnd() {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if SButtons[BtnEndTurn].hover {
			if epoch+1 >= len(epochs) {
				return
			}
			epoch++
			Messages.AddMessage(fmt.Sprintf("You advanced to the %s", epochs[epoch]))
		}
	}
	return
}

func (g *Game) Update() error {

	ResetFrameState()

	// TODO bind to window. ebiten seems to track outside of the window too
	mx, my = ebiten.CursorPosition()

	// this also finds which tile the mouse is on
	UpdateDrawLocations()

	UpdateInputs()

	// we definitely shouldn't accept any user input after this until the next loop
	HandleTurnEnd()

	// TODO UpdateUi

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

			// draw north arrow. for debugging only atm
			if x == 0 {
				screen.DrawImage(north, world.tiles[x][y].op)
			}
		}
	}

	// TODO dedupe functionality between these two loops
	// draw grass layer // why isn't this i nthe above block?
	for x := 0; x < len(world.tiles); x++ {
		for y := 0; y < len(world.tiles[x]); y++ {
			if world.tiles[x][y].category == TGrass {

				// colour tile differently based on selection
				if world.tiles[x][y].selected {
					if validMouseSelection {
						world.tiles[x][y].op.ColorM.Scale(0.6, 1, 0.6, 1)
					} else {
						world.tiles[x][y].op.ColorM.Scale(1, 0.6, 0.6, 1)
					}
				}

				screen.DrawImage(grass.flat, world.tiles[x][y].op)

				// if the west adjacent tile is lower, draw the west side
				if y == 0 || (world.tiles[x][y-1].category < world.tiles[x][y].category) {
					screen.DrawImage(grass.west, world.tiles[x][y].op)
				}

				// if the south adjacent tile is lower, draw the south side
				if x < len(world.tiles) || (world.tiles[x+1][y].category < world.tiles[x][y].category) {
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
				// constructions in progress will be transparent, with their opacity increasing as they near construction
				if !world.things[x][y].completed {
					// do not animate things under construction as it more clearly indicates that it's not in operation
					world.tiles[x][y].op.ColorM.Scale(1, 1, 1, 0.4)
					screen.DrawImage(world.things[x][y].animated.sprites[0], world.tiles[x][y].op)
				} else {
					screen.DrawImage(world.things[x][y].animated.sprites[world.things[x][y].animated.frame], world.tiles[x][y].op)
				}
				// reset scale in case we changed it
				world.tiles[x][y].op.ColorM.Scale(1, 1, 1, 1)
			}
		}
	}

	// debugging only
	// mx, my := ebiten.CursorPosition()
	// ebitenutil.DrawRect(screen, float64(mx-32), float64(my-16), 64, 32, color.Opaque)
}

// DrawButton handles text and button sizing and positioning
// TODO button state variable

func CreateButton(img *UiSprite, str string, x, y int) (*Button, int) {

	b := Button{
		content: str,
		x:       x,
		y:       y,
		img:     img,
		hover:   false,
	}

	// TODO calculate the text padding/button side size instead of using a magic number
	w := text.BoundString(fontDetail, str).Dx() + 8

	AllButtons = append(AllButtons, &b)
	return &b, w
}

func (b *Button) DrawButton(screen *ebiten.Image) {

	// TODO cache state so we don't need to recalculate if there are no changes.
	// 	use redraw variable for this and a member function for move or update string
	strRect := text.BoundString(fontDetail, b.content)
	strWidth := strRect.Size().X

	// TODO reuse this so we don't have to set scale multiple times
	op := &ebiten.DrawImageOptions{}
	if b.hover {
		op.ColorM.Scale(0.8, 0.8, 0.8, 1)
	}
	op.GeoM.Translate(float64(b.x), float64(b.y))
	screen.DrawImage(btn.left, op)
	w, _ := btn.left.Size()

	op = &ebiten.DrawImageOptions{}
	if b.hover {
		op.ColorM.Scale(0.8, 0.8, 0.8, 1)
	}
	op.GeoM.Scale(float64(strWidth), 1)
	op.GeoM.Translate(float64(b.x+w), float64(b.y))
	screen.DrawImage(btn.middle, op)

	op = &ebiten.DrawImageOptions{}
	if b.hover {
		op.ColorM.Scale(0.8, 0.8, 0.8, 1)
	}
	bx := b.x + w + strWidth
	op.GeoM.Translate(float64(bx), float64(b.y))
	screen.DrawImage(btn.right, op)

	// draw button text
	text.Draw(screen, b.content, fontDetail, b.x+w, b.y+12, color.White)

	// TODO write unit tests for this
	b.bounds = image.Rectangle{
		image.Point{
			b.x,
			b.y,
		},
		image.Point{
			bx + btn.right.Bounds().Bounds().Size().X,
			b.y + btn.right.Bounds().Bounds().Size().Y,
		},
	}
}

// TODO message display
func DrawUi(screen *ebiten.Image) {

	// the font should totally upgrade with each age
	text.Draw(screen, epochs[epoch], fontTitle, 8, 16, color.White)

	// smaller font for more detailed information
	text.Draw(screen, fmt.Sprintf("Citizens: %d/%d", len(civilization.citizens), civilization.max), fontDetail, 8, 30, color.White)

	for i := 0; i < len(AllButtons); i++ {
		AllButtons[i].DrawButton(screen)
	}
	// newer messages should be at the bottom of the screen and older messages should fade
	Messages.DrawMessages(screen)
}

func (g *Game) Draw(screen *ebiten.Image) {
	// render
	screen.Fill(color.Black)

	DrawWorld(screen, &world)
	DrawUi(screen)

	// TODO if debug
	// TODO don't calculate mouse pos on the draw call. this is for debugging only
	mx, my := ebiten.CursorPosition()
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse pos: %d,%d", mx, my), 16, 60)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse on tile: %d, %d", mtx, mty), 16, 80)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Cursor on tile: %d, %d", ctx, cty), 16, 100)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	sWidth = outsideWidth / 2
	sHeight = outsideHeight / 2
	return sWidth, sHeight
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
	tiles := IslandWorldTiles()

	w := World{
		tiles: tiles,
	}

	w.CreateThings()
	return w
}

func AddDebugThings(t [][]Thing) {
	t[7][7] = Thing{
		completed: true,
		animated:  &village,
		nothing:   false,
	}
	t[0][0] = Thing{
		completed: true,
		animated:  &village,
		nothing:   false,
	}
}

func CreateSpawnSettlement() Thing {

	c := []Citizen{}
	for i := 0; i < int(math.Floor(float64(SettlementCapacityVillage))); i++ {
		c = append(c, Citizen{age: 18})
	}

	return Thing{
		completed: true,
		animated:  &village,
		nothing:   false,
	}
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

	// spawn village in the middle(ish) of the map
	t[3][4] = CreateSpawnSettlement()

	//AddDebugThings(t)
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

func CreateUi() {

	Messages = CreateMessages()
	SButtons = make(map[string]*Button)

	// "static" button
	var bw int
	bx := 6
	// this won't work as sHeight hasn't been set yet. it's set when the game is run,
	// so you may have to conditionally run Init() at the top of the update function.
	// maybe you could put this in Layout()?
	// using the expected values for now
	by := (WHeight / 2) - 22
	SButtons[BtnEndTurn], bw = CreateButton(&btn, "End turn", bx, by)
	bx += bw
	SButtons[BtnShowBuildings], bw = CreateButton(&btn, "Buildings", bx, by)
	bx += bw
	// "anonymous" button
	CreateButton(&btn, "BALLS BALLS BALLS", bx, by)
}

// LoadUISprite assumes that the path contains left.png, middle.png and right.png
func LoadUISprite(path string) UiSprite {

	left, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("%s/left.png", path))
	if err != nil {
		log.Fatal(err)
	}

	middle, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("%s/middle.png", path))
	if err != nil {
		log.Fatal(err)
	}

	right, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("%s/right.png", path))
	if err != nil {
		log.Fatal(err)
	}

	return UiSprite{
		left:   left,
		middle: middle,
		right:  right,
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

	// load north sprite
	north, _, err = ebitenutil.NewImageFromFile("img/tiles/north.png")
	if err != nil {
		log.Fatal(err)
	}

	btn = LoadUISprite("img/ui/button")
}

func Init() {
	LoadFonts()
	LoadSprites()
	CreateUi()
	world = CreateWorld()
	civilization = CreateCivilization()
}

func main() {
	fmt.Println("hello world")

	Init()

	// do this with a function. it's to make the screen size fit the map (assuming 8x8) like minesweeper
	ebiten.SetWindowSize(WWidth, WHeight)
	ebiten.SetWindowTitle("Kingdom")

	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
