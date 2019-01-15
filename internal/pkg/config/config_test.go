package config

import (
	"reflect"
	"testing"
)

var (
	configFilePath = "../../../configs/testConfigs/"
)

func TestReadConfig(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name: "TestingWithCorrectValues",
			args: args{filePath: configFilePath + "correctAwsKubernetesConfig.yaml"},
			want: Config{
				SyncInterval:    "10s",
				RemoveUnknownIp: true,
				IpProviders: []IpProvider{
					IpProvider{
						Name: "kubernetes",
						Params: map[interface{}]interface{}{
							"FromPort":   0,
							"ToPort":     65535,
							"IpProtocol": "tcp",
						},
					},
				},
				Provider: Provider{
					Name: "aws",
					Params: map[interface{}]interface{}{
						"RoleArn": "arn:aws:iam::111111111111:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling",
						"Region":  "us-west-2",
					},
				},
				Filter: Filter{
					LabelName:  "whitelister",
					LabelValue: "true",
				},
			},
		},
		{
			name: "TestingWithEmptyFile",
			args: args{filePath: configFilePath + "Empty.yaml"},
			want: Config{},
		},
		{
			name:    "TestingWithFileNotPresent",
			args:    args{filePath: configFilePath + "FileNotFound.yaml"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadConfig(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
