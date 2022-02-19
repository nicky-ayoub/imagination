package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"

	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/nicky-ayoub/imagination/internal/pkg/imagefs"
)

func getName(g *Game) {
	randomIndex := rand.Intn(len(g.paths))
	g.name = g.paths[randomIndex]
}

func openImage(g *Game) (err error) {
	getName(g)
	g.img, _, err = ebitenutil.NewImageFromFile(g.root + "/" + g.name)

	return err
}

func NewGame() *Game {
	var err error
	g := &Game{}
	seed := time.Now().Unix()
	rand.Seed(seed)
	fmt.Println("Seed : ", seed)
	g.root = "../../assets"
	g.paths = imagefs.AllJpgFiles(g.root)

	err = openImage(g)
	if err != nil {
		fmt.Println("Fatal in init()")
		log.Fatal(err)
	}
	return g
}

type Game struct {
	img   *ebiten.Image
	name  string
	paths []string
	root  string
}

func TryNextImage(g *Game) (bool, bool) {

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		return true, false
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) || ebiten.IsKeyPressed(ebiten.KeyQ) {
		return true, true
	}
	return false, false
}

func (g *Game) Update() (err error) {

	go_to_next, quit := TryNextImage(g)
	if quit {
		log.Fatal("Exiting...")
		return err
	}
	if go_to_next {
		err = openImage(g) // Generate a new image
		if err != nil {
			fmt.Println("Fatal in init()")
			log.Fatal(err)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebiten.SetWindowSize(g.img.Bounds().Dx(), g.img.Bounds().Dy())
	ebiten.SetWindowTitle("Showing - " + g.name)
	screen.DrawImage(g.img, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.img.Bounds().Dx(), g.img.Bounds().Dy()
}

func main() {
	g := NewGame()
	ebiten.SetWindowSize(g.img.Bounds().Dx(), g.img.Bounds().Dy())
	ebiten.SetWindowTitle("Showing - " + g.name)
	fmt.Println(g.img.Bounds())
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
