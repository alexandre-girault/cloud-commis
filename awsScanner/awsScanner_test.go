package awsScanner

import (
	"cloud-commis/storage"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestEc2ScanParse(t *testing.T) {
	var testData *ec2.DescribeInstancesOutput

	want := storage.VirtualMachine{
		InstanceId:               "i-08cb2bdddc801f6d8",
		Name:                     "aws-ec2-test-vm",
		Architecture:             "arm64",
		LaunchTime:               time.Date(2024, 01, 02, 15, 04, 05, 000000000, time.UTC),
		UsageOperationUpdateTime: time.Date(2024, 01, 02, 15, 04, 05, 000000000, time.UTC),
		PlatformDetails:          "Linux/UNIX",
		ImageId:                  "ami-01c5300f289d64643",
		InstanceType:             "t4g.nano",
		State:                    "stopped",
	}

	testFile, err := os.ReadFile("../testData/awsDescribeInstanceOutput-01.json")
	if err != nil {
		t.Error(err.Error())
	}

	err = json.Unmarshal(testFile, &testData)
	if err != nil {
		t.Error(err.Error())
	}

	result := ec2ScanParse(testData, "eu-west-1")

	if !reflect.DeepEqual(want, result.VirtualMachines[0]) {
		t.Error("Parsing fail to found the correct ec2 attributes, \nwant : \n" + fmt.Sprintf("%+v", want) + "\nget : \n" + fmt.Sprintf("%+v", result.VirtualMachines[0]))
	}
}
