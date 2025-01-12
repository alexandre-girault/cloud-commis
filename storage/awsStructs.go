package storage

import "time"

type Aws_scans struct {
	AwsScanDate time.Time
	AwsAccounts map[int]Aws_account_scan
	AwsImages   map[string]AwsImage
}

type Aws_account_scan struct {
	AwsRegions map[string]Aws_region_scan
}

type Aws_region_scan struct {
	VirtualMachines map[string]VirtualMachine
}

type VirtualMachine struct {
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
	Description     string
	OwnerId         string
	DeprecationTime string
}
