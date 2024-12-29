package storage

import (
	"cloud-commis/config"
	"cloud-commis/logger"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var dataDir = config.ParsedData.String("localStoragePath")

const (
	localFileName = "ccdata.json"
)

type localStorage struct {
}

func createDir() bool {

	logger.Log.Debug("data directory is " + dataDir)

	mkerr := os.MkdirAll(dataDir, 0750)
	if mkerr != nil {
		logger.Log.Error(mkerr.Error())
	}

	return true
}

func (localS localStorage) Read() (Aws_scans, bool) {
	var jsonData Aws_scans

	readFile, err := os.ReadFile(dataDir + localFileName)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	logger.Log.Debug("reading file " + dataDir + localFileName)
	err = json.Unmarshal(readFile, &jsonData.Data)
	if err != nil {
		logger.Log.Error("Can't read json file : " + err.Error())
	}
	fmt.Println(jsonData)

	return jsonData, true
}

func (localS localStorage) Write(data Aws_scans) bool {

	if !createDir() {
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

	return true
}

func (localS localStorage) Delete() bool {

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
