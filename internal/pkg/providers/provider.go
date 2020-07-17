package providers

import (
	"github.com/sirupsen/logrus"

	"github.com/stakater/Whitelister/internal/pkg/config"
	"github.com/stakater/Whitelister/internal/pkg/providers/aws"
	"github.com/stakater/Whitelister/internal/pkg/utils"
)

// Provider interface so that providers like aws, google cloud can implement this
type Provider interface {
	Init(map[interface{}]interface{}) error
	WhiteListIpsByLoadBalancer(resourceIds []string, ipPermissions []utils.IpPermission) error
	WhiteListIpsBySecurityGroup(filterLabel []string, ipPermissions []utils.IpPermission) error
}

// PopulateFromConfig populates the IpProvider from config
func PopulateFromConfig(configProvider config.Provider) Provider {
	providerToAdd := MapToProvider(configProvider.Name)
	if providerToAdd != nil {
		err := providerToAdd.Init(configProvider.Params)
		if err != nil {
			logrus.Errorf("%v", err)
		}
		return providerToAdd
	}
	return nil
}

// MapToIpProvider maps the IP provider name to the actual IpProvider type
func MapToProvider(providerName string) Provider {
	ipProvider, ok := providerMap[providerName]
	if !ok {
		logrus.Errorf("Cannot find an provider for : %s", providerName)
		return nil
	}
	return ipProvider
}

var providerMap = map[string]Provider{
	"aws": &aws.Aws{},
}
