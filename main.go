package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"

	_ "image/png"

	"github.com/hajimehoshi/ebiten"
)

type Game struct{}

type Citizen struct {
	name      string
	age       int
	gender    string
	education int
	genetics  int
	// TODO home settlement, tile on last turn
	// home settlement could provide a buff to effort
}

// UI elements
type UiSprite struct {
	left *ebiten.Image
	// middle MUST be one pixel wide for text fit scaling
	middle *ebiten.Image
	right  *ebiten.Image
}

type Window struct {
	// dimensions
	width  int
	height int
	// window image
	canvas *ebiten.Image
	// position
	px, py float64
	// has there been a change or has the window just spawned? if so, redraw
	redraw bool
}

// likely need to use embedding for dynamic UIs
type SettlementUi struct {
	window *Window
	// selected settlement coordinates
	sx, sy int
	// has the user got something selected?
	focused bool
	// has there been a change or has the window just spawned? if so, redraw
	redraw bool
	// selectCtzButtons
	selectCtzButtons []*Button
	selectedCtz      *Citizen
	// selectJobButtons
	selectJobButtons []*Button
	selectedJob      *Job
	// jobs
	jobs []*Job
}

type Button struct {
	redraw     bool
	selected   bool
	content    string
	x, y       int
	width      int
	hover      bool // should default to false
	img        *UiSprite
	bounds     image.Rectangle
	windowed   bool
	window     *Window
	disabled   bool
	executable bool
	exec       func() string
	destroy    bool
	cachedImg  *ebiten.Image
}

type Message struct {
	content string
	dupes   int
}

type MessageQueue struct {
	// TODO duplicate message indicator
	queue []*Message
	max   int
}

// Tiles
type Tile struct {
	kind        int
	building    int
	selected    bool
	moved       bool
	height      int
	liquid      bool
	hovered     bool
	highlighted bool
	// cache
	tx       float64
	ty       float64
	opsFlat  *ebiten.DrawImageOptions
	opsWest  *ebiten.DrawImageOptions
	opsSouth *ebiten.DrawImageOptions
}

// TODO TileType struct?
type TileSprite struct {
	flat     *ebiten.Image
	south    *ebiten.Image
	west     *ebiten.Image
	southMid *ebiten.Image
	westMid  *ebiten.Image
}

// Animations
type Animation struct {
	frame   int
	sprites []*ebiten.Image
}

func (a *Animation) Animate() {
	a.frame++
	if a.frame == len(a.sprites) {
		a.frame = 0
	}
}

// World
type SettlementKind struct {
	name      string
	effort    float64
	animation Animation
	nothing   bool
	popcap    int
}

type Settlement struct {
	kind      *SettlementKind
	progress  float64
	completed bool
	citizens  []Citizen
}

type Work struct {
	x, y int
}

type Job struct {
	x, y int
	kind string
	work *Work
}

type World struct {
	tiles          [][]Tile
	settlementList []*Settlement
	settlementGrid [][]*Settlement
	xOffset        int
	yOffset        int
	redraw         bool
	cachedImg      *ebiten.Image
}

type Research struct {
	// Husbandry when researched, allows moving a citizen to an adjacent tile without having to wait a turn
	Husbandry bool
	// Transit when researched, allows moving a citizen to any tile without having to wait a turn
	Transit bool
}

const (
	Scale = 2
	// MaxMemAlloc maximum MiB we want to allow to be allocated before we crash the program
	MaxMemAlloc = 128
	// WindowWidth is the minimum supported window width
	WindowWidth = 1024
	// WindowHeight is the minimum supported window height
	WindowHeight = 768
	// TWater water tile type index
	TWater = 0
	// TGrass grass tile type index
	TGrass = 1
	// TileWidth width of tiles in pixels (unscaled)
	TileWidth = 64
	// TileHeight height of tiles in pixels (unscaled)
	TileHeight = 32
	// BtnEndTurn is the button map key for ending a turn
	BtnEndTurn       = "END_TURN"
	BtnShowBuildings = "SHOW_BUILDINGS"
)

var (

	// constant vars (they're vars but we treat them as constants. see defs())
	settlementKinds map[string]*SettlementKind
	nothing         Settlement
	epochs          = []string{
		"Neolithic Age", "Roman Age", "Classical Age",
		"Age of Steam", "Modern Age", "Transhuman Age",
		"Apocalyptic Age",
	}
	tileSprites map[string]TileSprite

	// meta game state
	initialised bool

	// world images
	tilesLayer     *ebiten.Image
	thingsLayer    *ebiten.Image
	highlightLayer *ebiten.Image
	uiLayer        *ebiten.Image

	// ui stuff
	fontTitle    font.Face
	fontDetail   font.Face
	fontSmall    font.Face
	btn          UiSprite
	settlementUi SettlementUi

	// actual vars now
	sHeight   int
	sWidth    int
	lastFrame int = 0
	ticks     int = 0
	year      int = 1
	epoch     int = 0
	world     World
	research  Research
	north     *ebiten.Image
	highlight *ebiten.Image

	// ctx and cty are the coordinate of the tile that the cursor is on
	ctx                 int  = 0
	cty                 int  = 0
	mx                  int  = 0
	my                  int  = 0
	mouseMoved          bool = true
	mtx                 int  = -1
	mty                 int  = -1
	validMouseSelection bool = false
	focusedSettlementX  int  = -1
	focusedSettlementY  int  = -1

	messages MessageQueue

	// AllButtons is the master list of buttons. Used by renderer and mouse picker
	AllButtons []*Button

	// SButtons "static buttons", a map of buttons usually defined once and part of the main UI
	SButtons map[string]*Button

	// AButtons "anonymous buttons", array of buttons usually created on the fly and handled differently
	AButtons []*Button

	// debugging
	debug             bool
	renderTilesLayer  bool
	renderThingsLayer bool
)

// TODO this should return some kind of tile build status object, e.g has building, can build on, etc
//  or perhaps adjacent tiles
func IsTileSelectionValid() bool {

	// TODO check this works for a non-square world
	return world.settlementGrid[mtx][mty].completed ||
		(mtx > 0 && world.settlementGrid[mtx-1][mty].completed) ||
		(mty > 0 && world.settlementGrid[mtx][mty-1].completed) ||
		(mtx < len(world.settlementGrid)-1 && world.settlementGrid[mtx+1][mty].completed) ||
		(mtx < len(world.settlementGrid) && mty < len(world.settlementGrid[mtx])-1 && world.settlementGrid[mtx][mty+1].completed)
}

// ResetFrameState is a handy function that will reset any variables that should not persist between updates
func ResetFrameState() {
	// remove any buttons marked for removal
	for i := 0; i < len(AllButtons); i++ {
		if AllButtons[i].destroy {
			// this will not preserve order. hope that's ok
			if len(AllButtons) > 1 {
				last := len(AllButtons) - 1
				// move the last button in the array to the index of the one we want to remove
				AllButtons[i] = AllButtons[last]
				// remove the last element of the array
				AllButtons = AllButtons[:last]
			} else {
				// if there will be no elements left in the array, just create a new array
				AllButtons = []*Button{}
			}
		}
	}
	validMouseSelection = false
	mtx = -1
	mty = -1
}

func CreateMessages() MessageQueue {
	q := make([]*Message, 0)
	return MessageQueue{
		queue: q,
		max:   8,
	}
}

func (m *Message) ToString() string {
	if m.dupes > 0 {
		return fmt.Sprintf("%s (%d)", m.content, m.dupes+1)
	}
	return m.content
}

func (m *MessageQueue) AddMessage(content string) {

	// get the last message. if the content matches, increment the value
	if len(m.queue) > 0 {
		lm := m.queue[len(m.queue)-1]
		if lm.content == content {
			lm.dupes++
			return
		}
	}

	queue := append(m.queue, &Message{
		content: content,
		dupes:   0,
	})
	if len(queue) > m.max {
		// dequeue
		queue = queue[1:]
	}
	m.queue = queue
}

func (m *MessageQueue) DrawMessages(screen *ebiten.Image) {

	x := 300
	y := 0

	for i := 0; i < len(messages.queue); i++ {

		// TODO calculate longest string, justify each message and translate
		// 	the canvas accordingly (instead of using static coordinates)
		canvas := ebiten.NewImage(600, 100)
		text.Draw(canvas, messages.queue[len(messages.queue)-1-i].ToString(), fontDetail, 20, 20, color.White)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y+(i*14)))
		alpha := 1 - (0.3 * float64(i))
		op.ColorM.Scale(1, 1, 1, alpha)
		screen.DrawImage(canvas, op)
	}
}

// TODO rename to UpdateWorld or something
func UpdateDrawLocations() {

	// north will be top left
	xOffset := 0
	yOffset := 120

	// mouse position must have been updated by now
	mxf, myf := float64(mx), float64(my)
	// TODO height offset for higher tiles
	mouseFound := false

	for x := 0; x < len(world.tiles); x++ {
		for y := 0; y < len(world.tiles[x]); y++ {

			tile := &world.tiles[x][y]
			tx := tile.tx
			ty := tile.ty

			if tile.moved {

				// TODO use tile width vars or consts
				// tx and ty are where the tile will be drawn from
				// arguably this may not belong in here
				// Recalculating tile position on screen
				tx = float64(xOffset) + float64(y*32) + float64(x*32)
				ty = float64(yOffset) - float64(16*y) + float64(x*16)
				ty = ty - float64(tile.height)

				// update dem positions
				tile.tx = tx
				tile.ty = ty
			}

			tile.hovered = false
			if !mouseFound {
				// this matches a box in the centre of the sprite. needs to actually fit the iso
				// if you treat what the player sees as a rectangle, it won't work correctly
				// can use rect and Point::in for this I think
				// I'd like to evaluate all this and find the one with the pointer closest to the center tbh as a long term solution
				if (tx+16 < mxf) && (mxf < tx+48) && (ty+8 < myf) && (myf < ty+24) {

					tile.hovered = true
					world.redraw = true
					mtx, mty = x, y
					mouseFound = true
					validMouseSelection = IsTileSelectionValid()
				}
			}
		}
	}
}

func DefocusSettlement() {
	if settlementUi.focused {
		HighlightAvailableTiles(settlementUi.sx, settlementUi.sy, false)
		ClearSettlementUi()
	}
}

// TODO consider making this "select settlement or something. focused might be a bit ambiguous"
// 	also, it might not be a settlement, it could be a building. Settlements are also buildings
// 	maybe it should even be focus tile? could then remove the buildings button...
// 	buildings could pop up on an empty tile?
// FocusSettlement returns true if a new settlement has come into focus
func FocusSettlement(x, y int) bool {
	if settlementUi.focused && settlementUi.sx == x && settlementUi.sy == y {
		// do nothing
		return false
	}

	// TODO move UI with mouse
	settlementUi.sx = x
	settlementUi.sy = y

	fmt.Println(fmt.Sprintf("Focused %d,%d", x, y))

	return true
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
			// only run this logic if the button is not already hovered over
			if !button.hover {
				button.hover = true
				button.SetRedraw()
			}
		} else {
			// if we are hovering already, set hover to false and mark the window for redraw if necessary
			if button.hover {
				button.hover = false
				button.SetRedraw()
			}
		}
	}

	justFocusedSettlement := false
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

		HandleButtonClicks()

		if settlementUi.focused && validMouseSelection && settlementUi.selectedCtz != nil && world.tiles[mtx][mty].highlighted {
			fmt.Println(fmt.Sprintf("Assigned citizen %s to %s", settlementUi.selectedCtz.name, world.tiles[mtx][mty].kind))
		} else if validMouseSelection && world.tiles[mtx][mty].kind == TGrass {

			clickedSettlement := world.settlementGrid[mtx][mty]

			if clickedSettlement.kind.nothing {
				// TODO instead spawn the buildings UI
				world.settlementGrid[mtx][mty] = world.CreateSettlement(settlementKinds["VILLAGE"])
			} else if !clickedSettlement.completed {
				DefocusSettlement() // TODO change to defocus selection? idk
			} else {
				DefocusSettlement()
				justFocusedSettlement = FocusSettlement(mtx, mty)
			}
		}
	}

	if justFocusedSettlement {
		HighlightAvailableTiles(mtx, mty, true)
		CreateSettlementUi()
	} else {
		UpdateSettlementUi()
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {

		if validMouseSelection && world.tiles[mtx][mty].kind == TGrass {

			if world.settlementGrid[mtx][mty].kind.nothing {
				world.settlementGrid[mtx][mty] = world.CreateSettlement(settlementKinds["SUBURB"])
			}
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		// TODO destroy settlement. hopefully go gc is good
		world.settlementGrid[mtx][mty] = &nothing
	}

	// move cursor north
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		DefocusSettlement()
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
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		fmt.Println("Thanks for playing")
		os.Exit(0)
	}

	// TODO ignore tile hover/click if blocked by a UI

	// debugging
	if inpututil.IsKeyJustPressed(ebiten.KeyGraveAccent) {
		fmt.Println("Debugging toggled. Use '\\' to toggle again")
		debug = !debug
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		renderTilesLayer = !renderTilesLayer
		if renderTilesLayer {
			fmt.Println("Tiles layer shown")
		} else {
			fmt.Println("Tiles layer hidden")
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		renderThingsLayer = !renderThingsLayer
		if renderThingsLayer {
			fmt.Println("Things layer shown")
		} else {
			fmt.Println("Things layer hidden")
		}
	}
}

func HandleButtonClicks() {
	// could actually store hovered button as a variable, but this will do for now
	for _, b := range AllButtons {
		if b.hover {
			if b.executable {
				fmt.Println(b.exec())
			}
		}
	}
}

func HandleTurnEnd() {

	if !(inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && SButtons[BtnEndTurn].hover) {
		return
	}

	year++

	// every ten years for now
	if year%10 == 0 && epoch+1 < len(epochs) {
		epoch++
		messages.AddMessage(fmt.Sprintf("You advanced to the %s", epochs[epoch]))
	}

	// iterate through constructions
	for i := 0; i < len(world.settlementList); i++ {
		s := world.settlementList[i]
		if !s.kind.nothing {
			if !s.completed {
				// TODO use manpower of adjacent settlement. obviously this will be
				//  a problem with multiple adjacent builds/settlements, for future
				// 	jakey
				s.progress += s.kind.effort
				if s.progress >= 1 {
					s.completed = true
					messages.AddMessage(fmt.Sprintf("Construction completed on '%s'", s.kind.name))
				}
			}
		}
	}
}

func (g *Game) Update() error {

	// initialise the game state if it's not
	if initialised {
		ResetFrameState()
	} else {
		initialised = true
		Init()
	}

	// TODO clamp to window. ebiten seems to track outside of the window too
	// 	i.e if mx < 0, mx = 0,
	if ebiten.IsFocused() {
		newMx, newMy := ebiten.CursorPosition()
		if newMx == mx && newMy == my {
			mouseMoved = false
		} else {
			mx, my = newMx, newMy
			mouseMoved = true
		}
	}

	// this also finds which tile the mouse is on
	UpdateDrawLocations()

	UpdateInputs()

	UpdateSettlementUi()

	// we definitely shouldn't accept any user input after this until the next loop
	HandleTurnEnd()

	// TODO delta
	// if the game is running at half speed, the delta should be 2
	// if the game is running at normal speed, the delta should be 1 etc
	if ticks == 0 {

		MonitorMemory()

		// then resume with game stuff
		settlementKinds["VILLAGE"].animation.Animate()
		settlementKinds["SUBURB"].animation.Animate()
	}
	// possible to use a float here for proper delta time?
	ticks++

	// will mean we update twice per second
	if ticks > 30 {
		ticks = 0
	}

	return nil
}

// Am I doing this right?
func Copy(g *ebiten.GeoM) ebiten.GeoM {
	geom := ebiten.GeoM{}
	geom.Translate(g.Element(0, 2), g.Element(1, 2))
	return geom
}

func GetT(g *ebiten.GeoM) (float64, float64) {
	return g.Element(0, 2), g.Element(1, 2)
}

// TODO consider making this a function of Tile
// 	although tiles do need the context of surrounding tiles provided by world
func DrawTile(colour *ebiten.ColorM, layer *ebiten.Image, world *World, ttype string, x, y int, highlighted bool) {

	tile := &world.tiles[x][y]

	if tile.moved || tile.hovered {

		const extraTileHeight = 16

		geomFlat := &ebiten.GeoM{}
		geomFlat.Translate(tile.tx, tile.ty)

		geomWest := &ebiten.GeoM{}
		// order is important. scale _then_ translate
		geomWest.Scale(1, float64(tile.height+extraTileHeight))
		geomWest.Translate(tile.tx, tile.ty+16) // magic number

		geomSouth := &ebiten.GeoM{}
		geomSouth.Scale(1, float64(tile.height+extraTileHeight)) // TODO something else so it goes _below_ the neighbouring tiles if applicable
		geomSouth.Translate(tile.tx+30, tile.ty+16)

		opsFlat := &ebiten.DrawImageOptions{
			GeoM:   *geomFlat,
			ColorM: *colour,
		}
		opsWest := &ebiten.DrawImageOptions{
			GeoM:   *geomWest,
			ColorM: *colour,
		}
		opsSouth := &ebiten.DrawImageOptions{
			GeoM:   *geomSouth,
			ColorM: *colour,
		}

		// cache
		tile.opsFlat = opsFlat
		tile.opsWest = opsWest
		tile.opsSouth = opsSouth
	} else if !tile.hovered {
		tile.opsFlat.ColorM.Reset()
		tile.opsWest.ColorM.Reset()
		tile.opsSouth.ColorM.Reset()
	}

	if y == 0 || (world.tiles[x][y-1].height < world.tiles[x][y].height) {
		if ttype != "water" {
			layer.DrawImage(tileSprites[ttype].westMid, tile.opsWest)
		}
		layer.DrawImage(tileSprites[ttype].west, tile.opsFlat)
	}

	// if the south adjacent tile is lower, draw the south side
	if x < len(world.tiles) || (world.tiles[x+1][y].height < world.tiles[x][y].height) {
		if ttype != "water" {
			layer.DrawImage(tileSprites[ttype].southMid, tile.opsSouth)
		}
		layer.DrawImage(tileSprites[ttype].south, tile.opsFlat)
	}

	layer.DrawImage(tileSprites[ttype].flat, tile.opsFlat)
	// TODO add TGrass property to tile. should be able to loop through tile types
}

func DrawWorld(layer *ebiten.Image, world *World) {

	// don't redraw if map doesn't change between frames
	// layer.Clear()

	// north is top left
	// might want to store this layer globally or something in between frames and reuse it
	if world.redraw || world.cachedImg == nil {

		world.cachedImg = ebiten.NewImage(sWidth, sHeight)

		for x := 0; x < len(world.tiles); x++ {
			for y := len(world.tiles[x]) - 1; y > -1; y-- {

				tile := &world.tiles[x][y]

				var ttype string
				colour := &ebiten.ColorM{}

				// tile type specific shading
				if tile.kind == TWater {

					ttype = "water"

					if x == ctx && y == cty {
						colour.Scale(0.6, 1, 0.6, 1)
					}

				} else if tile.kind == TGrass {

					ttype = "grass"

					// colour tile differently based on selection
					if tile.hovered && (!settlementUi.focused || tile.highlighted) {
						if validMouseSelection {
							colour.Scale(0.6, 1, 0.6, 1)
						} else {
							colour.Scale(1, 0.6, 0.6, 1)
						}
					}
				}

				DrawTile(colour, world.cachedImg, world, ttype, x, y, tile.highlighted)

				// we're done with the tile move state. on to the next frame
				tile.moved = false
			}
		}
	}

	layer.DrawImage(world.cachedImg, &ebiten.DrawImageOptions{})
	world.redraw = false

	// debugging only
	// mx, my := ebiten.CursorPosition()
	// ebitenutil.DrawRect(screen, float64(mx-32), float64(my-16), 64, 32, color.Opaque)
}

func DrawThings(layer *ebiten.Image) {

	layer.Clear()

	for x := 0; x < len(world.settlementGrid); x++ {
		for y := 0; y < len(world.settlementGrid[x]); y++ {
			if !world.settlementGrid[x][y].kind.nothing {
				s := world.settlementGrid[x][y]
				// constructions in progress will be transparent, with their opacity increasing as they near construction

				var frame *ebiten.Image
				ops := &ebiten.DrawImageOptions{}
				ops.GeoM.Translate(world.tiles[x][y].tx, world.tiles[x][y].ty)

				// TODO i broke !s.completed scaling
				if !s.completed {
					ops.ColorM.Scale(1, 1, 1, 0.4)
					// do not animate things under construction as it more clearly indicates that it's not in operation
					frame = s.kind.animation.sprites[0]
				} else {
					frame = s.kind.animation.sprites[world.settlementGrid[x][y].kind.animation.frame]
				}

				layer.DrawImage(frame, ops)
			}
		}
	}
}

// TODO don't draw if not required. skip any draw calls
func DrawHighlightLayer(layer *ebiten.Image) {

	layer.Clear()

	// TODO redraw property
	for x := 0; x < len(world.tiles); x++ {
		for y := len(world.tiles[x]) - 1; y > -1; y-- {
			tile := &world.tiles[x][y]
			if tile.highlighted {

				geom := Copy(&tile.opsFlat.GeoM)
				ops := &ebiten.DrawImageOptions{
					GeoM: geom,
				}

				if tile.selected {
					// orange
					const div = 255
					r, g, b := float64(255), float64(191), float64(0)
					dr, dg, db := r/div, g/div, b/div
					ops.ColorM.Scale(dr, dg, db, 1)
				}

				layer.DrawImage(highlight, ops)

				settlement := world.settlementGrid[x][y]
				if !settlement.completed {
					words := fmt.Sprintf("%.1f %%", settlement.progress*100)
					// TODO show progress on next turn
					width := text.BoundString(fontSmall, words).Dx()
					text.Draw(layer, words, fontSmall, int(tile.tx)+32-(width/2), int(tile.ty)+16, color.White)
				}

			}
		}
	}
}

func DrawLayers(screen *ebiten.Image) {

	if renderTilesLayer {
		// TODO if defocused, i.e saved or blocking dialogue, dim
		screen.DrawImage(tilesLayer, &ebiten.DrawImageOptions{})
	}
	if renderThingsLayer {
		// TODO if defocused, i.e saved or blocking dialogue, dim
		screen.DrawImage(thingsLayer, &ebiten.DrawImageOptions{})
	}
	screen.DrawImage(highlightLayer, &ebiten.DrawImageOptions{})
	screen.DrawImage(uiLayer, &ebiten.DrawImageOptions{})
}

func (c *Citizen) ToString() string {
	return fmt.Sprintf("%s, %s, %d, no job", c.name, c.gender, c.age)
}

func (c *Citizen) ToTerseString() string {
	// %s for strings, %c for chars
	return fmt.Sprintf("Citizen, %c%d", strings.ToUpper(c.gender)[0], c.age)
}

// TODO button state variable
// CreateButton appends it to the global buttons list, returns the button and the text width (TODO maybe make this whole button width?)
func CreateButton(img *UiSprite, str string, x, y int) (*Button, int) {

	b := Button{
		content: str,
		// TODO remove coordinates and move into draw cycle
		x:         x,
		y:         y,
		img:       img,
		hover:     false,
		selected:  false,
		destroy:   false,
		disabled:  false,
		redraw:    true,
		cachedImg: nil,
	}

	// TODO calculate the text padding/button side size instead of using a magic number
	// 	i.e use &btn.left.Dx()...
	w := text.BoundString(fontDetail, str).Dx() + 8
	b.width = w

	AllButtons = append(AllButtons, &b)
	return &b, w
}

// SetWindow assigns a button to a window. Useful for redraw logic, i.e hover
func (b *Button) SetWindow(window *Window) {
	b.windowed = true
	b.window = window
}

// SetRedraw instructs button and its window (if set) to redraw
func (b *Button) SetRedraw() {
	if b.windowed && b.window != nil {
		b.window.redraw = true
	}
	b.redraw = true
}

// DrawButtonAt draws button at given coordinates
func (b *Button) DrawButtonAt(layer *ebiten.Image, x, y int) {

	// oldX := b.x
	// oldY := b.y

	b.x = x
	b.y = y

	b.DrawButton(layer)

	// b.x = oldX
	// b.y = oldY
}

// DrawButton draws button at stored button coordinates. If you wish to
// 	specify otherwise, use DrawButtonAt
func (b *Button) DrawButton(layer *ebiten.Image) {

	var totalWidth int
	if b.redraw || b.cachedImg == nil {

		lw, lh := btn.left.Size()
		strWidth := text.BoundString(fontDetail, b.content).Size().X
		rw, _ := btn.right.Size() // what would be rh should be same size as lh

		totalWidth = lw + strWidth + rw

		buttonImg := ebiten.NewImage(totalWidth, lh)
		textImg := ebiten.NewImage(totalWidth, lh)
		b.cachedImg = ebiten.NewImage(totalWidth, lh)

		ops := &ebiten.DrawImageOptions{}
		buttonImg.DrawImage(btn.left, ops)

		// scale 1px wide middle button body and translate
		ops.GeoM.Scale(float64(strWidth), 1)
		ops.GeoM.Translate(float64(lw), 0)
		buttonImg.DrawImage(btn.middle, ops)

		// create new ops to get rid of the scaling we just did
		ops = &ebiten.DrawImageOptions{}
		rx := lw + strWidth
		ops.GeoM.Translate(float64(rx), 0)
		buttonImg.DrawImage(btn.right, ops)

		// draw button text
		text.Draw(textImg, b.content, fontDetail, lw, 12, color.White)

		// combine image and text
		ops = &ebiten.DrawImageOptions{}
		if b.selected {
			ops.ColorM.Scale(1, 2, 0, 1)
		}

		if b.hover && !b.disabled {
			ops.ColorM.Scale(0.8, 0.8, 0.8, 1)
		}
		b.cachedImg.DrawImage(buttonImg, ops)

		ops = &ebiten.DrawImageOptions{}
		if b.disabled {
			ops.ColorM.Scale(0.8, 0.8, 0.8, 1)
		}
		b.cachedImg.DrawImage(textImg, ops)
		b.redraw = false
	} else {
		totalWidth, _ = b.cachedImg.Size()
	}

	// TODO write unit tests for this
	windowOffsetX := 0
	windowOffsetY := 0

	if b.windowed {
		windowOffsetX += int(b.window.px)
		windowOffsetY += int(b.window.py)
	}

	ops := &ebiten.DrawImageOptions{}
	ops.GeoM.Translate(float64(b.x), float64(b.y))
	layer.DrawImage(b.cachedImg, ops)

	b.bounds = image.Rectangle{
		image.Point{
			b.x + windowOffsetX,
			b.y + windowOffsetY,
		},
		image.Point{
			b.x + windowOffsetX + totalWidth,
			b.y + windowOffsetY + btn.right.Bounds().Bounds().Size().Y,
		},
	}
}

func DrawUi(layer *ebiten.Image) {

	layer.Clear()

	// the font should totally upgrade with each age
	text.Draw(layer, epochs[epoch], fontTitle, 8, 16, color.White)

	// smaller font for more detailed information
	// TODO cache this value in update
	// TODO previous frame state (so we can avoid unnecessary calculations)
	civs := 0
	for i := 0; i < len(world.settlementList); i++ {
		civs += len(world.settlementList[i].citizens)
	}
	text.Draw(layer, fmt.Sprintf("Citizens: %d", civs), fontDetail, 8, 30, color.White)
	text.Draw(layer, fmt.Sprintf("Year: %d", year), fontDetail, 8, 44, color.White)

	for _, v := range SButtons {
		v.DrawButton(layer)
	}
	// newer messages should be at the bottom of the screen and older messages should fade
	messages.DrawMessages(layer)

}

func HighlightAvailableTiles(x, y int, highlighted bool) {

	// TODO don't forget array bounds
	// TODO use a different or simpler type
	works := make(map[string]*Work)
	works["here"] = &Work{x: x, y: y}
	works["north"] = &Work{x: x - 1, y: y}
	works["east"] = &Work{x: x, y: y + 1}
	works["south"] = &Work{x: x + 1, y: y}
	works["west"] = &Work{x: x, y: y - 1}

	world.tiles[x][y].selected = highlighted
	for _, v := range works {
		settlement := world.settlementGrid[v.x][v.y]
		tile := &world.tiles[v.x][v.y]

		if settlement.kind.nothing {
			continue
		} else {
			tile.highlighted = highlighted
		}
	}
}

func GetAvailableJobs(x, y int) []*Job {

	// TODO don't forget array bounds
	works := make(map[string]*Work)
	works["here"] = &Work{x: x, y: y}
	works["north"] = &Work{x: x - 1, y: y}
	works["east"] = &Work{x: x, y: y + 1}
	works["south"] = &Work{x: x + 1, y: y}
	works["west"] = &Work{x: x, y: y - 1}

	// TODO jobs type and jobs list (assigned, etc)
	jobs := []*Job{}
	here := world.settlementGrid[x][y]

	// if unfinished settlement, no jobs should be available
	if !here.completed {
		return jobs
	}

	// TODO tool tip of job info
	// TODO warn on excess effort
	for k, v := range works {

		settlement := world.settlementGrid[v.x][v.y]

		if settlement.kind.nothing {
			continue
		}

		if settlement.completed {

			// don't show move for "here"
			if v.x != x || v.y != y {
				jobs = append(jobs, &Job{
					kind: fmt.Sprintf("move %s", k),
					work: v,
				})
			}
		} else if v.x != x || v.y != y {
			jobs = append(jobs, &Job{
				kind: fmt.Sprintf("build %s", k),
				work: v,
			})
		}
	}

	return jobs
}

func (ui *SettlementUi) SelectCitizen(idx int) {

	for i := 0; i < len(ui.selectCtzButtons); i++ {
		b := ui.selectCtzButtons[i]
		if i == idx {
			b.selected = !b.selected
		} else {
			b.selected = false
		}
		b.SetRedraw()
	}

	settlementUi.DisableJobSelection(!ui.selectCtzButtons[idx].selected)
}

func (ui *SettlementUi) DisableJobSelection(disabled bool) {
	for i := 0; i < len(ui.selectJobButtons); i++ {
		j := ui.selectJobButtons[i]
		j.disabled = disabled
		j.SetRedraw()
	}
}

func ClearSettlementUi() {
	settlementUi.focused = false
	settlementUi.selectedCtz = nil
	// TODO may also need to do this in UpdateSettlementUi
	for i := 0; i < len(settlementUi.selectCtzButtons); i++ {
		settlementUi.selectCtzButtons[i].destroy = true
	}
	for i := 0; i < len(settlementUi.selectJobButtons); i++ {
		settlementUi.selectJobButtons[i].destroy = true
	}
	settlementUi.selectCtzButtons = []*Button{}
	settlementUi.selectJobButtons = []*Button{}
	fmt.Println(fmt.Sprintf("Defocused"))
}

func CreateSettlementUi() {

	settlementUi.focused = true
	settlementUi.redraw = true
	settlementUi.selectCtzButtons = []*Button{}
	settlementUi.selectJobButtons = []*Button{}

	// citizens
	settlement := world.settlementGrid[settlementUi.sx][settlementUi.sy]
	for i := 0; i < len(settlement.citizens); i++ {
		citizenText := settlement.citizens[i].ToTerseString()

		b, _ := CreateButton(&btn, citizenText, 0, 0)
		b.executable = true
		idx := i
		b.exec = func() string {
			if b.selected {
				// we're going to flip this in SelectCitizen. this is a bit shit
				settlementUi.selectedCtz = nil
			} else {
				settlementUi.selectedCtz = &settlement.citizens[idx]
			}
			settlementUi.SelectCitizen(idx)
			return "Selected citizen"
		}
		b.SetWindow(settlementUi.window)
		settlementUi.selectCtzButtons = append(settlementUi.selectCtzButtons, b)
	}

	// jobs
	jobs := GetAvailableJobs(settlementUi.sx, settlementUi.sy)

	for j := 0; j < len(jobs); j++ {
		jobsText := jobs[j].kind
		// TODO make this a function of Window?
		b, _ := CreateButton(&btn, jobsText, 0, 0)
		b.executable = true
		b.disabled = true
		b.exec = func() string {
			if settlementUi.selectedCtz == nil {
				return "No citizen selected"
			}
			return fmt.Sprintf("Assigned job '%s' to %s", jobsText, settlementUi.selectedCtz.name)
		}
		b.SetWindow(settlementUi.window)
		settlementUi.selectJobButtons = append(settlementUi.selectJobButtons, b)
	}

	// might not be necessary to store this
	settlementUi.jobs = jobs
}

func UpdateSettlementUi() {

}

func DrawSettlementUi(screen *ebiten.Image) {
	if settlementUi.focused && (settlementUi.redraw || settlementUi.window.redraw) {

		width := settlementUi.window.width
		height := settlementUi.window.height

		settlement := world.settlementGrid[settlementUi.sx][settlementUi.sy]

		canvas := ebiten.NewImage(width, height)
		canvas.Fill(color.Black)

		titleText := "Manage Settlement"
		titleWidth := text.BoundString(fontTitle, titleText).Dx()
		text.Draw(canvas, titleText, fontTitle, width/2-titleWidth/2, 20, color.White)

		// draw citizens UI
		x := 4
		y := 40

		citizensText := fmt.Sprintf("Citizens (%d)", len(settlement.citizens))
		text.Draw(canvas, citizensText, fontDetail, x, y, color.White)

		y += 10
		if len(settlementUi.selectCtzButtons) == 0 {

			text.Draw(canvas, "No citizens", fontDetail, x, y, color.White)
		} else {
			for i := 0; i < len(settlementUi.selectCtzButtons); i++ {
				b := settlementUi.selectCtzButtons[i]
				b.DrawButtonAt(canvas, x, y)
				y += 20
			}
		}

		// no use for the BALLS button right now
		b, _ := CreateButton(&btn, "BALLS BALLS BALLS", x, height-20)
		b.DrawButton(canvas)

		// draw jobs UI
		x = 100
		y = 40

		text.Draw(canvas, "Jobs", fontDetail, x, y, color.White)

		y += 10
		if len(settlementUi.jobs) == 0 {

			text.Draw(canvas, "No jobs", fontDetail, x, y, color.White)
		} else {
			for i := 0; i < len(settlementUi.selectJobButtons); i++ {
				j := settlementUi.selectJobButtons[i]
				j.DrawButtonAt(canvas, x, y)
				y += 20
			}
		}

		settlementUi.window.canvas = canvas

		ops := &ebiten.DrawImageOptions{}
		ops.ColorM.Scale(1, 1, 1, 0.95)
		ops.GeoM.Translate(settlementUi.window.px, settlementUi.window.py)
	}

	if settlementUi.focused {

		// not necessary every frame. maybe cache and use a moved parameter?
		// 	but it's pretty cheap for now
		ops := &ebiten.DrawImageOptions{}
		ops.ColorM.Scale(1, 1, 1, 0.95)
		ops.GeoM.Translate(settlementUi.window.px, settlementUi.window.py)
		screen.DrawImage(settlementUi.window.canvas, ops)
	}

	settlementUi.redraw = false
	settlementUi.window.redraw = false
}

func (g *Game) Draw(screen *ebiten.Image) {
	// render
	screen.Fill(color.Black)

	DrawWorld(tilesLayer, &world)
	DrawThings(thingsLayer)
	DrawHighlightLayer(highlightLayer)
	DrawUi(uiLayer)
	DrawSettlementUi(uiLayer)
	DrawLayers(screen)

	// TODO don't calculate mouse pos on the draw call. this is for debugging only
	mx, my := ebiten.CursorPosition()
	if debug {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse pos: %d,%d", mx, my), 16, 60)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse on tile: %d, %d", mtx, mty), 16, 80)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Current FPS: %f", ebiten.CurrentFPS()), 16, 100)
	}
}

// Layout is called every frame, which I didn't know until now...
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {

	sWidth = outsideWidth / Scale
	sHeight = outsideHeight / Scale

	if !initialised {
		CreateLayers(sWidth, sHeight)
	}

	return sWidth, sHeight
}

func CreateTile(kind int, height int, liquid bool) Tile {
	return Tile{
		kind:        kind,
		selected:    false,
		hovered:     false,
		moved:       true,
		height:      height,
		liquid:      liquid,
		highlighted: false,
	}
}

func CreateResearch() Research {

	return Research{
		Husbandry: false,
		Transit:   false,
	}
}

func CreateWorld() World {

	w := World{
		tiles:     IslandWorldTiles(),
		redraw:    true,
		cachedImg: nil,
	}

	w.CreateSettlements()

	return w
}

// CreateSettlement creates a settlement add it to the world's settlement
// 	list (if it's not nothing) and return the settlement so it can be
// 	added to the world grid location by the calling code
func (w *World) CreateSettlement(kind *SettlementKind) *Settlement {

	s := &Settlement{
		kind:      kind,
		completed: false,
		progress:  0,
		citizens:  []Citizen{},
	}

	if !kind.nothing {
		w.settlementList = append(w.settlementList, s)
	}

	return s
}

func CreateSpawnSettlement() *Settlement {

	sk := settlementKinds["VILLAGE"]
	c := []Citizen{}

	for i := 0; i < sk.popcap/2; i++ {

		var gender string
		var name string

		if i%2 == 0 {
			gender = "female"
			// right exclusive, neat
			name = FirstNamesFemale[rand.Intn(len(FirstNamesFemale))]
		} else {
			gender = "male"
			name = FirstNamesMale[rand.Intn(len(FirstNamesMale))]
		}

		c = append(c, Citizen{
			name:     name,
			gender:   gender,
			genetics: 100,
			age:      18,
		})
	}

	return &Settlement{
		kind:      sk,
		completed: true,
		citizens:  c,
	}
}

func (w *World) CreateSettlements() {

	grid := [][]*Settlement{}
	list := []*Settlement{}

	for x := 0; x < len(w.tiles); x++ {
		txa := []*Settlement{}
		for y := 0; y < len(w.tiles[x]); y++ {
			txa = append(txa, &nothing)
		}
		grid = append(grid, txa)
	}

	// spawn village in the middle(ish) of the map
	s := CreateSpawnSettlement()
	list = append(list, s)
	grid[3][4] = s

	w.settlementList = list
	w.settlementGrid = grid
}

func CreateLayers(sWidth, sHeight int) {

	tilesLayer = ebiten.NewImage(sWidth, sHeight)
	thingsLayer = ebiten.NewImage(sWidth, sHeight)
	highlightLayer = ebiten.NewImage(sWidth, sHeight)
	uiLayer = ebiten.NewImage(sWidth, sHeight)
}

func CreateUi() {

	settlementUi = SettlementUi{
		// TODO move to bottom right?
		window: &Window{
			width:  150,
			height: 200,
			px:     16,
			py:     16,
			redraw: true,
		},
		sx:      -1,
		sy:      -1,
		redraw:  true,
		focused: false,
	}

	messages = CreateMessages()
	SButtons = make(map[string]*Button)

	// "static" button
	var bw int
	bx := 6
	// this won't work as sHeight hasn't been set yet. it's set when the game is run,
	// so you may have to conditionally run Init() at the top of the update function.
	// maybe you could put this in Layout()?
	// using the expected values for now
	by := (WindowHeight / 2) - 22
	SButtons[BtnEndTurn], bw = CreateButton(&btn, "End turn", bx, by)
	bx += bw
	SButtons[BtnShowBuildings], bw = CreateButton(&btn, "Buildings", bx, by)
	bx += bw
	// "anonymous" button
	CreateButton(&btn, "BALLS BALLS BALLS", bx, by)
}

// LoadUISprite assumes that the path contains left.png, middle.png and right.png
func LoadUISprite(path string) UiSprite {

	left, _, err := ebitenutil.NewImageFromFile(filepath.Join(path, "left.png"))
	if err != nil {
		log.Fatal(err)
	}

	middle, _, err := ebitenutil.NewImageFromFile(filepath.Join(path, "middle.png"))
	if err != nil {
		log.Fatal(err)
	}

	right, _, err := ebitenutil.NewImageFromFile(filepath.Join(path, "right.png"))
	if err != nil {
		log.Fatal(err)
	}

	return UiSprite{
		left:   left,
		middle: middle,
		right:  right,
	}
}

func LoadSprites() {

	// see loading github.com/rakyll/statik in NewImageFromFile documentation
	var err error

	// load north sprite
	north, _, err = ebitenutil.NewImageFromFile(filepath.Join("img", "tiles", "north.png"))
	if err != nil {
		log.Fatal(err)
	}

	highlight, _, err = ebitenutil.NewImageFromFile(filepath.Join("img", "tiles", "highlight.png"))

	// TODO alpha property
	btn = LoadUISprite("img/ui/button")
}

// because we can't use consts for stuff like this
func defs() {

	settlementKinds = make(map[string]*SettlementKind)

	settlementKinds["NOTHING"] = &SettlementKind{
		nothing: true,
		popcap:  0,
	}

	settlementKinds["VILLAGE"] = &SettlementKind{
		name:      "village",
		animation: LoadAnimatedSprite(filepath.Join("img", "sprites", "buildings"), "village", 2),
		popcap:    10,
		nothing:   false,
		// means it will take two person years to construct
		effort: 0.5,
	}

	settlementKinds["SUBURB"] = &SettlementKind{
		animation: LoadAnimatedSprite(filepath.Join("img", "sprites", "buildings"), "house", 2),
		popcap:    20,
		nothing:   false,
		effort:    0.2,
	}

	nothing = Settlement{
		kind: settlementKinds["NOTHING"],
	}

	tileSprites = make(map[string]TileSprite)
	tileSprites["grass"] = LoadTileSprite(filepath.Join("img", "tiles", "grass"))
	tileSprites["water"] = LoadTileSprite(filepath.Join("img", "tiles", "water"))
}

func Init() {
	debug = false
	defs()
	LoadFonts()
	LoadSprites()
	CreateUi()

	// new game
	world = CreateWorld()
	research = CreateResearch()
}

func main() {
	fmt.Println("Starting...")
	initialised = false
	renderTilesLayer = true
	renderThingsLayer = true

	// do this with a function. it's to make the screen size fit the map
	//  (assuming 8x8) like minesweeper
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle("Kingdom")

	// mainly for my development
	ebiten.SetWindowPosition(1400, 200)

	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
