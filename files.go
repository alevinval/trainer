package trainer

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func getFileContent(name string) (b []byte, err error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func FindFiles(rootPath string) (fileNames chan string, err error) {
	allPaths, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return fileNames, err
	}

	fileNames = make(chan string)
	go func() {
		for _, filePath := range allPaths {
			fileNames <- path.Join(rootPath, filePath.Name())
		}
		close(fileNames)
	}()

	return fileNames, nil
}

func FindFilesWithPrefix(rootPath string, prefix string) (fileNames chan string, err error) {
	files, err := FindFiles(rootPath)
	if err != nil {
		return fileNames, err
	}

	prefix = path.Join(rootPath, prefix)
	fileNames = make(chan string)
	go func() {
		for name := range files {
			if len(name) < len(prefix) {
				continue
			}
			if strings.Compare(name[:len(prefix)], prefix) != 0 {
				continue
			}
			fileNames <- name
		}
		close(fileNames)
	}()

	return fileNames, nil
}
