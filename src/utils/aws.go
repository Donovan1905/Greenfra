package utils

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type InstanceInfo struct {
	InstanceType string
	VCPUs        int32
	Memory       int64
}

func LoadAWSConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return cfg
}

func CreateEC2Client(cfg aws.Config) *ec2.Client {
	return ec2.NewFromConfig(cfg)
}

func DescribeInstanceTypes(svc *ec2.Client, instanceTypes []string) []InstanceInfo {
	input := &ec2.DescribeInstanceTypesInput{
		InstanceTypes: make([]types.InstanceType, len(instanceTypes)),
	}

	for i, instanceType := range instanceTypes {
		input.InstanceTypes[i] = types.InstanceType(instanceType)
	}

	result, err := svc.DescribeInstanceTypes(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to describe instance types, %v", err)
	}

	var instancesInfo []InstanceInfo
	for _, instance := range result.InstanceTypes {
		instancesInfo = append(instancesInfo, InstanceInfo{
			InstanceType: string(instance.InstanceType),
			VCPUs:        *instance.VCpuInfo.DefaultVCpus,
			Memory:       *instance.MemoryInfo.SizeInMiB,
		})
	}

	return instancesInfo
}
