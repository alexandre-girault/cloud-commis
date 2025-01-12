package awsScanner

import (
	"cloud-commis/logger"
	"cloud-commis/storage"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func awsAmiScan(awsClient *ec2.EC2, imageIds []string, region string) []storage.AwsImage {

	var awsImages []storage.AwsImage

	result, err := awsClient.DescribeImages(&ec2.DescribeImagesInput{
		ImageIds: aws.StringSlice(imageIds),
	})

	if err != nil {
		logger.Log.Error(err.Error())
	} else {
		for Ami := range result.Images {
			foundAmi := &storage.AwsImage{
				Name:            *result.Images[Ami].Name,
				Region:          region,
				Description:     *result.Images[Ami].Description,
				OwnerId:         *result.Images[Ami].OwnerId,
				DeprecationTime: *result.Images[Ami].DeprecationTime,
			}
			awsImages = append(awsImages, *foundAmi)
		}
	}
	return awsImages
}
