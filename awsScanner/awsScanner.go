package awsScanner

import (
	"cloud-commis/config"
	"cloud-commis/logger"
	"cloud-commis/storage"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func ScheduledScan() {

	var AwsTicker = time.NewTicker(time.Duration(config.ParsedData.Int("scanIntervalMin")) * time.Minute)
	done := make(chan bool)

	for range AwsTicker.C {
		logger.Log.Info("AWS scan started")
		Aws_instances_inventory()
	}

	//ticker.Stop()
	done <- true
}

// scan each AWS account for ec2 instances
func Aws_instances_inventory() {

	var scans storage.Aws_scans

	for _, profile := range config.AwsProfiles {

		logger.Log.Info("scanning AWS account " + profile.Name + " with role " + profile.RoleArn)

		aws_profile_session := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config:            aws.Config{Region: aws.String("eu-west-3")},
		}))

		assumedRole := stscreds.NewCredentials(aws_profile_session, profile.RoleArn)

		awsClient := ec2.New(aws_profile_session, &aws.Config{Credentials: assumedRole})

		scans.AwsAccounts = append(scans.AwsAccounts, aws_account_scan(awsClient, assumedRole))
		scans.AwsScanDate = time.Now()
	}

	logger.Log.Info("writing data to storage")
	storage.Data.Write(scans)

}

func aws_account_scan(awsClient *ec2.EC2, assumedRole *credentials.Credentials) storage.Aws_account_scan {

	var scan storage.Aws_account_scan
	var wg sync.WaitGroup

	region_result, err := awsClient.DescribeRegions(nil)
	if err != nil {
		logger.Log.Error("aws region scan failure")
		logger.Log.Error(err.Error())
	} else {
		// search for instances on every regions with go routines

		outputs := make(chan storage.Aws_region_scan, len(region_result.Regions))

		for count := range region_result.Regions {
			wg.Add(1)
			regionName := *region_result.Regions[count].RegionName

			//regionalAwsSessions := make(map[string]*session.Session)

			regionalAwsSessions, _ := session.NewSessionWithOptions(session.Options{
				Config: aws.Config{Region: aws.String(regionName)},
			})

			regionalEc2Client := ec2.New(regionalAwsSessions, &aws.Config{Credentials: assumedRole})
			go aws_region_scan(regionName, regionalEc2Client, outputs, &wg)
		}

		logger.Log.Info("Waiting for all regions to finish")
		wg.Wait()

		logger.Log.Info("End of AWS scan")

		scan.AwsAccountID, _, _ = config.AwsGetIdentity()

		for message := range len(region_result.Regions) {
			scan.AwsRegions = append(scan.AwsRegions, <-outputs)
			logger.Log.Debug(fmt.Sprint(message))
		}

	}
	return scan
}

func aws_region_scan(aws_region string, awsClient *ec2.EC2, channel chan storage.Aws_region_scan, wg *sync.WaitGroup) {

	defer wg.Done()
	logger.Log.Debug("scanning region " + aws_region)

	var aws_region_scan_result storage.Aws_region_scan
	aws_region_scan_result.RegionName = *awsClient.Config.Region

	// Call to get detailed information on each instance
	result, err := awsClient.DescribeInstances(nil)
	if err != nil {
		logger.Log.Error("Error" + err.Error())
	} else {
		logger.Log.Debug(result.String())
		aws_region_scan_result = ec2ScanParse(result, aws_region)
	}
	channel <- aws_region_scan_result
}

func ec2ScanParse(result *ec2.DescribeInstancesOutput, aws_region string) storage.Aws_region_scan {

	var parsedOutput storage.Aws_region_scan

	parsedOutput.RegionName = aws_region

	for reservation := range result.Reservations {
		for instance := range result.Reservations[reservation].Instances {
			logger.Log.Info("found virtualmachine " + *result.Reservations[reservation].Instances[instance].InstanceId)
			found_vm := storage.VirtualMachine{
				InstanceId:               *result.Reservations[reservation].Instances[instance].InstanceId,
				Architecture:             *result.Reservations[reservation].Instances[instance].Architecture,
				LaunchTime:               *result.Reservations[reservation].Instances[instance].LaunchTime,
				UsageOperationUpdateTime: *result.Reservations[reservation].Instances[instance].UsageOperationUpdateTime,
				PlatformDetails:          *result.Reservations[reservation].Instances[instance].PlatformDetails,
				ImageId:                  *result.Reservations[reservation].Instances[instance].ImageId,
				InstanceType:             *result.Reservations[reservation].Instances[instance].InstanceType,
				//PublicIpAddress:          string(*result.Reservations[reservation].Instances[instance].PublicIpAddress),
				State: *result.Reservations[reservation].Instances[instance].State.Name,
			}
			for _, tag := range result.Reservations[reservation].Instances[instance].Tags {
				if *tag.Key == "Name" {
					found_vm.Name = *tag.Value
				}
			}

			parsedOutput.VirtualMachines = append(parsedOutput.VirtualMachines, found_vm)
		}
	}
	logger.Log.Debug("End of data parsing")
	return parsedOutput
}
