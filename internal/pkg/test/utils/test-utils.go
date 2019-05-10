package utils

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/stakater/Whitelister/internal/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Node - creates a new node object with specified name(required) and ip address (optional)
func Node(name string, ipAddress string) *corev1.Node {
	addresses := []corev1.NodeAddress{}

	if ipAddress != "" {
		addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeExternalIP, Address: ipAddress})
	}

	return &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Status:     corev1.NodeStatus{Addresses: addresses},
	}
}

func CopyFile(sourceFile string, destinationFile string, destinationDir string) (bool, error) {
	result := true

	if _, err := os.Stat(destinationDir); os.IsNotExist(err) {
		os.Mkdir(destinationDir, 0777)
	}

	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	err = ioutil.WriteFile(destinationDir+"/"+destinationFile, input, 0777)
	if err != nil {
		fmt.Println("Error creating", destinationFile)
		fmt.Println(err)
		return false, err
	}
	return result, nil
}

func DeleteDir(path string) error {
	return os.RemoveAll(path)
}

func PrintIpPermissionsStructure(ipPerms []utils.IpPermission) {
	for _, ipPerm := range ipPerms {
		for _, iprange := range ipPerm.IpRanges {
			logrus.Println(*iprange.IpCidr)
			logrus.Println(*iprange.Description)
		}
		logrus.Println(*ipPerm.ToPort)
		logrus.Println(*ipPerm.FromPort)
		logrus.Println(*ipPerm.IpProtocol)
	}
}
