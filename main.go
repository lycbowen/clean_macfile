package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	origDir   string
	trash_dir string
)

// 递归扫描目录，返回所有需要归档的文件路径
func readDirRecursion(dirName string) ([]string, error) {
	var files []string
	fs, err := os.ReadDir(dirName)
	if err != nil {
		return nil, err
	}
	for _, f := range fs {
		fullPath := path.Join(dirName, f.Name())
		if f.IsDir() {
			childrenFiles, err := readDirRecursion(fullPath)
			if err != nil {
				return nil, err
			}
			files = append(files, childrenFiles...)
		} else {
			info, _ := f.Info()
			if (strings.HasPrefix(f.Name(), "._") && info.Size() == 4096) || f.Name() == ".DS_Store" {
				files = append(files, fullPath)
			}
		}
	}
	return files, nil
}

// 自动生成不冲突的目标文件名
func uniqueDest(dest string) string {
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		return dest
	}
	ext := filepath.Ext(dest)
	name := strings.TrimSuffix(filepath.Base(dest), ext)
	dir := filepath.Dir(dest)
	i := 1
	for {
		newDest := path.Join(dir, fmt.Sprintf("%s_%d%s", name, i, ext))
		if _, err := os.Stat(newDest); os.IsNotExist(err) {
			return newDest
		}
		i++
	}
}

// 归档文件
func moveFiles(files []string) error {
	for _, f := range files {
		dest := uniqueDest(path.Join(trash_dir, path.Base(f)))
		if err := os.Rename(f, dest); err != nil {
			return err
		}
		fmt.Println("已归档", f)
	}
	return nil
}

func main() {
	flag.StringVar(&origDir, "t", "", "Clean the target directory, default is current directory")
	flag.Parse()

	if origDir == "" {
		origDir, _ = os.Getwd()
	}
	trash_dir = path.Join(origDir, ".wait_clean")
	os.MkdirAll(trash_dir, 0755)

	files, err := readDirRecursion(origDir)
	if err != nil {
		fmt.Println("读取目录失败:", err)
		return
	}
	if len(files) == 0 {
		fmt.Println("没有找到可归档的文件")
		os.Remove(trash_dir)
		return
	}

	if err := moveFiles(files); err != nil {
		fmt.Println("归档文件失败:", err)
		return
	}

	fmt.Print("是否要删除已归档文件？(Y/n)：")
	var input string
	fmt.Scanln(&input)
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "y" || input == "" {
		if err := os.RemoveAll(trash_dir); err != nil {
			fmt.Println("删除失败:", err)
		} else {
			fmt.Println("已删除归档文件")
		}
	} else {
		fmt.Println("归档文件保存在:", trash_dir)
	}
}
