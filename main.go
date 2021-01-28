package main

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"runsync/API"
	"runsync/API/nike"
	"runsync/API/strava"
)

func main() {
	// Initialize Logger
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		PadLevelText:  true,
		FullTimestamp: true,
	})

	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
		log.Exit(1)
	}

	runs, err := nike.GetRunsFromNRC()
	if err != nil {
		log.Error("Error will loading data from Nike Run Club")
		log.Exit(1)
	}

	log.WithFields(
		log.Fields{
			"length": len(runs),
		},
	).Info("Runs retrieved from Nike Run Club")

	paths := []string{}

	for _, run := range runs {
		if API.Contains(run.MetricTypes, "longitude") && API.Contains(run.MetricTypes, "longitude") {
			gpx := nike.BuildGpxFromActivity(run)
			if nil != gpx {
				paths = append(paths, API.WriteGpxToFile(run.ID, gpx))
			}
		} else {
			tcx := nike.BuildTcxFromActivity(run)
			if nil != tcx {
				paths = append(paths, API.WriteTcxToFile(run.ID, tcx))
			}
		}
	}

	for _, path := range paths {
		strava.ImportDataFromFiles(path)
	}

}
