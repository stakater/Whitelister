package utils

import (
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
