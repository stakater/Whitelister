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
	clientset "k8s.io/client-go/kubernetes"

	"github.com/stakater/Whitelister/internal/pkg/utils"
)

// Aws provider class implementing the Provider interface
type Aws struct {
	ClientSet                 clientset.Interface
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
func (a *Aws) Init(params map[interface{}]interface{}, clientSet clientset.Interface) error {
	a.ClientSet = clientSet
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
func (a *Aws) WhiteListIps(filter config.Filter, ipPermissions []utils.IpPermission) error {

	// Initial credentials loaded from SDK's default credential chain. Such as
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
	securityGroups, err := a.fetchSecurityGroup(awsSession, roleCredentials, filter)

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
