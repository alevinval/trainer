package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/alevinval/trainer"
)

type inputArgs struct {
	cmd          string
	lookupPath   string
	searchPrefix string
}

func parseArgs() (args *inputArgs, err error) {
	flag.Usage = func() {
		fmt.Print("./trainer [command] [path] [prefix]\n" +
			"command: cluster|performance.\n" +
			"path: where to look for .gpx files.\n" +
			"prefix: prefix that must be satisfied by the *.gpx files.\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if len(flag.Args()) < 2 {
		flag.Usage()
		return nil, errors.New("invalid arguments")
	}
	args = &inputArgs{
		cmd:          flag.Arg(0),
		lookupPath:   flag.Arg(1),
		searchPrefix: flag.Arg(2),
	}
	return args, nil
}

func findActivities(lookupPath, prefix string) trainer.ActivityList {
	activities := trainer.ActivityList{}

	fileNames, err := findFilesWithPrefix(lookupPath, prefix)
	if err != nil {
		log.Printf("cannot find activities in %s: %s\n", lookupPath, err)
		return activities
	}

	wg := new(sync.WaitGroup)
	mux := new(sync.Mutex)
	for fileName := range fileNames {
		wg.Add(1)
		go func(fileName string) {
			defer wg.Done()
			activity, err := trainer.OpenFile(fileName)
			if err != nil {
				log.Printf("cannot open file %q: %s\n", fileName, err)
				return
			}
			mux.Lock()
			activities = append(activities, activity)
			mux.Unlock()
		}(fileName)
	}
	wg.Wait()
	return activities
}

func cluster(activities trainer.ActivityList) {
	for _, cluster := range activities.GetClusters() {
		avgPerf := cluster.Activities.DataPoints().AvgPerf()
		fmt.Printf("%s\nAvg.perf: %0.2f\n\n", cluster, avgPerf)
	}
	return
}

func performance(activities trainer.ActivityList, print bool) {
	histogram := activities.DataPoints().GetHistogram()
	trainer.PrintHistogram(histogram)
	output, _ := os.Create("global.csv")
	defer output.Close()
	trainer.WriteCsvTo(histogram, output)
}

func main() {
	args, err := parseArgs()
	if err != nil {
		return
	}
	activities := findActivities(args.lookupPath, args.searchPrefix)
	switch args.cmd {
	case "cluster":
		cluster(activities)
	case "performance":
		performance(activities, true)
	}
}
