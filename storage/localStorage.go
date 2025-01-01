package storage

import (
	"cloud-commis/config"
	"cloud-commis/logger"
	"encoding/json"
	"os"
	"strings"
)

const (
	localFileName = "ccdata.json"
)

type localStorage struct {
}

func createDir() error {
	dataDir := config.ParsedData.String("localStoragePath")

	logger.Log.Error("data directory is " + dataDir)

	mkerr := os.MkdirAll(dataDir, 0750)
	if mkerr != nil {
		logger.Log.Error(dataDir + mkerr.Error())
	}

	return mkerr
}

func (localS localStorage) Read() (Aws_scans, error) {
	var jsonData Aws_scans
	dataDir := config.ParsedData.String("localStoragePath")

	readFile, err := os.ReadFile(dataDir + localFileName)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	logger.Log.Debug("reading file " + dataDir + localFileName)
	err = json.Unmarshal(readFile, &jsonData)
	if err != nil {
		logger.Log.Error("Can't read json file : " + err.Error())
	}

	return jsonData, err
}

func (localS localStorage) Write(data Aws_scans) error {

	dataDir := config.ParsedData.String("localStoragePath")

	err := createDir()
	if err != nil {
		logger.Log.Error("Can't write in directory " + dataDir)
	}

	jsonString, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		logger.Log.Error(marshalErr.Error())
	}

	writeErr := os.WriteFile(dataDir+localFileName, jsonString, os.ModePerm)
	if writeErr != nil {
		logger.Log.Error(writeErr.Error())
	}

	return writeErr
}

func (localS localStorage) Delete() bool {

	dataDir := config.ParsedData.String("localStoragePath")

	files, err := os.ReadDir(dataDir)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".json") && !file.IsDir() {
			os.Remove(dataDir + file.Name())
			logger.Log.Debug("Deleting file " + file.Name())
		}
	}

	return true
}

//func GetFiles() []string {
//
//	var fileNames []string
//
//	files, err := os.ReadDir(dataDir)
//	if err != nil {
//		logger.Log.Error(err.Error())
//		panic(err)
//	}
//
//	for _, file := range files {
//		if strings.Contains(file.Name(), ".json") && !file.IsDir() {
//			fileNames = append(fileNames, file.Name())
//			logger.Log.Debug("Found file " + file.Name())
//		}
//	}
//	return fileNames
//}
