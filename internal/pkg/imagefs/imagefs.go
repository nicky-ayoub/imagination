package imagefs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func CountGoFiles(folder string, count int) int {
	files, err := os.ReadDir(folder)
	fmt.Println("Scanning  folder " + folder)
	if err != nil {
		return 0
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".go" {
			count++
		}
	}
	return count
}

func CountAllGoFiles(folder string) (count int) {
	return CountAllFilesByExt(folder, ".go")
}

func CountAllFilesByExt(folder string, ext string) (count int) {
	fsys := os.DirFS(folder)
	fmt.Println("Scanning  all in folder " + folder)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if filepath.Ext(p) == ext {
			count++
		}
		return nil
	})
	return count
}
