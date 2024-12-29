package main

import (
	"cloud-commis/awsScanner"
	"cloud-commis/config"
	"cloud-commis/logger"
	"cloud-commis/webui"
	"log"
	"net/http"
)

func main() {

	config.Read()
	// Start
	logger.Log.Info("loglevel = " + config.ParsedData.String("loglevel"))
	logger.SetLogLevel(config.ParsedData.String("loglevel"))

	logger.Log.Info("test : " + config.ParsedData.String("scan_interval_min"))

	go awsScanner.ScheduledScan()
	router := http.NewServeMux()

	webui.Start(router)
	log.Fatal(http.ListenAndServe(":"+config.ParsedData.String("httpPort"), router))
}
