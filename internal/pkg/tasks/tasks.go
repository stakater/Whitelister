package tasks

import (
	"github.com/sirupsen/logrus"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	combinedIpPermissions := []utils.IpPermission{}
	for _, ipProvider := range t.ipProviders {
		ipList, err := ipProvider.GetIPPermissions()
		if err != nil {
			logrus.Errorf("Error getting Ip list from provider: %s\n err: %v", ipProvider.GetName(), err)
		}
		combinedIpPermissions = append(combinedIpPermissions, ipList...)
	}

	loadBalancerNames := t.getLoadBalancerNames(t.config.Filter)

	if len(loadBalancerNames) > 0 {
		t.provider.WhiteListIps(loadBalancerNames, combinedIpPermissions)
		logrus.Infof("%v", loadBalancerNames)
	} else {
		logrus.Errorf("Cannot find any services with label name: " + t.config.Filter.LabelName +
			" , label value: " + t.config.Filter.LabelValue)
	}

}

// Get Load Balancer names
func (t *Task) getLoadBalancerNames(filter config.Filter) []string {
	services, err := t.clientset.Core().Services("").List(meta_v1.ListOptions{
		LabelSelector: filter.LabelName + "=" + filter.LabelValue},
	)

	if err != nil {
		logrus.Fatal(err)
	}

	loadBalancerNames := []string{}
	var loadBalancerDNSName string

	for _, service := range services.Items {
		if service.Spec.Type == "LoadBalancer" {
			loadBalancerDNSName =
				utils.GetLoadBalancerNameFromDNSName(service.Status.LoadBalancer.Ingress[0].Hostname)
			loadBalancerNames = append(loadBalancerNames, loadBalancerDNSName)
		} else {
			logrus.Error("Cannot process service : " + service.Name)
		}
	}
	return loadBalancerNames
}
