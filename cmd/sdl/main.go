package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/nicky-ayoub/imagination/internal/pkg/imagefs"
	"github.com/nicky-ayoub/imagination/internal/pkg/viewport"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SCREEN_WIDTH  int32 = 1440
	SCREEN_HEIGHT int32 = 900
)

type Game struct {
	image *sdl.Surface
	paths []string
	root  string
	name  string
}

func NewGame() *Game {
	g := &Game{}
	seed := time.Now().Unix()
	rand.Seed(seed)
	fmt.Println("Seed : ", seed)
	g.root = "../../assets"
	g.paths = imagefs.AllJpgFiles(g.root)
	fmt.Println(g.paths)
	return g
}

func GetName(g *Game) {
	randomIndex := rand.Intn(len(g.paths))
	g.name = g.paths[randomIndex]
}

func OpenImage(g *Game) (err error) {
	GetName(g)
	g.image, err = img.Load(g.root + "/" + g.name)
	return err
}

func DrawImage(g *Game, window *sdl.Window) (err error) {
	err = OpenImage(g)
	if err != nil {
		return err
	}

	window.SetSize(g.image.W, g.image.H) // Must be called before GetSurface()

	var surface *sdl.Surface
	if surface, err = window.GetSurface(); err != nil {
		return err
	}

	surface.FillRect(&surface.ClipRect, 0)
	// Draw the BMP image on the first half of the window
	g.image.Blit(nil, surface, &g.image.ClipRect)

	// Update the window surface with what we have drawn
	window.UpdateSurface()
	return err
}

func run(g *Game) (err error) {

	var event sdl.Event
	var tick uint32

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		os.Exit(1)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("SDL2 Power", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, SCREEN_WIDTH, SCREEN_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	DrawImage(g, window)
	running := true
	for running {
		event = sdl.WaitEventTimeout(1000) // wait here until an event is in the event queue
		if event == nil {
			fmt.Println("WaitEventTimeout timed out")
			tick = tick + 1
			if tick%5 == 0 {
				//g.image.Free() // is this needed?
				DrawImage(g, window)
			}
			continue
		}
		switch t := event.(type) {
		case *sdl.QuitEvent:
			running = false
		case *sdl.KeyboardEvent:
			keyCode := t.Keysym.Sym
			if keyCode == sdl.K_ESCAPE {
				running = false
				continue
			}
			if keyCode == sdl.K_SPACE {
				//g.image.Free() // is this needed?
				DrawImage(g, window)
			}
		}

	}

	return
}

func main() {
	g := NewGame()
	fmt.Println(g)
	if err := viewport.ViewportInitialize("Viewport", g.image); err != nil {
		os.Exit(1)
	}
	if err := run(g); err != nil {
		os.Exit(1)
	}
}
