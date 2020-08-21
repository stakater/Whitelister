package azure

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"

	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/stakater/Whitelister/internal/pkg/config"
	clientset "k8s.io/client-go/kubernetes"

	"github.com/stakater/Whitelister/internal/pkg/utils"
)

// Azure provider class implementing the Provider interface
type Azure struct {
	ClientSet                 clientset.Interface
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
	a.ClientSet = clientSet
	err := mapstructure.Decode(params, &a) //Converts the params to Azure struct fields
	if err != nil {
		return err
	}
	return nil
}

// WhiteListIps - Get List of IP addresses to whitelist
func (a *Azure) WhiteListIps(filter config.Filter, ipPermissions []utils.IpPermission) error {

	resourceGroups, err := fetchResourceGroups(a, filter)
	if err != nil {
		logrus.Error("Error fetching resource groups for given filter")
		return err
	}
	for _, resourceGroup := range resourceGroups.Values() {
		resourceName := *resourceGroup.Name
		logrus.Infof("Name of the resource %s", *resourceGroup.Name)

		securityGroup, err := fetchSecurityGroup(a, resourceName)
		if err != nil {
			logrus.Error("Error fetching security group for resource " + resourceName)
			return err
		}

		rulesClient, err := getSecurityRulesClient(a)
		if err != nil {
			logrus.Error("Error creating security rules client")
			return err
		}

		for _, existingSecurityRule := range *securityGroup.SecurityGroupPropertiesFormat.SecurityRules {
			existingRuleName := *existingSecurityRule.Name
			logrus.Infof("Found security group %s", existingRuleName)
			if !isSecurityRuleToBeRetained(existingRuleName, ipPermissions, a.KeepRuleDescriptionPrefix, a.RemoveRule) {
				err = deleteRule(rulesClient, a, resourceName, existingRuleName)
				if err != nil {
					logrus.Error("Error deleting the security rule " + existingRuleName)
					return err
				}
			}
		}

		for _, ipPermission := range ipPermissions {
			fromPortStr := strconv.FormatInt(*ipPermission.FromPort, 10)
			toPortStr := strconv.FormatInt(*ipPermission.ToPort, 10)

			for _, ipRange := range *&ipPermission.IpRanges {
				ipDescription := ipRange.Description
				ipCidr := ipRange.IpCidr

				err = createSecurityRule(a, rulesClient, resourceName, ipDescription, fromPortStr, toPortStr, ipCidr, *ipPermission.IpProtocol)
				if err != nil {
					logrus.Error("Error adding security rule for azure")
					return err
				}
			}
		}
	}
	return nil
}

func fetchResourceGroups(azure *Azure, filter config.Filter) (*resources.ListResultPage, error) {
	if filter.FilterType == config.SecurityGroup {
		return fetchResourceGroupsByFilter(azure, filter)
	} else {
		logrus.Error("filter type " + filter.FilterType.String() + " not supported for azure yet")
		return nil, errors.New("filter type " + filter.FilterType.String() + " not supported for azure yet")
	}
}

func fetchResourceGroupsByFilter(azure *Azure, filter config.Filter) (result *resources.ListResultPage, err error) {
	resourcesClient, err := getResourcesClient(azure)
	if err != nil {
		return nil, err
	}
	tagFilter := fmt.Sprintf("tagName eq '%s' and tagValue eq '%s'", filter.LabelName, filter.LabelValue)
	resourceGroups, err := resourcesClient.ListByResourceGroup(context.TODO(), azure.ResourceGroupName, tagFilter, "SecurityRules", nil)
	if err != nil {
		return nil, err
	}
	return &resourceGroups, nil
}

func fetchSecurityGroup(a *Azure, resourceName string) (securityGroup network.SecurityGroup, err error) {
	securityGroupClient, err := getSecurityGroupsClient(a)
	if err != nil {
		return network.SecurityGroup{}, err
	}
	securityGroup, err = securityGroupClient.Get(context.TODO(), a.ResourceGroupName, resourceName, "")
	if err != nil {
		return network.SecurityGroup{}, err
	}
	return securityGroup, err
}

func getSecurityRulesClient(a *Azure) (network.SecurityRulesClient, error) {

	authorizer, err := auth.NewClientCredentialsConfig(a.ClientID, a.ClientSecret, a.TenantID).Authorizer()
	if err != nil {
		return network.SecurityRulesClient{}, err
	}
	securityRulesClient := network.NewSecurityRulesClient(a.SubscriptionID)
	securityRulesClient.Authorizer = authorizer

	return securityRulesClient, nil
}

func deleteRule(rulesClient network.SecurityRulesClient, a *Azure, resourceName string, ipDescription string) error {
	futureT, err := rulesClient.Delete(context.TODO(), a.ResourceGroupName, resourceName, ipDescription)
	if err != nil {
		return err
	}
	err = futureT.WaitForCompletionRef(context.TODO(), rulesClient.Client)
	if err != nil {
		return err
	}
	return nil
}

func getResourcesClient(a *Azure) (resources.Client, error) {

	authorizer, err := auth.NewClientCredentialsConfig(a.ClientID, a.ClientSecret, a.TenantID).Authorizer()
	if err != nil {
		return resources.Client{}, err
	}
	resourcesClient := resources.NewClient(a.SubscriptionID)
	resourcesClient.Authorizer = authorizer

	return resourcesClient, nil
}

func getSecurityGroupsClient(a *Azure) (*network.SecurityGroupsClient, error) {

	authorizer, err := auth.NewClientCredentialsConfig(a.ClientID, a.ClientSecret, a.TenantID).Authorizer()
	if err != nil {
		return nil, err
	}
	nsgClient := network.NewSecurityGroupsClient(a.SubscriptionID)
	nsgClient.Authorizer = authorizer

	return &nsgClient, nil
}

func createSecurityRule(a *Azure, rulesClient network.SecurityRulesClient, resourceName string, ipDescription *string, fromPortStr string, toPortStr string, ipRange *string, ipProtocol string) error {

	var networkProtocol network.SecurityRuleProtocol

	switch ipProtocol {
	case "tcp":
		networkProtocol = network.SecurityRuleProtocolTCP
	case "udp":
		networkProtocol = network.SecurityRuleProtocolUDP
	case "icmp":
		networkProtocol = network.SecurityRuleProtocolIcmp
	case "esp":
		networkProtocol = network.SecurityRuleProtocolEsp
	case "*":
		networkProtocol = network.SecurityRuleProtocolAsterisk
	case "ah":
		networkProtocol = network.SecurityRuleProtocolAh
	default:
		return errors.New("ip protocol " + ipProtocol + " unidentified")
	}

	future, err := rulesClient.CreateOrUpdate(
		context.TODO(),
		a.ResourceGroupName,
		resourceName,
		*ipDescription,
		network.SecurityRule{
			SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
				Access:                   network.SecurityRuleAccessAllow,
				DestinationAddressPrefix: to.StringPtr("*"),
				DestinationPortRange:     to.StringPtr(fromPortStr + "-" + toPortStr),
				Direction:                network.SecurityRuleDirectionInbound,
				Description:              ipDescription,
				Priority:                 to.Int32Ptr(100),
				Protocol:                 networkProtocol,
				SourceAddressPrefix:      ipRange,
				SourcePortRange:          to.StringPtr("*"),
			},
		})
	if err != nil {
		return err
	}
	err = future.WaitForCompletionRef(context.TODO(), rulesClient.Client)
	if err != nil {
		return err
	}
	return nil
}

func isSecurityRuleToBeRetained(existingRuleName string, ipPermissions []utils.IpPermission, keepRuleDescriptionPrefix string, removeRule bool) bool {
	return !removeRule || strings.HasPrefix(existingRuleName, keepRuleDescriptionPrefix)
}
