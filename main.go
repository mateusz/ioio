package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/faiface/pixel/imdraw"

	"github.com/lafriks/go-tiled"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/mateusz/ioio/architecture"
	"github.com/mateusz/ioio/graphics"
	"github.com/mateusz/ioio/program"
	"github.com/mateusz/rtsian/piksele"
)

var (
	workDir          string
	monW             float64
	monH             float64
	pixSize          float64
	componentSprites piksele.Spriteset
	p1               player
	gameWorld        piksele.World
	gameSim          program.Simulation
	components       []*architecture.Component
	gameBlips        graphics.BlipList
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var err error
	workDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("Error checking working dir: %s\n", err)
		os.Exit(2)
	}

	gameWorld = piksele.World{}
	gameWorld.Load(fmt.Sprintf("%s/../assets/arch3.tmx", workDir))

	gameBlips = graphics.NewBlipList(&gameWorld)

	loadComponents()
	gameSim = program.NewSimulation(fmt.Sprintf("%s/../examples/prg3.yml", workDir), components, &gameWorld, &gameBlips)

	componentSprites, err = piksele.NewSpritesetFromTsx(fmt.Sprintf("%s/../assets", workDir), "components.tsx")
	if err != nil {
		fmt.Printf("Error loading component sprites: %s\n", err)
		os.Exit(2)
	}

	p1.position = pixel.Vec{
		X: float64(gameWorld.PixelWidth()) / 2.0,
		Y: float64(gameWorld.PixelHeight()) / 2.0,
	}
	p1.scrollSpeed = 200.0
	p1.scrollHotZone = 10.0

	pixelgl.Run(run)
}

func run() {
	monitor := pixelgl.PrimaryMonitor()

	monW, monH = monitor.Size()
	pixSize = 4.0

	cfg := pixelgl.WindowConfig{
		Title:   "Rtsian",
		Bounds:  pixel.R(0, 0, monW, monH),
		VSync:   true,
		Monitor: monitor,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Zoom in to get nice pixels
	win.SetSmooth(false)
	win.SetMatrix(pixel.IM.Scaled(pixel.ZV, pixSize))
	win.SetMousePosition(pixel.Vec{X: monW / 2.0, Y: monH / 2.0})

	mapCanvas := pixelgl.NewCanvas(pixel.R(0, 0, float64(gameWorld.PixelWidth()), float64(gameWorld.PixelHeight())))
	gameWorld.Draw(mapCanvas)

	for _, c := range components {
		componentSprites.Sprites[c.SpriteID].Draw(mapCanvas, pixel.IM.Moved(flipY(c.Position)))
	}

	p1view := pixelgl.NewCanvas(pixel.R(0, 0, monW/pixSize, monH/pixSize))
	blipCanvas := imdraw.New(nil)

	gameSim.Start()

	last := time.Now()
	for !win.Closed() {
		if win.Pressed(pixelgl.KeyEscape) {
			break
		}

		dt := time.Since(last).Seconds()
		last = time.Now()

		// Move player's view
		cam1 := pixel.IM.Moved(pixel.Vec{
			X: -p1.position.X + p1view.Bounds().W()/2,
			Y: -p1.position.Y + p1view.Bounds().H()/2,
		})
		p1view.SetMatrix(cam1)

		// Update world state
		p1.Input(win, cam1)
		p1.Update(dt)

		// Clean up for new frame
		win.Clear(colornames.Black)
		p1view.Clear(colornames.Lightslategray)
		blipCanvas.Clear()

		// Draw transformed map
		mapCanvas.Draw(p1view, pixel.IM.Moved(pixel.Vec{
			X: mapCanvas.Bounds().W() / 2.0,
			Y: mapCanvas.Bounds().H() / 2.0,
		}))

		// Draw blips
		blips := gameBlips.Get()
		for _, b := range blips {
			blipCanvas.Color = b.Color
			p := flipY(b.Pos)
			blipCanvas.Push(p)
			blipCanvas.Push(p.Add(b.Size))
			blipCanvas.Rectangle(0)
		}
		blipCanvas.Draw(p1view)

		// Blit player view
		p1view.Draw(win, pixel.IM.Moved(pixel.Vec{
			X: p1view.Bounds().W() / 2,
			Y: p1view.Bounds().H() / 2,
		}))

		// Present frame!
		win.Update()
	}
}

func loadComponents() {
	for _, o := range gameWorld.Tiles.ObjectGroups[0].Objects {
		lt, err := gameWorld.Tiles.TileGIDToTile(o.GID)
		if err != nil {
			log.Fatal(err)
		}

		p := gameWorld.AlignToTile(pixel.Vec{X: o.X + 10.0, Y: o.Y - 10.0})
		x, y := gameWorld.VecToTile(p)
		// Meh, we shouldn't be flipping at all when loading, just when blitting.
		y = gameWorld.Tiles.Height - y - 1
		tileDef := lt.Tileset.Tiles[lt.ID]

		c := &architecture.Component{
			Position: p,
			X:        x,
			Y:        y,
			SpriteID: lt.ID,
			Sched:    anyProp("sched", o.Properties, tileDef.Properties),
			Con:      anyProp("con", o.Properties, tileDef.Properties),
			Name:     anyProp("name", o.Properties, tileDef.Properties),
		}

		c.Lat, _ = time.ParseDuration(anyProp("lat", o.Properties, tileDef.Properties))
		c.Proc, _ = strconv.Atoi(anyProp("proc", o.Properties, tileDef.Properties))
		c.Cores, _ = strconv.Atoi(anyProp("cores", o.Properties, tileDef.Properties))
		c.Queue, _ = strconv.Atoi(anyProp("queue", o.Properties, tileDef.Properties))
		c.Workers, _ = strconv.Atoi(anyProp("workers", o.Properties, tileDef.Properties))

		components = append(components, c)
	}
}

func anyProp(propName string, p1 tiled.Properties, p2 tiled.Properties) string {
	if p := p1.GetString(propName); p != "" {
		return p
	}

	return p2.GetString(propName)
}

func flipY(pos pixel.Vec) pixel.Vec {
	pos.Y = float64(gameWorld.Tiles.Height*gameWorld.Tiles.TileHeight) - pos.Y
	return pos
}
