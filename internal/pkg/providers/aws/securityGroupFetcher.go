package aws

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/sirupsen/logrus"
)

func (a *Aws) getSecurityGroupsByLoadBalancer(session *session.Session, credentials *credentials.Credentials, resourceIds []string) ([]*ec2.SecurityGroup, error) {

	// Create an ELB service client.
	elbClient := elb.New(session, &aws.Config{
		Credentials: credentials,
		Region:      aws.String(a.Region),
	})

	result, err := elbClient.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{
		LoadBalancerNames: aws.StringSlice(resourceIds),
	})
	if err != nil {
		logrus.Errorf("%v", err)
		return nil, err
	}

	if len(result.LoadBalancerDescriptions) == 0 {
		return nil, errors.New("no load balancer found with AWS")
	}

	var securityGroupNames []*string
	for _, loadBalancerDescription := range result.LoadBalancerDescriptions {
		securityGroupNames = append(securityGroupNames, loadBalancerDescription.SourceSecurityGroup.GroupName)
	}

	ec2Client := getEc2Client(session, credentials, a)

	var vpcFilter = "vpc-id"
	var groupFilter = "group-name"

	securityGroupResult, err := ec2Client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   &vpcFilter,
				Values: []*string{result.LoadBalancerDescriptions[0].VPCId},
			},
			{
				Name:   &groupFilter,
				Values: securityGroupNames,
			},
		},
	})

	if err != nil {
		logrus.Errorf("%v", err)
		return nil, err
	}

	return securityGroupResult.SecurityGroups, nil
}

func (a *Aws) getSecurityGroupsByTagFilter(session *session.Session, credentials *credentials.Credentials, filterLabel []string) ([]*ec2.SecurityGroup, error) {

	ec2Client := getEc2Client(session, credentials, a)
	filters := a.getSearchFilterWithTag(filterLabel)

	securityGroupResult, err := ec2Client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{Filters: filters})

	if err != nil {
		logrus.Errorf("%v", err)
		return nil, err
	}

	return securityGroupResult.SecurityGroups, nil
}

func getEc2Client(session *session.Session, credentials *credentials.Credentials, a *Aws) *ec2.EC2 {
	return ec2.New(session, &aws.Config{
		Credentials: credentials,
		Region:      aws.String(a.Region),
	})
}

func (a *Aws) getSearchFilterWithTag(filterLabel []string) []*ec2.Filter {
	filters := make([]*ec2.Filter, 0)
	keyName := "tag:" + filterLabel[0]
	filter := ec2.Filter{
		Name:   &keyName,
		Values: []*string{&filterLabel[1]}}
	filters = append(filters, &filter)
	return filters
}
