package utils

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stakater/Whitelister/internal/pkg/config"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
)

// Get Load Balancer names
func GetLoadBalancerNames(filter config.Filter, clientSet clientset.Interface) []string {
	services, err := clientSet.CoreV1().Services("").List(context.TODO(), meta_v1.ListOptions{
		LabelSelector: filter.LabelName + "=" + filter.LabelValue},
	)

	if err != nil {
		logrus.Fatal(err)
	}

	var loadBalancerNames []string
	var loadBalancerDNSName string

	for _, service := range services.Items {
		if service.Spec.Type == "LoadBalancer" {
			loadBalancerDNSName = GetLoadBalancerNameFromDNSName(service.Status.LoadBalancer.Ingress[0].Hostname)
			loadBalancerNames = append(loadBalancerNames, loadBalancerDNSName)
		} else {
			logrus.Error("Cannot process service : " + service.Name)
		}
	}
	return loadBalancerNames
}
