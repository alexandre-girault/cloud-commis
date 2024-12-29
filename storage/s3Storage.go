package storage

import (
	"cloud-commis/logger"
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type CachedData struct {
	Data      Aws_scans
	CacheDate time.Time
}

var Cache = make(map[string]CachedData)

func clearCache() {
	for {
		logger.Log.Debug("Checking cache")
		for cacheKey, cacheValue := range Cache {
			if time.Since(cacheValue.CacheDate).Seconds() > 20 {
				logger.Log.Debug("Cached data for " + cacheKey + " too old, deleting")
				delete(Cache, cacheKey)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	logger.Log.Info(strconv.Itoa(int(elapsed.Milliseconds())) + "ms")
}

func init() {
	go clearCache()
}

func GetS3File(bucketName string, filePath string) Aws_scans {

	defer timeTrack(time.Now())

	//Check in Cache
	val, isInCache := Cache[bucketName+filePath]
	if isInCache {
		logger.Log.Debug("Data for " + bucketName + filePath + " found in Cache")
		return val.Data
	} else {
		logger.Log.Debug("Data for " + bucketName + filePath + " not found in Cache")
	}

	aws_session, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("eu-west-3"),
		},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		logger.Log.Error("Unable to create a new AWS session " + err.Error())
	}

	s3Session := s3.New(aws_session)
	remoteFile := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filePath),
	}

	result, err := s3Session.GetObject(remoteFile)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	var jsonData Aws_scans
	merr := json.Unmarshal(body, &jsonData.Data)
	if merr != nil {
		logger.Log.Error(merr.Error())
	}

	var newBufferData CachedData
	newBufferData.CacheDate = time.Now()
	newBufferData.Data = jsonData

	Cache[bucketName+filePath] = newBufferData

	return jsonData

}
