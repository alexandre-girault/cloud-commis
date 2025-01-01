package storage

import (
	"cloud-commis/config"
	"cloud-commis/logger"
)

type Storage interface {
	Read() (Aws_scans, error)
	Write(Aws_scans) error
	Delete() bool
}

var Data Storage

func Configure() {
	dataDir := config.ParsedData.String("localStoragePath")

	switch config.ParsedData.String("storage") {

	case "local":
		logger.Log.Info("Using local filesystem as storage")
		logger.Log.Info("Local file is " + dataDir + localFileName)
		Data = localStorage{}

	case "s3":
		logger.Log.Info("Using s3 bucket as storage")

	}

}
