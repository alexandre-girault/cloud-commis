package storage

import "time"

type Aws_scans struct {
	Data []Aws_region_scan
}

type Aws_region_scan struct {
	RegionName      string
	VirtualMachines []VirtualMachine
	AwsImages       []AwsImage
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
	ImageId         string
	Description     string
	OwnerId         string
	DeprecationTime string
}
