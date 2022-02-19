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

func run(g *Game) (err error) {

	var Frame_Starting_Time = uint32(0)
	var Elapsed_Time uint32
	var Mouse_X int32
	var Mouse_Y int32
	var Zoom_Factor = int32(1)
	Flipping_Mode := viewport.TViewportFlippingModeID(viewport.FLIPPING_MODE_ID_NORMAL)

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
	g.name = g.paths[1]
	g.image, err = img.Load(g.root + "/" + g.name)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Initialize modules (no need to display an error message if a module initialization fails because the module already did)
	if err = viewport.Initialize(g.title, g.image); err != nil {
		log.Fatal(err)
		return err
	} // TODO set initial viewport size and window decorations according to parameters saved on previous program exit ?

	viewport.DrawImage()
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
				// Tell the viewport that its size

				if t.Event == sdl.WINDOWEVENT_SIZE_CHANGED {
					fmt.Printf("Window size change to (%d, %d) %d %t\n", t.Data1, t.Data2, t.Event, t.Event == sdl.WINDOWEVENT_SIZE_CHANGED)
					viewport.SetDimensions(t.Data1, t.Data2)
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
				viewport.SetZoomedArea(Mouse_X, Mouse_Y, Zoom_Factor)

			case *sdl.KeyboardEvent:
				if t.Type == sdl.KEYDOWN {
					// Toggle image flipping
					switch t.Keysym.Sym {
					case sdl.K_f:
						// Set next available flipping mode
						Flipping_Mode++
						if Flipping_Mode >= viewport.FLIPPING_MODE_IDS_COUNT {
							Flipping_Mode = viewport.FLIPPING_MODE_ID_NORMAL
						}
						viewport.SetFlippingMode(Flipping_Mode)

						// Zoom has been reset when flipping the image
						Zoom_Factor = 1
					case sdl.K_q:
						fmt.Println("Application quit...")
						running = false
					case sdl.K_s:
						// Scale image to fit viewport
						viewport.ScaleImage()
						// Reset zoom
						Zoom_Factor = 1
					}
				}
			case *sdl.MouseMotionEvent:
				if t.Type == sdl.MOUSEMOTION {
					// Do not recompute everything when the image is not zoomed
					if Zoom_Factor > 1 {
						// Successively zoom to the current zoom level to make sure the internal ViewportSetZoomedArea() data are consistent
						i := int32(1)
						for i <= Zoom_Factor {
							viewport.SetZoomedArea(t.X, t.Y, i)
							i <<= 1
						}
					}
				}
			default:
				//fmt.Printf("[%d ms] Unknown\ttype:%d\n", t.GetTimestamp(), t.GetType())
			}
		}
		//fmt.Println("Drawing image")
		viewport.DrawImage()

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

}
