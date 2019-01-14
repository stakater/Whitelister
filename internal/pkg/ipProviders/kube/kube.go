package kube

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stakater/Whitelister/internal/pkg/utils"
	"github.com/stakater/Whitelister/pkg/kube"
)

// Kube Ip provider class implementing the IpProvider interface
type Kube struct {
	FromPort   *int64
	ToPort     *int64
	IpProtocol *string
}

func (k *Kube) GetName() string {
	return "Kubernetes"
}

// Init initializes the Kube Configuration like Tag name and value
func (k *Kube) Init(params map[interface{}]interface{}) error {
	err := mapstructure.Decode(params, &k) //Converts the params to kube struct fields
	if err != nil {
		return err
	}

	if k.FromPort == nil {
		return errors.New("Missing Kube From Port")
	}
	if k.ToPort == nil {
		return errors.New("Missing Kube To Port")
	}
	if k.IpProtocol == nil || *k.IpProtocol == "" {
		return errors.New("Missing Kube Ip Protocol")
	}
	return nil
}

// Get List of IP addresses to whitelist
func (k *Kube) GetIpPermissions() ([]utils.IpPermission, error) {
	clientset, err := kube.GetClient()
	if err != nil {
		logrus.Fatal(err)
	}

	nodes, err := clientset.Core().Nodes().List(meta_v1.ListOptions{})

	if err != nil {
		logrus.Fatal(err)
	}

	ipPermissions := []utils.IpPermission{}
	for _, node := range nodes.Items {
		for _, address := range node.Status.Addresses {
			if address.Type == "ExternalIP" {
				ipCidr := address.Address + "/32"
				ipPermissions = append(ipPermissions,
					utils.IpPermission{
						IpCidr:      &(ipCidr),
						FromPort:    k.FromPort,
						ToPort:      k.ToPort,
						IpProtocol:  k.IpProtocol,
						Description: (&(node.Name)),
					})
			}
		}
	}
	return ipPermissions, nil
}
