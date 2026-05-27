package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Recursively scans the target directory and returns macOS metadata files to clean.
func readDirRecursion(dirName, trashDir string) ([]string, error) {
	var files []string
	fs, err := os.ReadDir(dirName)
	if err != nil {
		return nil, err
	}
	for _, f := range fs {
		fullPath := filepath.Join(dirName, f.Name())
		if f.IsDir() {
			if fullPath == trashDir {
				continue
			}
			childrenFiles, err := readDirRecursion(fullPath, trashDir)
			if err != nil {
				return nil, err
			}
			files = append(files, childrenFiles...)
		} else {
			ok, err := isMacJunkFile(f)
			if err != nil {
				return nil, err
			}
			if ok {
				files = append(files, fullPath)
			}
		}
	}
	return files, nil
}

func isMacJunkFile(entry os.DirEntry) (bool, error) {
	if entry.Name() == ".DS_Store" {
		return true, nil
	}
	if !strings.HasPrefix(entry.Name(), "._") {
		return false, nil
	}
	info, err := entry.Info()
	if err != nil {
		return false, err
	}
	return info.Size() == 4096, nil
}

// Returns a destination path that does not conflict with existing files.
func uniqueDest(dest string) string {
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		return dest
	}
	ext := filepath.Ext(dest)
	name := strings.TrimSuffix(filepath.Base(dest), ext)
	dir := filepath.Dir(dest)
	i := 1
	for {
		newDest := filepath.Join(dir, fmt.Sprintf("%s_%d%s", name, i, ext))
		if _, err := os.Stat(newDest); os.IsNotExist(err) {
			return newDest
		}
		i++
	}
}

// Moves files to the archive directory.
func moveFiles(files []string, trashDir string) error {
	for _, f := range files {
		dest := uniqueDest(filepath.Join(trashDir, filepath.Base(f)))
		if err := os.Rename(f, dest); err != nil {
			return err
		}
		fmt.Println("Moved to archive:", f)
	}
	return nil
}

func main() {
	var origDir string
	flag.StringVar(&origDir, "t", "", "Clean the target directory, default is current directory")
	flag.Parse()

	if origDir == "" {
		var err error
		origDir, err = os.Getwd()
		if err != nil {
			fmt.Println("Failed to get current working directory:", err)
			return
		}
	}
	trashDir := filepath.Join(origDir, ".wait_clean")
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		fmt.Println("Failed to create archive directory:", err)
		return
	}

	files, err := readDirRecursion(origDir, trashDir)
	if err != nil {
		fmt.Println("Failed to read target directory:", err)
		return
	}
	if len(files) == 0 {
		fmt.Println("No macOS metadata files found")
		if err := os.Remove(trashDir); err != nil && !os.IsNotExist(err) {
			fmt.Println("Failed to remove empty archive directory:", err)
		}
		return
	}

	if err := moveFiles(files, trashDir); err != nil {
		fmt.Println("Failed to move files to archive:", err)
		return
	}

	fmt.Print("Delete archived files? (Y/n): ")
	var input string
	fmt.Scanln(&input)
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "y" || input == "" {
		if err := os.RemoveAll(trashDir); err != nil {
			fmt.Println("Failed to delete archived files:", err)
		} else {
			fmt.Println("Archived files deleted")
		}
	} else {
		fmt.Println("Archived files kept at:", trashDir)
	}
}
