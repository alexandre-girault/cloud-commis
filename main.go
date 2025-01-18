package main

import (
	"cloud-commis/awsScanner"
	"cloud-commis/config"
	"cloud-commis/logger"
	"cloud-commis/storage"
	"cloud-commis/webui"
	"log"
	"net/http"
)

func main() {

	config.Read(config.ParsedData)
	// Start
	logger.Log.Info("loglevel = " + config.ParsedData.String("loglevel"))
	logger.SetLogLevel(config.ParsedData.String("loglevel"))
	logger.Log.Info("scan_interval_min : " + config.ParsedData.String("scanIntervalMin"))

	if !storage.Configure() {
		logger.Log.Info("No data found, triggering a first scan now")
		awsScanner.Aws_instances_inventory()
	}

	go awsScanner.ScheduledScan()
	router := http.NewServeMux()

	webui.Start(router)
	log.Fatal(http.ListenAndServe(":"+config.ParsedData.String("httpPort"), router))
}
