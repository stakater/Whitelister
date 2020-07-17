package kube

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"

	testUtils "github.com/stakater/Whitelister/internal/pkg/test/utils"
	"github.com/stakater/Whitelister/internal/pkg/utils"
)

func TestKubeInit(t *testing.T) {
	var fromPort int64
	fromPort = int64(0)

	var toPort int64
	toPort = int64(0)

	emptyIPProtocol := ""
	tcpIPProtocol := "tcp"

	tests := []struct {
		name     string
		args     map[interface{}]interface{}
		want     *Kube
		wantErr  bool
		errValue error
	}{
		{
			name:     "Missing config",
			args:     nil,
			want:     &Kube{},
			wantErr:  true,
			errValue: errors.New("missing kube params in config"),
		},
		{
			name:     "Empty Config",
			args:     map[interface{}]interface{}{},
			want:     &Kube{},
			wantErr:  true,
			errValue: errors.New("Missing Kube From Port"),
		},
		{
			name: "From Port Only",
			args: map[interface{}]interface{}{
				"FromPort": 0,
			},
			want:     &Kube{FromPort: &fromPort},
			wantErr:  true,
			errValue: errors.New("Missing Kube To Port"),
		},
		{
			name: "From and To Port Only",
			args: map[interface{}]interface{}{
				"FromPort": fromPort,
				"ToPort":   toPort,
			},
			want:     &Kube{FromPort: &fromPort, ToPort: &toPort},
			wantErr:  true,
			errValue: errors.New("Missing Kube Ip Protocol"),
		},
		{
			name: "From Port, To Port and Empty Ip Protocol",
			args: map[interface{}]interface{}{
				"FromPort":   fromPort,
				"ToPort":     toPort,
				"IpProtocol": emptyIPProtocol,
			},
			want:     &Kube{FromPort: &fromPort, ToPort: &toPort, IpProtocol: &emptyIPProtocol},
			wantErr:  true,
			errValue: errors.New("Missing Kube Ip Protocol"),
		},
		{
			name: "From Port, To Port and Tcp Ip Protocol",
			args: map[interface{}]interface{}{
				"FromPort":   fromPort,
				"ToPort":     toPort,
				"IpProtocol": tcpIPProtocol,
			},
			want:     &Kube{FromPort: &fromPort, ToPort: &toPort, IpProtocol: &tcpIPProtocol},
			wantErr:  false,
			errValue: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Kube{}
			err := got.Init(tt.args)

			if err != nil && tt.wantErr {
				if err.Error() != tt.errValue.Error() {
					t.Errorf("Got Err: %v, Wanted Err: %v", err, tt.errValue)
					return
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Got = %v, wanted %v", got, tt.want)
			}
		})
	}
}

func TestGetNodesIPPermissions(t *testing.T) {

	kube := Kube{}
	kube.Init(map[interface{}]interface{}{"FromPort": int64(0), "ToPort": int64(65535), "IpProtocol": "tcp"})

	var nodesList []runtime.Object
	var ipRanges []*utils.IpRange

	for i := 1; i < 3; i++ {
		ipCidr := fmt.Sprintf("127.0.0.%d/32", i)
		ipAddr := fmt.Sprintf("127.0.0.%d", i)
		name := fmt.Sprintf("node-%d", i)

		nodesList = append(nodesList, testUtils.Node(name, ipAddr))

		ipRanges = append(ipRanges, &utils.IpRange{
			IpCidr:      &ipCidr,
			Description: &name,
		})
	}

	ipPermissions := []utils.IpPermission{
		{
			FromPort:   kube.FromPort,
			ToPort:     kube.ToPort,
			IpProtocol: kube.IpProtocol,
			IpRanges:   ipRanges,
		},
	}

	tests := []struct {
		name     string
		args     []runtime.Object
		want     []utils.IpPermission
		wantErr  bool
		errValue error
	}{
		{
			name: "get 0 nodes IP",
			args: []runtime.Object{},
			want: []utils.IpPermission{
				{
					FromPort:   kube.FromPort,
					ToPort:     kube.ToPort,
					IpProtocol: kube.IpProtocol,
				},
			},
			wantErr:  false,
			errValue: nil,
		},
		{
			name:     "get 3 nodes IP",
			args:     nodesList,
			want:     ipPermissions,
			wantErr:  false,
			errValue: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := fake.NewSimpleClientset(tt.args...)

			got, err := kube.getNodesIPPermissions(client.CoreV1())

			if err != nil && tt.wantErr {
				if err.Error() != tt.errValue.Error() {
					t.Errorf("Got Err: %v, Wanted Err: %v", err, tt.errValue)
					return
				}
			}

			for _, gotPermission := range got {
				contains := false
				for _, wantPermission := range tt.want {
					if gotPermission.Equal(&wantPermission) {
						contains = true
					}
				}
				if !contains {
					t.Errorf("Got: %v, Wanted: %v", got, tt.want)
					return
				}
			}
		})
	}
}

func TestGetNodeIPPermissions(t *testing.T) {

	kube := Kube{}
	kube.Init(map[interface{}]interface{}{"FromPort": int64(0), "ToPort": int64(65535), "IpProtocol": "tcp"})
	ipAddr := "127.0.0.1"
	ipCidr := fmt.Sprintf("%s/32", ipAddr)
	name := "name"

	tests := []struct {
		name     string
		args     coreV1.Node
		want     *utils.IpRange
		wantErr  bool
		errValue error
	}{
		{
			name:     "Node without External IP",
			args:     *testUtils.Node(name, ""),
			wantErr:  true,
			errValue: fmt.Errorf("No ExternalIP for Node: %s", name),
		},
		{
			name: "Node with External IP",
			args: *testUtils.Node(name, ipAddr),
			want: &utils.IpRange{
				IpCidr:      &ipCidr,
				Description: &name,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := kube.getNodeIPRange(tt.args)

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
