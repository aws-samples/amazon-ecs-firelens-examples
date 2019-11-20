package main

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	for true {
		log.WithFields(log.Fields{
			"requestID": "45234523",
			"path":      "/",
		}).Info("Got a request")

		log.WithFields(log.Fields{
			"requestID": "546745643",
			"path":      "/tardis",
			"user":      "TheMaster",
		}).Warn("Access denied")

		log.WithFields(log.Fields{
			"requestID": "546745643",
			"path":      "/tardis",
			"user":      "TheDoctor",
		}).Debug("Admin access")
		time.Sleep(100 * time.Millisecond)
	}
}
