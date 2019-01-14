package aws

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"

	"github.com/stakater/Whitelister/internal/pkg/utils"
)

// Aws provider class implementing the Provider interface
type Aws struct {
	RoleArn string
	Region  string
}

var vpcFilter string = "vpc-id"
var groupFilter string = "group-name"

func (a *Aws) GetName() string {
	return "Amazon Web Services"
}

// Init initializes the Aws Provider Configuration like Access Token and Reion
func (a *Aws) Init(params map[interface{}]interface{}) error {
	err := mapstructure.Decode(params, &a) //Converts the params to Aws struct fields
	if err != nil {
		return err
	}
	if a.RoleArn == "" || a.Region == "" {
		return errors.New("Missing Aws Assume Role ARN or Region")
	}
	return nil
}

// Get List of IP addresses to whitelist
func (a *Aws) WhiteListIps(resourceIds []string, ipPermissions []utils.IpPermission) error {

	// Initial credentials loaded from SDK's default credential chain. Such as
	// the environment, shared credentials (~/.aws/credentials), or EC2 Instance
	// Role. These credentials will be used to to make the STS Assume Role API.
	session, err := session.NewSession()
	if err != nil {
		logrus.Errorf("%v", err)
		return err
	}

	// Create the credentials from AssumeRoleProvider to assume the role
	// referenced by the "myRoleARN" ARN.
	roleCredentials := stscreds.NewCredentials(session, a.RoleArn)

	securityGroups, err := getSecurityGroups(session, roleCredentials, resourceIds)
	if err != nil {
		logrus.Errorf("%v", err)
		return err
	}

	if len(securityGroups) != 1 {
		logrus.Errorf("FIX ME")
		return nil
	}

	ec2IpPermissions := getEc2IpPermissions(ipPermissions)

	ec2Client := ec2.New(session, &aws.Config{Credentials: roleCredentials})

	for _, securityGroup := range securityGroups {
		updateSecurityGroup(ec2Client, securityGroup, ec2IpPermissions)
	}
	return nil
}

func getSecurityGroups(session *session.Session, credentials *credentials.Credentials,
	resourceIds []string) ([]*ec2.SecurityGroup, error) {

	// Create an ELB service client.
	elbClient := elb.New(session, &aws.Config{Credentials: credentials})

	result, err := elbClient.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{
		LoadBalancerNames: aws.StringSlice(resourceIds),
	})
	if err != nil {
		logrus.Errorf("%v", err)
		return nil, err
	}

	if len(result.LoadBalancerDescriptions) == 0 {
		return nil, errors.New("No Load Balancer Found with AWS")
	}

	securityGroupNames := []*string{}
	for _, loadBalancerDescription := range result.LoadBalancerDescriptions {
		securityGroupNames = append(securityGroupNames, loadBalancerDescription.SourceSecurityGroup.GroupName)
	}

	ec2Client := ec2.New(session)

	securityGroupResult, err := ec2Client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   &vpcFilter,
				Values: []*string{result.LoadBalancerDescriptions[0].VPCId},
			},
			{
				Name:   &groupFilter,
				Values: aws.StringSlice(resourceIds),
			},
		},
	})

	return securityGroupResult.SecurityGroups, nil
}

func updateSecurityGroup(client *ec2.EC2, securityGroup *ec2.SecurityGroup,
	ipPermissions []*ec2.IpPermission) error {

	removeSecurityRules(client, securityGroup, ipPermissions)
	addSecurityRules(client, securityGroup, ipPermissions)

	return nil
}

func removeSecurityRules(client *ec2.EC2, securityGroup *ec2.SecurityGroup,
	ipPermissions []*ec2.IpPermission) {
	var removeIpPermission bool
	ipPermissionsToRemove := []*ec2.IpPermission{}

	for _, securityGroupIpPermission := range securityGroup.IpPermissions {
		removeIpPermission = true
		for _, ipPermission := range ipPermissions {
			if utils.IsEc2IpPermissionEqual(ipPermission, securityGroupIpPermission) {
				removeIpPermission = false
				break
			}
		}
		if removeIpPermission {
			ipPermissionsToRemove = append(ipPermissionsToRemove, securityGroupIpPermission)
		}
	}
	if len(ipPermissionsToRemove) > 0 {
		logrus.Infof("Removing security rules : %v for security group :%s", ipPermissionsToRemove, *securityGroup.GroupName)
		err := removeSecurityGroupIngresses(client, securityGroup, ipPermissionsToRemove)
		if err != nil {
			logrus.Errorf("Error removing security rules for security group %s : %v", *securityGroup.GroupName, err)
		}
	} else {
		logrus.Infof("No security rules to remove for security group : %s", *securityGroup.GroupName)
	}
}

func addSecurityRules(client *ec2.EC2, securityGroup *ec2.SecurityGroup,
	ipPermissions []*ec2.IpPermission) {
	var ipPermissionExists bool
	ipPermissionsToAdd := []*ec2.IpPermission{}

	for _, ipPermission := range ipPermissions {
		ipPermissionExists = false
		for _, securityGroupIpPermission := range securityGroup.IpPermissions {
			if utils.IsEc2IpPermissionEqual(ipPermission, securityGroupIpPermission) {
				ipPermissionExists = true
				break
			}
		}
		if !ipPermissionExists {
			ipPermissionsToAdd = append(ipPermissionsToAdd, ipPermission)
		}
	}
	if len(ipPermissionsToAdd) > 0 {
		logrus.Infof("Adding security rules : %v for security group :%s", ipPermissionsToAdd, *securityGroup.GroupName)
		err := addSecurityGroupIngresses(client, securityGroup, ipPermissionsToAdd)
		if err != nil {
			logrus.Errorf("Error adding security rules for security group %s : %v", *securityGroup.GroupName, err)
		}
	} else {
		logrus.Infof("No security rules to add for security group : %s", *securityGroup.GroupName)
	}

}

func removeSecurityGroupIngresses(client *ec2.EC2, securityGroup *ec2.SecurityGroup,
	ipPermissions []*ec2.IpPermission) error {

	// _, err := client.RevokeSecurityGroupIngress(&ec2.RevokeSecurityGroupIngressInput{
	// 	GroupId:       securityGroup.GroupId,
	// 	IpPermissions: ipPermissions,
	// })

	return nil
}

func addSecurityGroupIngresses(client *ec2.EC2, securityGroup *ec2.SecurityGroup,
	ipPermissions []*ec2.IpPermission) error {

	// _, err := client.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
	// 	GroupId:       securityGroup.GroupId,
	// 	IpPermissions: ipPermissions,
	// })

	return nil
}

func getEc2IpPermissions(ipPermissions []utils.IpPermission) []*ec2.IpPermission {

	ec2IpPermissions := []*ec2.IpPermission{}
	for _, ipPermission := range ipPermissions {
		ec2IpPermissions = append(ec2IpPermissions,
			(&ec2.IpPermission{}).
				SetIpProtocol(*ipPermission.IpProtocol).
				SetFromPort(*ipPermission.FromPort).
				SetToPort(*ipPermission.ToPort).
				SetIpRanges([]*ec2.IpRange{
					{
						CidrIp:      ipPermission.IpCidr,
						Description: ipPermission.Description,
					},
				}),
		)
	}

	return ec2IpPermissions
}
