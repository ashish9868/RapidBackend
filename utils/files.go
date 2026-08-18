package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func FileExists(embed fs.FS, path string) bool {
	_, err := fs.Stat(embed, path)
	if err == nil {
		return true
	}
	return false
}

func PrintFiles(embed fs.FS) error {
	// Walk the root directory "." to visit every file and folder
	err := fs.WalkDir(embed, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip printing the root folder itself
		if path == "." {
			return nil
		}

		// Distinguish between files and directories
		if d.IsDir() {
			fmt.Printf("[DIR]  %s\n", path)
		} else {
			fmt.Printf("[FILE] %s\n", path)
		}
		return nil
	})

	return err
}

func SubFs(embed fs.FS, path string) *fs.FS {
	subFs, err := fs.Sub(embed, path)
	if err == nil {
		return &subFs
	}
	println("Error reading sub fs :", err.Error())
	return nil
}

func GetCWD() string {
	// Gets the absolute path of the running executable
	ex, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	return ex
}

func GetPathFromRoot(path string) string {
	return filepath.Join(GetCWD(), path)
}

func SafeCreateFile(path string, content string) error {

	_, err := os.Stat(path)

	if err == nil {
		return nil // File exists
	}
	if errors.Is(err, fs.ErrNotExist) {
		pathDir := filepath.Dir(path)
		err = SafeCreateFolder(pathDir)
		if err == nil {
			err := os.WriteFile(path, []byte(content), 0644)
			if err == nil {
				return nil
			}
			return err
		}
		return err // File explicitly does not exist
	}

	return nil
}

func SafeCreateFolder(folder string) error {
	root := GetCWD()
	folders := append([]string{root}, strings.Split(folder, "/")...)
	fullPath := path.Join(folders...)

	fmt.Println("Creating path", fullPath)
	stat, err := os.Stat(fullPath)

	createFolder := true
	if err == nil {
		if stat.IsDir() {
			createFolder = false
		}
	}

	if createFolder {
		err := os.MkdirAll(fullPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
