package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nicky-ayoub/imagination/internal/pkg/imagefs"
)

func main() {
	cmd := os.Args[0]
	fmt.Println("Imagination Suite")
	fmt.Println(cmd)

	extPtr := flag.String("ext", ".go", "a file extention")
	flag.Parse()
	fmt.Println("Extension:", *extPtr)

	dir := "."
	if len(flag.Args()) > 0 {
		dir = flag.Args()[0]
	}
	fmt.Println("Scanning", dir)

	fmt.Println(imagefs.CountAllGoFiles("/home/nicky/go"))
	fmt.Println(imagefs.CountAllGoFiles(dir))
	fmt.Println(imagefs.CountAllFilesByExt(dir, *extPtr))

	for _, file := range imagefs.AllFilesByExt(dir, *extPtr) {
		fmt.Println(file)
	}
	fmt.Println(imagefs.CountAllJpgFiles(dir))
	// aPath := dir + "/**/*" + *extPtr

	// fmt.Println("Print globbing", aPath)
	// files, err := filepath.Glob(aPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, file := range files {

	// 	fmt.Println(file)
	// }
}
