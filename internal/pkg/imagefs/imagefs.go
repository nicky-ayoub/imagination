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
	return CountAllFilesByExt(folder, ".jpg") + CountAllFilesByExt(folder, ".jpeg")
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
	return AllFilesByExt(folder, []string{".jpg", ".jpeg"})
}

func AllGoFiles(folder string) (files []string) {
	return AllFilesByExt(folder, []string{".go"})
}

func AllFilesByExt(folder string, exts []string) (files []string) {
	valid := make(map[string]bool)
	for _, s := range exts {
		valid[s] = true
	}
	fsys := os.DirFS(folder)
	fmt.Println("Scanning  all in folder " + folder)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if valid[filepath.Ext(p)] {
			files = append(files, p)
		}
		return nil
	})
	return files
}
