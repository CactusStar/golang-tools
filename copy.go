package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	source := "C:\\Users\\XHe1\\Desktop\\TestCopy\\folder1"
	destination := "C:\\Users\\XHe1\\Desktop\\TestCopy\\folder2"

	// test1 := "C:\\Users\\XHe1\\Desktop\\TestCopy\\folder1\\abc_test"
	// test2 := "C:\\Users\\XHe1\\Desktop\\TestCopy\\folder1"
	// name := strings.Split(test1, source)[1]
	// fmt.Println(name)
	// CopyPath(source, destination)
	GetCountAndCompare(source, destination)
	// fmt.Scanf("h")
}

func Listfiles(src string) map[string]int {
	filecount := make(map[string]int)
	files, _ := ioutil.ReadDir(src)
	for _, file := range files {
		if file.IsDir() {
			// fmt.Printf("this is directory: %s\n", file.Name())
			path := src + "\\" + file.Name()
			filecount = Listfiles(path)
		}
	}
	filecount[src] = len(files)
	return filecount
}

func CompareCount(before map[string]int, after map[string]int, src string) string {
	var name string
	for beforename := range before {
		for aftername := range after {
			if beforename == aftername {
				if before[beforename] != after[aftername] {
					fmt.Printf("File count changed for %s\n", beforename)
					name = strings.Split(beforename, src)[1]
					fmt.Println(name)
					return name
				}
			}
		}
	}

	return "nochange"
}

func GetCountAndCompare(src string, dst string) {
	before := Listfiles(src)
	for {
		time.Sleep(time.Duration(5) * time.Second)
		after := Listfiles(src)
		changed := CompareCount(before, after, src)
		
		if changed != "nochange" {
			fmt.Println("start copy")
			srcchanged := src + changed
			dstchanged := dst + changed
			fmt.Printf("source changed folder: %s\n", srcchanged)
			fmt.Printf("destination changed folder: %s\n", dstchanged)
			CopyFileonly(srcchanged, dstchanged)
		} else {
			fmt.Println("no change")
		}
		
		before = after
	}
}

func GetFileInfo(src string) os.FileInfo {
	if fileInfo, e := os.Stat(src); e != nil {
		if os.IsNotExist(e) {
			return nil
		}
		return nil
	} else {
		return fileInfo
	}
}

func CopyFile(src, dst string) bool {
	if len(src) == 0 || len(dst) == 0 {
		return false
	}
	srcFile, e := os.OpenFile(src, os.O_RDONLY, os.ModePerm)
	if e != nil {
		fmt.Printf("copyfile1 %s\n", e)
		return false
	}
	defer srcFile.Close()

	dst = strings.Replace(dst, "\\", "/", -1)
	dstPathArr := strings.Split(dst, "/")
	dstPathArr = dstPathArr[0 : len(dstPathArr)-1]
	dstPath := strings.Join(dstPathArr, "/")

	dstFileInfo := GetFileInfo(dstPath)

	if dstFileInfo == nil {
		if e := os.MkdirAll(dstPath, os.ModePerm); e != nil {
			fmt.Printf("copyfile2 %s\n", e)
			return false
		}
	}
	//这里要把O_TRUNC 加上，否则会出现新旧文件内容出现重叠现象
	fmt.Printf("======================destination location is=================> %s\n", dst)
	dstFile, e := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if e != nil {
		fmt.Printf("copyfile3 %s\n", e)
		return false
	}
	defer dstFile.Close()
	//fileInfo, e := srcFile.Stat()
	//fileInfo.Size() > 1024
	//byteBuffer := make([]byte, 10)
	if _, e := io.Copy(dstFile, srcFile); e != nil {
		fmt.Printf("copyfile4 %s\n", e)
		return false
	} else {
		return true
	}

}

func CopyPath(src, dst string) bool {
	srcFileInfo := GetFileInfo(src)
	if srcFileInfo == nil || !srcFileInfo.IsDir() {
		fmt.Printf("not dir\n")
		return false
	}
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("copypath1 %s\n", err)
			return err
		}
		relationPath := strings.Replace(path, src, "/", -1)
		dstPath := strings.TrimRight(strings.TrimRight(dst, "/"), "\\") + relationPath
		if !info.IsDir() {
			if CopyFile(path, dstPath) {
				return nil
			} else {
				return errors.New(path + " copy fail")
			}
		} else {
			if _, err := os.Stat(dstPath); err != nil {
				if os.IsNotExist(err) {
					if err := os.MkdirAll(dstPath, os.ModePerm); err != nil {
						fmt.Printf("copypath2 %s\n", err)
						return err
					} else {
						return nil
					}
				} else {
					fmt.Printf("copypath3 %s\n", err)
					return err
				}
			} else {
				return nil
			}
		}
	})

	if err != nil {
		fmt.Printf("copypath4 %s\n", err)
		return false
	}
	return true

}

func CopyFileonly(src, dst string) {
	files, _ := ioutil.ReadDir(src)
	for _, file := range files {
		if !file.IsDir() {
			sourcefile := src + "\\" + file.Name()
			fmt.Printf("This is file to copy ===> %s\n", sourcefile)
			destinationfile := dst + "\\" + file.Name()
			CopyFile(sourcefile, destinationfile)
		}
	}
}

func GetChangedName(src string, changedfolder string) string {
	sourceChanged := strings.Split(src, changedfolder)[0]
	if sourceChanged == "" {
		sourceChanged = src
	}
	return sourceChanged
}