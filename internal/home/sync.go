package home

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/alevinval/trainer/internal/utils"
)

// Sync checks lookupPath against the home folder and copies
// missing activities.
func Sync(lookupPath string) {
	pathsToSync, err := getMissingPaths(lookupPath, ActivitiesPath)
	if err != nil {
		log.Printf("error synchronising: %s\n", err)
		return
	}
	if len(pathsToSync) == 0 {
		log.Printf("no activities pending to be synchronised\n")
		return
	}

	activities := utils.ActivitiesFromPaths(pathsToSync)
	for _, activity := range activities {
		fname := filepath.Base(activity.Metadata().DataSource.Name)
		dstPath := filepath.Join(ActivitiesPath, fname)
		log.Printf("synchronising activity: %s\n", fname)

		f, err := os.Create(dstPath)
		if err != nil {
			log.Printf("[error] cannot create path: %s\n", err)
			continue
		}
		var dstWriter io.WriteCloser = f
		defer dstWriter.Close()

		if filepath.Ext(fname) == ".gz" {
			dstWriter = gzip.NewWriter(dstWriter)
		}

		_, err = dstWriter.Write(activity.Bytes())
		if err != nil {
			log.Printf("[error] cannot synchronise file: %s\n", err)
			continue
		}
	}
}

// Given a source and destination path, finds paths present
// in source folder and missing at destination folder.
func getMissingPaths(srcDir, dstDir string) (paths []string, err error) {
	srcPaths, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return nil, err
	}

	dstPaths, err := ioutil.ReadDir(dstDir)
	if err != nil {
		return nil, err
	}

	dstMap := map[string]struct{}{}
	for i := range dstPaths {
		dstMap[dstPaths[i].Name()] = struct{}{}
	}

	paths = []string{}
	for _, srcPath := range srcPaths {
		_, present := dstMap[srcPath.Name()]
		if !present {
			fullPath := filepath.Join(srcDir, srcPath.Name())
			paths = append(paths, fullPath)
		}
	}

	return paths, nil
}
