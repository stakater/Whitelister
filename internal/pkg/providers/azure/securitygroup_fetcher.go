package azure

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/network/mgmt/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
	"github.com/sirupsen/logrus"
	"github.com/stakater/Whitelister/internal/pkg/config"
)

func (a *Azure) fetchResources(filter config.Filter) ([]resources.GenericResourceExpanded, error) {
	if filter.FilterType == config.SecurityGroup {
		return fetchSecurityGroupResources(a, filter)
	}
	logrus.Error("filter type " + filter.FilterType.String() + " not supported for azure yet")
	return nil, errors.New("filter type " + filter.FilterType.String() + " not supported for azure yet")
}

func (a *Azure) fetchSecurityGroup(securityGroupName string) (*network.SecurityGroup, error) {

	securityGroup, err := a.securityGroupClient.Get(context.TODO(), a.ResourceGroupName, securityGroupName, "")
	if err != nil {
		return nil, err
	}
	return &securityGroup, nil
}

func fetchSecurityGroupResources(a *Azure, filter config.Filter) (filteredResources []resources.GenericResourceExpanded, err error) {

	resources, err := a.resourcesClient.ListByResourceGroup(context.TODO(), a.ResourceGroupName, "resourceType eq 'Microsoft.Network/networkSecurityGroups'", "", nil)
	if err != nil {
		return nil, err
	}
	logrus.Infof("filter type %v", len(resources.Values()))
	
	for resources.NotDone() {
		filteredResources = append(filteredResources, filterResourcesByTag(resources.Values(), filter)...)
		if err := resources.NextWithContext(context.TODO()); err != nil {
			return nil, err
		}
	}
	return filteredResources, nil
}

func filterResourcesByTag(resources []resources.GenericResourceExpanded, filter config.Filter) (filteredResources []resources.GenericResourceExpanded) {
	for _, resource := range resources {
		if !hasMatchingTag(filter.LabelName, filter.LabelValue, resource.Tags) {
			continue
		}
		filteredResources = append(filteredResources, resource)
	}
	return filteredResources
}

func hasMatchingTag(tagName string, tagValue string, tagsMap map[string]*string) bool {
	for key, value := range tagsMap {
		if key == tagName && *value == tagValue {
			return true
		}
	}
	return false
}
