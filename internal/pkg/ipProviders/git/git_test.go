package git

import (
	"errors"
	"github.com/stakater/Whitelister/internal/pkg/utils"
	"reflect"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	testUtils "github.com/stakater/Whitelister/internal/pkg/test/utils"
)

var (
	configFilePath = "../../../../configs/testConfigs/"
	accessToken    = "access_token"
	url            = "https://github.com/"
	configFile     = "config.yaml"
	testFile       = "sampleConfig.yaml"
	emptyFile      = "Empty.yaml"
	ipCidr         = "127.0.0.1/32"
	description    = "Sample address"
	fromPort       = int64(80)
	toPort         = int64(80)
	ipProtocol     = "tcp"
)

func TestGitInit(t *testing.T) {

	tests := []struct {
		name     string
		args     map[interface{}]interface{}
		want     *Git
		wantErr  bool
		errValue error
	}{
		{
			name: "access token Only",
			args: map[interface{}]interface{}{
				"AccessToken": accessToken,
			},
			want:     &Git{AccessToken: accessToken},
			wantErr:  true,
			errValue: errors.New("Missing Git URL"),
		},
		{
			name: "URL only",
			args: map[interface{}]interface{}{
				"URL": url,
			},
			want:     &Git{URL: url},
			wantErr:  true,
			errValue: errors.New("Missing Git Access Token"),
		},
		{
			name: "url and access token Only",
			args: map[interface{}]interface{}{
				"AccessToken": accessToken,
				"URL":         url,
			},
			want:     &Git{AccessToken: accessToken, URL: url, Config: configFile},
			wantErr:  false,
			errValue: nil,
		},
		{
			name: "access token, url and config",
			args: map[interface{}]interface{}{
				"AccessToken": accessToken,
				"URL":         url,
				"Config":      configFile,
			},
			want:     &Git{AccessToken: accessToken, URL: url, Config: configFile},
			wantErr:  false,
			errValue: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Git{}
			err := got.Init(tt.args)

			if err != nil && tt.wantErr {
				if err.Error() != tt.errValue.Error() {
					t.Errorf("Got Err: %v, Wanted Err: %v", err, tt.errValue)
					return
				}
			}
			if !got.Equal(tt.want) {
				t.Errorf("Got = %v, wanted %v", got, tt.want)
			}
		})
	}
}

func TestReadConfig(t *testing.T) {

	result, err := testUtils.CopyFile(configFilePath+testFile, testFile, path)
	if !result && err != nil {
		t.Errorf("Cannot copy file. Error: %v", err)
	}

	result, err = testUtils.CopyFile(configFilePath+emptyFile, emptyFile, path)
	if !result && err != nil {
		t.Errorf("Cannot copy Emptyfile. Error: %v", err)
	}

	var ipRanges []*utils.IpRange
	ipRanges = append(ipRanges, &utils.IpRange{
		IpCidr:      &ipCidr,
		Description: &description,
	})

	ipPermissions := []utils.IpPermission{
		{
			IpRanges:   ipRanges,
			FromPort:   &fromPort,
			ToPort:     &toPort,
			IpProtocol: &ipProtocol,
		},
	}

	tests := []struct {
		name       string
		args       Git
		wantPerm   []utils.IpPermission
		wantConfig Config
		wantErr    bool
		errValue   error
	}{
		{
			name:     "Wrong Config Path",
			args:     Git{AccessToken: accessToken, Config: testFile + ".wrong", URL: url},
			wantErr:  true,
			errValue: errors.New("no such file or directory"),
		},
		{
			name:       "Empty Config",
			args:       Git{AccessToken: accessToken, Config: emptyFile, URL: url},
			wantConfig: Config{},
			wantErr:    false,
			errValue:   nil,
		},
		{
			name:     "Correct Config",
			args:     Git{AccessToken: accessToken, Config: testFile, URL: url},
			wantPerm: ipPermissions,
			wantErr:  false,
			errValue: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args
			config, err := got.readConfig()

			if err != nil && tt.wantErr {
				if !strings.Contains(err.Error(), tt.errValue.Error()) {
					t.Errorf("Got = %v, wanted %v", err, tt.errValue.Error())
				}
				return
			}

			if err == nil && !tt.wantErr {
				if !reflect.DeepEqual(config, tt.wantConfig) {
					for _, gotPermission := range config.IpPermissions {
						contains := false
						for _, wantPermission := range tt.wantPerm {
							if gotPermission.Equal(&wantPermission) {
								contains = true
							}
						}
						if !contains {
							t.Errorf("Mismatch")
							logrus.Println("Got:")
							testUtils.PrintIpPermissionsStructure(config.IpPermissions)
							logrus.Println("Wanted:")
							testUtils.PrintIpPermissionsStructure(tt.wantPerm)
							return
						}
					}
				}
			}
		})
	}
	err = testUtils.DeleteDir(path)
	if err != nil {
		t.Error(err.Error())
	}
}
