package awsScanner

import (
	"cloud-commis/logger"
	"encoding/base64"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func AwsConsoleLogOutput(awsClient *ec2.EC2, instanceId string) (os string, kernel string) {

	returnOs := ""
	returnKernel := ""

	var compiledOsRegex []*regexp.Regexp
	var compiledKernelRegex []*regexp.Regexp

	osRegex := []string{
		"Booting `(Amazon Linux .+ 20[0-9][0-9])",
		"Welcome to .+(Ubuntu.*)!",
		"Welcome to .*(Red Hat .+)!",
	}
	kernelRegex := []string{
		"(vmlinuz-[0-9].[0-9].[0-9]-.+) root",
		"(Kernel .+ on an .+).*",
	}

	for i := range osRegex {
		re := regexp.MustCompile(osRegex[i])
		compiledOsRegex = append(compiledOsRegex, re)
	}

	for i := range kernelRegex {
		re := regexp.MustCompile(kernelRegex[i])
		compiledKernelRegex = append(compiledKernelRegex, re)
	}

	result, err := awsClient.GetConsoleOutput(&ec2.GetConsoleOutputInput{
		InstanceId: &instanceId})
	if err != nil {
		fmt.Println("Error", err)
	} else {
		decoded, err := base64.StdEncoding.DecodeString(*result.Output)
		if err != nil {
			fmt.Println("decode error:", err, decoded)
			return
		}

		for cOsRegex := range compiledOsRegex {
			osVersion := compiledOsRegex[cOsRegex].FindStringSubmatch(string(decoded))
			if len(osVersion) > 0 {
				returnOs = osVersion[0]
			}
		}
		for cKernRegex := range compiledKernelRegex {
			kernelVersion := compiledKernelRegex[cKernRegex].FindStringSubmatch(string(decoded))
			if len(kernelVersion) > 0 {
				returnKernel = kernelVersion[0]
			}
		}
	}
	logger.Log.Debug("console logs infos : " + returnKernel + returnOs)
	return returnOs, returnKernel
}
