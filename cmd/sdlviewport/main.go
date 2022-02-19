package main

import (
	"fmt"
	"log"
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
	title string
	seed  int64
}

func NewGame() *Game {
	g := &Game{}
	g.title = "SDL Image Viewer - "
	g.seed = time.Now().Unix()
	rand.Seed(g.seed)
	g.root = "../../assets"
	g.paths = imagefs.AllJpgFiles(g.root)
	return g
}

func GetName(g *Game) {
	randomIndex := rand.Intn(len(g.paths))
	g.name = g.paths[randomIndex]
}

func run(g *Game) (err error) {

	var Frame_Starting_Time = uint32(0)
	var Elapsed_Time uint32
	var Mouse_X int32
	var Mouse_Y int32
	var Zoom_Factor = int32(1)
	Flipping_Mode := viewport.TViewportFlippingModeID(viewport.VIEWPORT_FLIPPING_MODE_ID_NORMAL)

	g.name = g.paths[0]

	// Initialize SDL before everything else, so other SDL libraries can be safely initialized
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		os.Exit(1)
	}
	defer sdl.Quit()

	// Try to initialize the SDL image library
	err = img.Init(img.INIT_JPG | img.INIT_PNG | img.INIT_TIF)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer img.Quit()

	// Try to load the image before creating the viewport
	g.image, err = img.Load(g.root + "/" + g.name)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Initialize modules (no need to display an error message if a module initialization fails because the module already did)
	if err = viewport.ViewportInitialize(g.title, g.image); err != nil {
		return err
	} // TODO set initial viewport size and window decorations according to parameters saved on previous program exit ?
	g.image.Free()

	// Process incoming SDL events
	running := true
	for running {
		// Keep the time corresponding to the frame rendering beginning
		Frame_Starting_Time = sdl.GetTicks()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Application closed...")
				running = false

			case *sdl.WindowEvent:
				// Tell the viewport that its size changed
				if event.GetType() == sdl.WINDOWEVENT_SIZE_CHANGED {
					viewport.ViewportSetDimensions(t.Data1, t.Data2)
					Zoom_Factor = 1 // Zoom has been reset when resizing the window
				}
			case *sdl.MouseWheelEvent:
				// Wheel is rotated toward the user, increment the zoom factor
				if t.Y > 0 {
					if Zoom_Factor < viewport.CONFIGURATION_VIEWPORT_MAXIMUM_ZOOM_FACTOR {
						Zoom_Factor *= 2
					}
				} else { // Wheel is rotated away from the user, decrement the zoom factor
					if Zoom_Factor > 1 {
						Zoom_Factor /= 2
					}
				}
				// Start zooming area from the mouse coordinates
				Mouse_X, Mouse_Y, _ = sdl.GetMouseState()
				viewport.ViewportSetZoomedArea(Mouse_X, Mouse_Y, Zoom_Factor)

			case *sdl.KeyboardEvent:
				if t.Type == sdl.KEYDOWN {
					// Toggle image flipping
					switch t.Keysym.Sym {
					case sdl.K_f:
						// Set next available flipping mode
						Flipping_Mode++
						if Flipping_Mode >= viewport.VIEWPORT_FLIPPING_MODE_IDS_COUNT {
							Flipping_Mode = 0
						}
						viewport.ViewportSetFlippingMode(Flipping_Mode)

						// Zoom has been reset when flipping the image
						Zoom_Factor = 1
					case sdl.K_q:
						fmt.Println("Application quit...")
						running = false
					case sdl.K_s:
						// Scale image to fit viewport
						viewport.ViewportScaleImage()
						// Reset zoom
						Zoom_Factor = 1
					}
				}
			case *sdl.MouseMotionEvent:
				// Do not recompute everything when the image is not zoomed
				if Zoom_Factor > 1 {
					// Successively zoom to the current zoom level to make sure the internal ViewportSetZoomedArea() data are consistent
					i := int32(1)
					for i <= Zoom_Factor {
						viewport.ViewportSetZoomedArea(t.X, t.Y, i)
						i <<= 1
					}
				}
			default:
				fmt.Printf("[%d ms] Unknown\ttype:%d\n",
					t.GetTimestamp(), t.GetType())
			}
		}
		//fmt.Println("Drawing image")
		viewport.ViewportDrawImage()

		// Wait enough time to get a 60Hz refresh rate
		Elapsed_Time = sdl.GetTicks() - Frame_Starting_Time
		if Elapsed_Time < viewport.CONFIGURATION_DISPLAY_REFRESH_RATE_PERIOD {
			sdl.Delay(viewport.CONFIGURATION_DISPLAY_REFRESH_RATE_PERIOD - Elapsed_Time)
		}
	}
	return
}

func main() {
	g := NewGame()
	fmt.Println(g)
	if err := run(g); err != nil {
		os.Exit(1)
	}
	return
}
