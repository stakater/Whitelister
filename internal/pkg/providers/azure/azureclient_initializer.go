package azure

import (
	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/network/mgmt/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/sirupsen/logrus"
)

func (a *Azure) initializeAzureClients() error {

	resourcesClient, err := getResourcesClient(a)
	if err != nil {
		logrus.Error("Error creating azure resources client")
		return err
	}
	a.resourcesClient = *resourcesClient

	securityGroupsClient, err := getSecurityGroupsClient(a)
	if err != nil {
		logrus.Error("Error creating azure security group client")
		return err
	}
	a.securityGroupClient = *securityGroupsClient

	securityRulesClient, err := getSecurityRulesClient(a)
	if err != nil {
		logrus.Error("Error creating azure security rules client")
		return err
	}
	a.securityRulesClient = *securityRulesClient

	authorizer, err := auth.NewClientCredentialsConfig(a.ClientID, a.ClientSecret, a.TenantID).Authorizer()
	if err != nil {
		logrus.Error("Error creating azure authorizer using new client credentials config")
		return err
	}
	a.authorizer = authorizer

	return nil
}

func getResourcesClient(a *Azure) (*resources.Client, error) {

	resourcesClient := resources.NewClient(a.SubscriptionID)
	resourcesClient.Authorizer = a.authorizer

	return &resourcesClient, nil
}

func getSecurityGroupsClient(a *Azure) (*network.SecurityGroupsClient, error) {

	nsgClient := network.NewSecurityGroupsClient(a.SubscriptionID)
	nsgClient.Authorizer = a.authorizer

	return &nsgClient, nil
}

func getSecurityRulesClient(a *Azure) (*network.SecurityRulesClient, error) {

	securityRulesClient := network.NewSecurityRulesClient(a.SubscriptionID)
	securityRulesClient.Authorizer = a.authorizer

	return &securityRulesClient, nil
}
