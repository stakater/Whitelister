package kube

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"

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

// GetIPPermissions - Get List of IP addresses to whitelist
func (k *Kube) GetIPPermissions() ([]utils.IpPermission, error) {
	client, err := kube.GetClient()
	if err != nil {
		return nil, err
	}
	return k.getNodesIPPermissions(client.CoreV1())
}

func (k *Kube) getNodesIPPermissions(client v1.CoreV1Interface) ([]utils.IpPermission, error) {

	nodes, err := client.Nodes().List(meta_v1.ListOptions{})

	if err != nil {
		return nil, err
	}

	ipPermissions := []utils.IpPermission{}
	for _, node := range nodes.Items {
		ipPermission, err := k.getNodeIPPermissions(node)
		if err != nil {
			logrus.Error(err)
		} else {
			ipPermissions = append(ipPermissions, *ipPermission)
		}
	}
	return ipPermissions, nil
}

// getNodeIPPermissions - Get IP permission based on ExternalIP of node
func (k *Kube) getNodeIPPermissions(node corev1.Node) (*utils.IpPermission, error) {
	for _, address := range node.Status.Addresses {
		if address.Type == "ExternalIP" {
			ipCidr := address.Address + "/32"
			return &utils.IpPermission{
				IpCidr:      &(ipCidr),
				FromPort:    k.FromPort,
				ToPort:      k.ToPort,
				IpProtocol:  k.IpProtocol,
				Description: (&(node.Name)),
			}, nil
		}
	}
	return nil, fmt.Errorf("No ExternalIP for Node: %s", node.Name)
}
