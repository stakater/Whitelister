package aws

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/stakater/Whitelister/internal/pkg/config"

	"github.com/stakater/Whitelister/internal/pkg/utils"
)

// Aws provider class implementing the Provider interface
type Aws struct {
	RoleArn                   string
	Region                    string
	RemoveRule                bool
	KeepRuleDescriptionPrefix string
}

// GetName Returns name of provider
func (a *Aws) GetName() string {
	return "Amazon Web Services"
}

// Init initializes the Aws Provider Configuration like Access Token and Region
func (a *Aws) Init(params map[interface{}]interface{}) error {
	err := mapstructure.Decode(params, &a) //Converts the params to Aws struct fields
	if err != nil {
		return err
	}
	if a.RoleArn == "" || a.Region == "" {
		return errors.New("missing Aws Assume Role ARN or Region")
	}
	return nil
}

// WhiteListIps - Get List of IP addresses to whitelist
func (a *Aws) WhiteListIps(filterType config.FilterType, resourceIds []string, ipPermissions []utils.IpPermission) error {

	// Initial credentials loaded default credential chain from SDK. Such as
	// the environment, shared credentials (~/.aws/credentials), or EC2 Instance
	// Role. These credentials will be used to to make the STS Assume Role API.
	awsSession, err := session.NewSession()
	if err != nil {
		logrus.Errorf("%v", err)
		return err
	}

	// Create the credentials from AssumeRoleProvider to assume the role
	// referenced by the "myRoleARN" ARN.
	roleCredentials := stscreds.NewCredentials(awsSession, a.RoleArn)
	securityGroups, err := a.fetchSecurityGroup(filterType, awsSession, roleCredentials, resourceIds)

	if err != nil {
		logrus.Errorf("%v", err)
		return err
	}

	if len(securityGroups) != 1 {
		logrus.Errorf("FIX ME : %v", securityGroups)
		return nil
	}

	ec2IpPermissions := getEc2IpPermissions(ipPermissions)

	ec2Client := ec2.New(awsSession, &aws.Config{
		Credentials: roleCredentials,
		Region:      aws.String(a.Region),
	})

	for _, securityGroup := range securityGroups {
		err := a.updateSecurityGroup(ec2Client, securityGroup, ec2IpPermissions)
		if err != nil {
			logrus.Errorf("%v", err)
		}
	}
	return nil
}
