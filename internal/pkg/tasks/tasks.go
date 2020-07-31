package tasks

import (
	"github.com/sirupsen/logrus"
	clientset "k8s.io/client-go/kubernetes"

	"github.com/stakater/Whitelister/internal/pkg/config"
	"github.com/stakater/Whitelister/internal/pkg/ipProviders"
	"github.com/stakater/Whitelister/internal/pkg/providers"
	"github.com/stakater/Whitelister/internal/pkg/utils"
)

// Task represents the actual tasks and actions to be taken by Whitelister
type Task struct {
	clientset   clientset.Interface
	ipProviders []ipProviders.IpProvider
	provider    providers.Provider
	config      config.Config
}

// NewTask creates a new Task object
func NewTask(clientSet clientset.Interface, ipProviders []ipProviders.IpProvider,
	provider providers.Provider, conf config.Config) *Task {
	return &Task{
		clientset:   clientSet,
		ipProviders: ipProviders,
		provider:    provider,
		config:      conf,
	}
}

// PerformTasks handles all tasks
func (t *Task) PerformTasks() {
	combinedIPPermissions := []utils.IpPermission{}
	for _, ipProvider := range t.ipProviders {
		ipList, err := ipProvider.GetIPPermissions()
		if err != nil {
			logrus.Errorf("Error getting Ip list from provider: %s\n err: %v", ipProvider.GetName(), err)
		}
		combinedIPPermissions = utils.CombineIpPermission(combinedIPPermissions, ipList)
	}

	_ = t.provider.WhiteListIps(t.config.Filter, combinedIPPermissions)
}
