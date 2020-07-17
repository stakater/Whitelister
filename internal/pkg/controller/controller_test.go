package controller

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	testClient "k8s.io/client-go/kubernetes/fake"

	"github.com/stakater/Whitelister/internal/pkg/config"
	testUtils "github.com/stakater/Whitelister/internal/pkg/test/utils"
)

var (
	configFilePath = "../../../configs/testConfigs/"
)

func TestController(t *testing.T) {

	noIPProviderConfig, _ := config.ReadConfig(configFilePath + "noIPProviderConfig.yaml")
	noProviderConfig, _ := config.ReadConfig(configFilePath + "noProviderConfig.yaml")

	var nodes []v1.Node

	//create 3 nodes
	for i := 0; i < 3; i++ {
		nodeName := fmt.Sprintf("%s%d", "node", i)
		nodes = append(nodes, *testUtils.Node(nodeName, fmt.Sprintf("127.0.0.%d", i)))
	}
	// create the clientset
	clientset := testClient.NewSimpleClientset(&v1.NodeList{Items: nodes})

	tests := []struct {
		name     string
		args     config.Config
		want     Controller
		wantErr  bool
		errValue error
	}{
		{
			name:     "Get New Controller without Ip provider",
			args:     noIPProviderConfig,
			wantErr:  true,
			errValue: errors.New("No Ip Provider specified"),
		},
		{
			name:     "Get New Controller without provider",
			args:     noProviderConfig,
			wantErr:  true,
			errValue: errors.New("No Provider specified"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewController(clientset, tt.args)

			if err != nil && tt.wantErr {
				if err.Error() != tt.errValue.Error() {
					t.Errorf("Got Err: %v, Wanted Err: %v", err, tt.errValue)
					return
				}
				if got != nil {
					t.Errorf("Got: %v, Wanted: nil", got)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Got: %v, Wanted: %v", got, tt.want)
			}
		})
	}
}
