package ipProviders

import (
	"github.com/sirupsen/logrus"

	"github.com/stakater/Whitelister/internal/pkg/config"
	"github.com/stakater/Whitelister/internal/pkg/ipProviders/git"
	"github.com/stakater/Whitelister/internal/pkg/ipProviders/kube"
	"github.com/stakater/Whitelister/internal/pkg/utils"
)

// IpProvider interface so that other IpProvider like github can implement this
type IpProvider interface {
	Init(map[interface{}]interface{}) error
	GetIPPermissions() ([]utils.IpPermission, error)
	GetName() string
}

// PopulateFromConfig populates the IpProvider from config
func PopulateFromConfig(configIpProviders []config.IpProvider) []IpProvider {
	var populatedIpProviders []IpProvider
	for _, configIpProvider := range configIpProviders {
		ipProviderToAdd := MapToIpProvider(configIpProvider.Name)
		if ipProviderToAdd != nil {
			err := ipProviderToAdd.Init(configIpProvider.Params)
			if err != nil {
				logrus.Errorf("%v", err)
			} else {
				populatedIpProviders = append(populatedIpProviders, ipProviderToAdd)
			}
		}
	}
	return populatedIpProviders
}

// MapToIpProvider maps the IP provider name to the actual IpProvider type
func MapToIpProvider(ipProviderName string) IpProvider {
	switch ipProviderName {
	case "kubernetes":
		return &kube.Kube{}
	case "git":
		return &git.Git{}
	}
	logrus.Errorf("Cannot find an ip provider for : %s", ipProviderName)
	return nil
}
