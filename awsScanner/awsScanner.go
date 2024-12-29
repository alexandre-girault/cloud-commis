package awsScanner

import (
	"cloud-commis/config"
	"cloud-commis/logger"
	"cloud-commis/storage"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
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

// scan aws regions for ec2 instances
func Aws_instances_inventory() {

	var wg sync.WaitGroup

	aws_session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String("eu-west-3")},
	}))

	ec2Svc := ec2.New(aws_session)
	regionalAwsSessions := make(map[string]*session.Session)

	region_result, err := ec2Svc.DescribeRegions(nil)
	if err != nil {
		logger.Log.Error("aws region scan failure")
		logger.Log.Error(err.Error())
	} else {
		// search for instances on every regions with go routines
		outputs := make(chan storage.Aws_region_scan, len(region_result.Regions))
		for count := range region_result.Regions {
			wg.Add(1)
			regionName := *region_result.Regions[count].RegionName

			regionalAwsSessions[regionName], _ = session.NewSessionWithOptions(session.Options{
				Config: aws.Config{Region: aws.String(regionName)},
			})
			regionaEc2Client := ec2.New(regionalAwsSessions[regionName])
			go aws_region_scan(regionName, regionaEc2Client, outputs, &wg)
		}
		wg.Wait()

		logger.Log.Info("End of AWS scan")
		var scan storage.Aws_scans
		for message := range len(region_result.Regions) {
			scan.Data = append(scan.Data, <-outputs)
			logger.Log.Debug(fmt.Sprint(message))
		}
		storage.Data.Write(scan)
	}
}

func aws_region_scan(aws_region string, awsClient *ec2.EC2, channel chan storage.Aws_region_scan, wg *sync.WaitGroup) {

	defer wg.Done()
	logger.Log.Debug("scanning region " + aws_region)

	var aws_region_scan_result storage.Aws_region_scan
	aws_region_scan_result.RegionName = *awsClient.Config.Region

	// Call to get detailed information on each instance
	result, err := awsClient.DescribeInstances(nil)
	if err != nil {
		fmt.Println("Error", err)
	} else {
		fmt.Println(result)
		aws_region_scan_result = ec2ScanParse(result)
	}
	channel <- aws_region_scan_result
}

func ec2ScanParse(result *ec2.DescribeInstancesOutput) storage.Aws_region_scan {

	var parsedOutput storage.Aws_region_scan

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
