package main

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"runsync/API/nike"
)

func main() {
	// Initialize Logger
	log.SetFormatter(&log.TextFormatter{
		ForceColors:  true,
		PadLevelText: true,
		FullTimestamp: true,

	})

	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
		log.Exit(1)
	}

	activities, err := nike.GetActivitiesFromNRC()
	if err != nil {
		log.Error("Error will loading data from Nike Run Club")
		log.Exit(1)
	}

	log.WithFields(
		log.Fields{
			"activitiesLength": len(activities),
		},
	).Info("Activities retrieved from Nike Run Club")

}
