package main

import (
	"io/ioutil"
	"path"
	"strings"
)

func findFiles(root string) (filePaths chan string, err error) {
	paths, err := ioutil.ReadDir(root)
	if err != nil {
		return
	}
	filePaths = make(chan string)
	go func() {
		for _, filePath := range paths {
			filePaths <- path.Join(root, filePath.Name())
		}
		close(filePaths)
	}()
	return
}

func findFilesWithPrefix(root string, prefix string) (filePaths chan string, err error) {
	paths, err := findFiles(root)
	if err != nil {
		return
	}
	prefix = path.Join(root, prefix)
	filePaths = make(chan string)
	go func() {
		for filePath := range paths {
			if len(filePath) < len(prefix) {
				continue
			}
			if strings.Compare(filePath[:len(prefix)], prefix) != 0 {
				continue
			}
			filePaths <- filePath
		}
		close(filePaths)
	}()
	return
}
