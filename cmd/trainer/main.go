package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

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

func findActivities(lookupPath, prefix string) <-chan *trainer.Activity {
	ch := make(chan *trainer.Activity)
	go func() {
		fileNames, err := findFilesWithPrefix(lookupPath, prefix)
		if err != nil {
			log.Printf("cannot find activities in %s: %s\n", lookupPath, err)
			close(ch)
		}
		for fileName := range fileNames {
			activity, err := trainer.OpenFile(fileName)
			if err != nil {
				log.Printf("cannot open file %q: %s\n", fileName, err)
				continue
			}
			ch <- activity
		}
		close(ch)
	}()
	return ch
}

func cluster(activities trainer.ActivityList) {
	for _, cluster := range activities.GetClusters() {
		hist := trainer.Histogram{}
		hist.Reset()
		hist.Feed(cluster.Activities)
		flat := hist.Flatten()
		fmt.Printf("%s\nAvg.perf: %0.2f\n\n", cluster, flat.GetAvgPerf())
	}
	return
}

func performance(activities trainer.ActivityList, print bool) {
	globalHist := &trainer.Histogram{}
	globalHist.Reset()
	for _, activity := range activities {
		fmt.Printf("Processing %q\n", activity.Metadata().Name)
		globalHist.Feed(activity)
	}
	globalHist.PrintRaw()
	output, _ := os.Create("global.csv")
	defer output.Close()
	globalHist.WriteTo(output)
}

func main() {
	args, err := parseArgs()
	if err != nil {
		return
	}
	activities := trainer.ActivityList{}
	for activity := range findActivities(args.lookupPath, args.searchPrefix) {
		activities = append(activities, activity)
	}
	switch args.cmd {
	case "cluster":
		cluster(activities)
	case "performance":
		performance(activities, true)
	}
}
