package main

import (
	"fmt"
	"os"

	"github.com/nicky-ayoub/imagination/internal/pkg/imagefs"
)

func main() {
	cmd := os.Args[0]
	fmt.Println("Imagination Suite")
	fmt.Println(cmd)

	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	fmt.Println("Scanning", dir)

	fmt.Println(imagefs.CountAllGoFiles("/home/nicky/go"))
	fmt.Println(imagefs.CountAllGoFiles(dir))
	fmt.Println(imagefs.CountAllFilesByExt(dir, ".keep"))
}
