package storage

import "time"

type Aws_scans struct {
	AwsScanDate time.Time
	AwsAccounts []Aws_account_scan
	AwsImages   []AwsImage
}

type Aws_account_scan struct {
	AwsAccountID int
	AwsRegions   []Aws_region_scan
}

type Aws_region_scan struct {
	RegionName      string
	VirtualMachines []VirtualMachine
}

type VirtualMachine struct {
	InstanceId               string
	Name                     string
	Architecture             string
	LaunchTime               time.Time
	UsageOperationUpdateTime time.Time
	PlatformDetails          string
	ImageId                  string
	InstanceType             string
	PublicIpAddress          string
	State                    string
	KernelImage              string
	BootImage                string
}

type AwsImage struct {
	Name            string
	Region          string
	ImageId         string
	Description     string
	OwnerId         string
	DeprecationTime string
}
