package main

import (
	"log"

	"fmt"
	_ "image/jpeg"
	_ "image/png"

	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/nicky-ayoub/imagination/internal/pkg/imagefs"
)

var img *ebiten.Image
var name string
var paths []string

func getName() (pick string) {

	rand.Seed(time.Now().Unix())

	randomIndex := rand.Intn(len(paths))
	pick = paths[randomIndex]
	return pick
}

func openImage() (image *ebiten.Image, err error) {
	name = getName()
	image, _, err = ebitenutil.NewImageFromFile(name)

	if err != nil {
		return nil, err
	}
	return image, nil
}

func init() {
	var err error

	fmt.Println("Scanning cwd")
	paths = imagefs.AllJpgFiles(".")

	img, err = openImage()
	if err != nil {
		fmt.Println("Fatal in init()")
		log.Fatal(err)
	}
}

type Game struct{}

func (g *Game) Update() (err error) {
	img, err = openImage() // Generate a new image
	if err != nil {
		fmt.Println("Fatal in init()")
		log.Fatal(err)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebiten.SetWindowSize(img.Bounds().Dx(), img.Bounds().Dy())
	ebiten.SetWindowTitle("Render an image - " + name)
	screen.DrawImage(img, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return img.Bounds().Dx(), img.Bounds().Dy()
}

func main() {
	ebiten.SetWindowSize(img.Bounds().Dx(), img.Bounds().Dy())
	ebiten.SetWindowTitle("Render an image - " + name)
	fmt.Println(img.Bounds())
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
