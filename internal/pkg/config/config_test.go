package config

import (
	"errors"
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
		name     string
		args     args
		want     Config
		wantErr  bool
		errValue error
	}{
		{
			name: "TestingWithCorrectValues",
			args: args{filePath: configFilePath + "correctAwsKubernetesConfig.yaml"},
			want: Config{
				SyncInterval: "10s",
				IpProviders: []IpProvider{
					{
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
						"RoleArn":                   "arn:aws:iam::111111111111:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling",
						"Region":                    "us-west-2",
						"RemoveRule":                true,
						"KeepRuleDescriptionPrefix": "DO NOT REMOVE -",
					},
				},
				Filter: Filter{
					FilterType: LoadBalancer,
					LabelName:  "whitelister",
					LabelValue: "true",
				},
			},
			wantErr: false,
		},
		{
			name: "TestingWithCorrectValuesForSecurityGroupFilter",
			args: args{filePath: configFilePath + "correctAwsGitConfigWithSG.yaml"},
			want: Config{
				SyncInterval: "10s",
				IpProviders: []IpProvider{
					{
						Name: "git",
						Params: map[interface{}]interface{}{
							"AccessToken": "access-token",
							"URL":         "http://github.com/stakater/whitelister-config.git",
							"Config":      "config.yaml",
						},
					},
				},
				Provider: Provider{
					Name: "aws",
					Params: map[interface{}]interface{}{
						"RoleArn":                   "arn:aws:iam::111111111111:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling",
						"Region":                    "us-west-2",
						"RemoveRule":                true,
						"KeepRuleDescriptionPrefix": "DO NOT REMOVE -",
					},
				},
				Filter: Filter{
					FilterType: SecurityGroup,
					LabelName:  "whitelister",
					LabelValue: "true",
				},
			},
			wantErr: false,
		},
		{
			name: "TestingWithCorrectValuesForSecurityGroupFilterOnAzure",
			args: args{filePath: configFilePath + "correctAzureGitConfigWithSG.yaml"},
			want: Config{
				SyncInterval: "10s",
				IpProviders: []IpProvider{
					{
						Name: "git",
						Params: map[interface{}]interface{}{
							"AccessToken": "access-token",
							"URL":         "http://github.com/stakater/whitelister-config.git",
							"Config":      "config.yaml",
						},
					},
				},
				Provider: Provider{
					Name: "azure",
					Params: map[interface{}]interface{}{
						"RemoveRule":                true,
						"KeepRuleDescriptionPrefix": "DO NOT REMOVE -",
						"SubscriptionID":            "47c9180a-967d-4ba0-bfc0-7b12762f0779",
						"ClientID":                  "4ab0b7f2-197f-4b14-bf22-1856a6f095aa",
						"ClientSecret":              "thisisthesecret",
						"TenantID":                  "73cf1f9c-03d0-4709-8434-b50bc8440454",
						"ResourceGroupName":         "my-resource-group",
					},
				},
				Filter: Filter{
					FilterType: SecurityGroup,
					LabelName:  "whitelister",
					LabelValue: "true",
				},
			},
			wantErr: false,
		},
		{
			name:     "TestingWithIncorrectFilterType",
			args:     args{filePath: configFilePath + "configWithIncorrectFilterType.yaml"},
			wantErr:  true,
			errValue: errors.New("incorrect FilterType :InCorrectType provided"),
		},
		{
			name:     "TestingWithIncorrectFilterTypeAzure",
			args:     args{filePath: configFilePath + "configWithIncorrectFilterTypeAzure.yaml"},
			wantErr:  true,
			errValue: errors.New("incorrect FilterType :InCorrectType provided"),
		},
		{
			name:    "TestingWithEmptyFile",
			args:    args{filePath: configFilePath + "Empty.yaml"},
			want:    Config{},
			wantErr: false,
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
			if (err == nil && tt.wantErr) || (!tt.wantErr && err != nil) {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr && tt.errValue != nil {
				if err.Error() != tt.errValue.Error() {
					t.Errorf("ReadConfig() error %v, wantErr %v", err, tt.errValue)
					return
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
