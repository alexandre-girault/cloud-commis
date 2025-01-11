package config

import (
	"cloud-commis/logger"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const (
	defaultConfigPath = "/etc/cloud-commis/config.yaml"
)

var ParsedData = koanf.New(".")
var AwsProfiles []awsProfile

type awsProfile struct {
	Name         string
	AwsAccountID int
	RoleArn      string
}

func Read(config *koanf.Koanf) {

	Flags.Read()

	// default config
	err := config.Load(confmap.Provider(map[string]interface{}{
		"configPath":       defaultConfigPath,
		"disable_ui":       false,
		"httpPort":         8080,
		"localStoragePath": "/data/cloud-commis",
		"loglevel":         "info",
		"s3BucketName":     "",
		"s3BucketPath":     "",
		"scanAws":          true,
		"scanIntervalMin":  60,
		"storage":          "local",
	}, "."), nil)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	// default config is overwritten by yaml config
	if _, err := os.Stat(Flags.configFile); err == nil {
		if err := config.Load(file.Provider(Flags.configFile), yaml.Parser()); err != nil {
			logger.Log.Error(err.Error())
		} else {
			err := config.Set("configPath", Flags.configFile)
			if err != nil {
				logger.Log.Error(err.Error())
			}
		}
	} else {
		logger.Log.Info("No config file found")
	}

	// aws roles to assume, from yaml config
	for i := 0; i < len(config.Slices("awsAssumedRoles")); i++ {

		p := config.Slices("awsAssumedRoles")[i]

		profile := awsProfile{
			Name:    p.String("name"),
			RoleArn: p.String("roleArn")}
		AwsProfiles = append(AwsProfiles, profile)

	}

	_, profileInfo, err := AwsGetIdentity(nil)
	if err != nil {
		logger.Log.Error(err.Error())
	} else {
		logger.Log.Info("AWS identity is " + profileInfo)
	}

	// yaml config is overwrittent by env
	err = config.Load(env.Provider("CC_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "CC_")), "_", ".", -1)
	}), nil)
	if err != nil {
		logger.Log.Error(err.Error())
	}

}

func AwsGetIdentity(assumedRole *credentials.Credentials) (int, string, error) {
	aws_session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String("eu-west-3")},
	}))

	var stscli *sts.STS
	if assumedRole != nil {
		aws_session.Config.Credentials = assumedRole
		stscli = sts.New(aws_session, &aws.Config{Credentials: assumedRole})
	} else {
		stscli = sts.New(aws_session)
	}

	input := &sts.GetCallerIdentityInput{}

	result, err := stscli.GetCallerIdentity(input)
	if err != nil {
		logger.Log.Error(err.Error())
		return 0, "", err

	} else {
		accountId, _ := strconv.Atoi(*result.Account)
		return accountId, *result.Arn, nil
	}
}
