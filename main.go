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
	logger.Log.Info("test : " + config.ParsedData.String("scan_interval_min"))
	storage.Configure()

	go awsScanner.ScheduledScan()
	router := http.NewServeMux()

	webui.Start(router)
	log.Fatal(http.ListenAndServe(":"+config.ParsedData.String("httpPort"), router))
}
