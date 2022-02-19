package imagefs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func CountGoFiles(folder string) (count int) {
	return CountGoFilesByExt(folder, ".go")
}

func CountGoFilesByExt(folder string, ext string) (count int) {
	files, err := os.ReadDir(folder)
	fmt.Println("Scanning  folder " + folder)
	if err != nil {
		return 0
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ext {
			count++
		}
	}
	return count
}

func CountAllGoFiles(folder string) (count int) {
	return CountAllFilesByExt(folder, ".go")
}

func CountAllJpgFiles(folder string) (count int) {
	return CountAllFilesByExt(folder, ".jpg")
}

func CountAllFilesByExt(folder string, ext string) (count int) {
	fsys := os.DirFS(folder)
	fmt.Println("Scanning  all " + ext + " in folder " + folder)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if filepath.Ext(p) == ext {
			count++
		}
		return nil
	})
	return count
}

func AllJpgFiles(folder string) (files []string) {
	return AllFilesByExt(folder, ".jpg")
}

func AllGoFiles(folder string) (files []string) {
	return AllFilesByExt(folder, ".go")
}

func AllFilesByExt(folder string, ext string) (files []string) {

	fsys := os.DirFS(folder)
	fmt.Println("Scanning  all in folder " + folder)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if filepath.Ext(p) == ext {
			files = append(files, p)
		}
		return nil
	})
	return files
}
