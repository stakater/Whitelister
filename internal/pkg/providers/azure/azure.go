package azure

import (
	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/network/mgmt/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
	"github.com/Azure/go-autorest/autorest"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/stakater/Whitelister/internal/pkg/config"
	clientset "k8s.io/client-go/kubernetes"

	"github.com/stakater/Whitelister/internal/pkg/utils"
)

// Azure provider class implementing the Provider interface
type Azure struct {
	ClientSet                 clientset.Interface
	resourcesClient           resources.Client
	securityGroupClient       network.SecurityGroupsClient
	securityRulesClient       network.SecurityRulesClient
	authorizer                autorest.Authorizer
	SubscriptionID            string
	ClientID                  string
	ClientSecret              string
	TenantID                  string
	ResourceGroupName         string
	RemoveRule                bool
	KeepRuleDescriptionPrefix string
}

// GetName Returns name of provider
func (a *Azure) GetName() string {
	return "Azure"
}

// Init initializes the Azure Provider Configuration like client id and client secret
func (a *Azure) Init(params map[interface{}]interface{}, clientSet clientset.Interface) error {
	err := mapstructure.Decode(params, &a) //Converts the params to Azure struct fields
	if err != nil {
		return err
	}

	err = a.initializeAzureClients() // initializes azure clients for whitelisting ips
	if err != nil {
		return err
	}

	return nil
}

// WhiteListIps - Get List of IP addresses to whitelist
func (a *Azure) WhiteListIps(filter config.Filter, ipPermissions []utils.IpPermission) error {

	resources, err := a.fetchResources(filter)
	if err != nil {
		logrus.Error("Error fetching resources for the given filter")
		return err
	}

	for _, resource := range resources {
		logrus.Infof("Name of the security group %s", *resource.Name)
		securityGroup, err := a.fetchSecurityGroup(*resource.Name)
		if err != nil {
			logrus.Errorf("Error fetching security group for resource %s", *resource.Name)
			return err
		}
		err = a.updateSecurityRules(*securityGroup, ipPermissions)
		if err != nil {
			logrus.Errorf("Error whitelisting ips for security group %s", *resource.Name)
			return err
		}
	}
	return nil
}
