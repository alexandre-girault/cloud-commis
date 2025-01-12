package awsScanner

import (
	"cloud-commis/config"
	"cloud-commis/logger"
	"cloud-commis/storage"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"golang.org/x/exp/maps"
)

var imageScans map[string]storage.AwsImage

func ScheduledScan() {

	var AwsTicker = time.NewTicker(time.Duration(config.ParsedData.Int("scanIntervalMin")) * time.Minute)
	done := make(chan bool)

	for range AwsTicker.C {
		logger.Log.Info("AWS scan started")
		Aws_instances_inventory()
	}

	done <- true
}

// scan each AWS account for ec2 instances
func Aws_instances_inventory() {

	var scans storage.Aws_scans
	scans.AwsAccounts = make(map[int]storage.Aws_account_scan)
	scans.AwsImages = make(map[string]storage.AwsImage)
	imageScans = make(map[string]storage.AwsImage)
	scans.AwsScanDate = time.Now()

	//scan accounts with environment credentials

	accountId, identity, err := config.AwsGetIdentity(nil)
	if err != nil {
		logger.Log.Error("no valid aws credentials found")
	} else {
		logger.Log.Info("scanning AWS account " + strconv.Itoa(accountId) + "with identity : " + identity)

		aws_session := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config:            aws.Config{Region: aws.String("eu-west-3")},
		}))

		awsClient := ec2.New(aws_session)

		scans.AwsAccounts[accountId] = aws_account_scan(awsClient, nil)
	}

	//scan accounts with assumed roles from config
	if len(config.AwsProfiles) != 0 {
		for _, profile := range config.AwsProfiles {

			logger.Log.Info("scanning AWS account " + profile.Name + " with role " + profile.RoleArn)

			aws_profile_session := session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
				Config:            aws.Config{Region: aws.String("eu-west-3")},
			}))

			assumedRole := stscreds.NewCredentials(aws_profile_session, profile.RoleArn)

			awsClient := ec2.New(aws_profile_session, &aws.Config{Credentials: assumedRole})

			accountID, _, _ := config.AwsGetIdentity(assumedRole)
			scans.AwsAccounts[accountID] = aws_account_scan(awsClient, assumedRole)
		}
	} else {
		logger.Log.Info("No assumed roles found in config")
	}

	for awsImage := range imageScans {
		logger.Log.Debug("awsImage: " + imageScans[awsImage].Name)
		scans.AwsImages[awsImage] = imageScans[awsImage]
	}

	logger.Log.Info("writing data to storage")
	storage.Data.Write(scans)

}

func aws_account_scan(awsClient *ec2.EC2, assumedRole *credentials.Credentials) storage.Aws_account_scan {

	var scan storage.Aws_account_scan
	var wg sync.WaitGroup
	scan.AwsRegions = make(map[string]storage.Aws_region_scan)

	region_result, err := awsClient.DescribeRegions(nil)
	if err != nil {
		logger.Log.Error("aws region scan failure")
		logger.Log.Error(err.Error())
	} else {
		// search for instances on each regions with a go routine

		outputs := make(chan map[string]storage.Aws_region_scan, len(region_result.Regions))

		var regionalEc2Client *ec2.EC2
		for count := range region_result.Regions {
			wg.Add(1)
			regionName := *region_result.Regions[count].RegionName

			regionalAwsSessions, _ := session.NewSessionWithOptions(session.Options{
				Config: aws.Config{Region: aws.String(regionName)},
			})

			if assumedRole != nil {
				regionalEc2Client = ec2.New(regionalAwsSessions, &aws.Config{Credentials: assumedRole})
			} else {
				regionalEc2Client = ec2.New(regionalAwsSessions)
			}
			go aws_region_scan(regionName, regionalEc2Client, outputs, &wg)
		}

		logger.Log.Info("Waiting for all regions to finish")
		wg.Wait()

		logger.Log.Info("End of AWS scan")

		for message := range len(region_result.Regions) {

			data := <-outputs
			regions := maps.Keys(data)
			for i := range regions {
				scan.AwsRegions[regions[i]] = data[regions[i]]
			}

			logger.Log.Debug(fmt.Sprint(message))
		}

	}
	return scan
}

func aws_region_scan(aws_region string, awsClient *ec2.EC2, channel chan map[string]storage.Aws_region_scan, wg *sync.WaitGroup) {

	defer wg.Done()
	logger.Log.Debug("scanning region " + aws_region)

	aws_region_scan_result := make(map[string]storage.Aws_region_scan)
	//aws_region_scan_result.RegionName = *awsClient.Config.Region

	// Call to get detailed information on each instance
	result, err := awsClient.DescribeInstances(nil)
	if err != nil {
		logger.Log.Error("Error" + err.Error())
	} else {
		logger.Log.Debug(result.String())
		aws_region_scan_result[aws_region] = ec2ScanParse(result)
	}

	// get details of AMIs
	for _, vm := range aws_region_scan_result[aws_region].VirtualMachines {

		_, isScanned := imageScans[vm.ImageId]
		if isScanned {
			logger.Log.Debug("Image already scanned: " + vm.ImageId)
		} else {
			logger.Log.Debug("Scanning image " + vm.ImageId)
			imageScans[vm.ImageId] = awsAmiScan(awsClient, []string{vm.ImageId}, aws_region)[0]
		}
	}
	channel <- aws_region_scan_result
}

func ec2ScanParse(result *ec2.DescribeInstancesOutput) storage.Aws_region_scan {

	var parsedOutput storage.Aws_region_scan
	parsedOutput.VirtualMachines = make(map[string]storage.VirtualMachine)
	var instanceId string

	for reservation := range result.Reservations {
		for instance := range result.Reservations[reservation].Instances {
			instanceId = *result.Reservations[reservation].Instances[instance].InstanceId
			logger.Log.Info("found virtualmachine " + *result.Reservations[reservation].Instances[instance].InstanceId)
			found_vm := storage.VirtualMachine{
				//InstanceId:               *result.Reservations[reservation].Instances[instance].InstanceId,
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

			parsedOutput.VirtualMachines[instanceId] = found_vm
		}
	}
	logger.Log.Debug("End of data parsing")
	return parsedOutput
}
