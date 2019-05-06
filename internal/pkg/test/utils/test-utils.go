package utils

import (
	"fmt"
	"io/ioutil"

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

func CopyFile(sourceFile string, destinationFile string) (bool, error) {
	result := true
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		fmt.Println("Error creating", destinationFile)
		fmt.Println(err)
		return false, err
	}
	return result, nil
}
