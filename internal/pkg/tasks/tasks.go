package tasks

import (
	"context"

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
	combinedIPPermissions := []utils.IpPermission{}
	for _, ipProvider := range t.ipProviders {
		ipList, err := ipProvider.GetIPPermissions()
		if err != nil {
			logrus.Errorf("Error getting Ip list from provider: %s\n err: %v", ipProvider.GetName(), err)
		}
		combinedIPPermissions = utils.CombineIpPermission(combinedIPPermissions, ipList)
	}

	if t.config.Filter.FilterType.String() == config.LoadBalancerStr {
		loadBalancerNames := t.getLoadBalancerNames(t.config.Filter)
		logrus.Info("load balancer names: ", loadBalancerNames[0])

		if len(loadBalancerNames) > 0 {
			_ = t.provider.WhiteListIpsByLoadBalancer(loadBalancerNames, combinedIPPermissions)
		} else {
			logrus.Errorf("Cannot find any services with label name: " + t.config.Filter.LabelName +
				" , label value: " + t.config.Filter.LabelValue)
		}
	} else if t.config.Filter.FilterType.String() == config.SecurityGroupStr {
		filterLabel := []string{t.config.Filter.LabelName, t.config.Filter.LabelValue}
		_ = t.provider.WhiteListIpsBySecurityGroup(filterLabel, combinedIPPermissions)
	} else {
		logrus.Errorf("Unrecognized filter " + t.config.Filter.LabelName)
	}
}

// Get Load Balancer names
func (t *Task) getLoadBalancerNames(filter config.Filter) []string {
	services, err := t.clientset.CoreV1().Services("").List(context.TODO(), meta_v1.ListOptions{
		LabelSelector: filter.LabelName + "=" + filter.LabelValue},
	)

	if err != nil {
		logrus.Fatal(err)
	}

	var loadBalancerNames []string
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
