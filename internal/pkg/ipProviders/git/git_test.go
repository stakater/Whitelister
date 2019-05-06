package git

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stakater/Whitelister/internal/pkg/test/utils"
)

var (
	configFilePath = "../../../../configs/testConfigs/"
)

func TestGitInit(t *testing.T) {

	var accessToken string
	accessToken = "ABC"

	var url string
	url = "https://example.com/"

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
			want:     &Git{AccessToken: accessToken, URL: url, Config: "config.yaml"},
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
			if !(got.URL == tt.want.URL && got.AccessToken == tt.want.AccessToken) {
				t.Errorf("Got = %v, wanted %v", got, tt.want)
			}
		})
	}
}

func TestReadConfig(t *testing.T) {

	testFile := "sampleConfig.yaml"
	EmptyFile := "Empty.yaml"

	result, err := utils.CopyFile(configFilePath+testFile, "/tmp/whitelister-config/"+testFile)
	if !result && err != nil {
		t.Errorf("Cannot copy file. Error: %v", err)
	}

	result1, err1 := utils.CopyFile(configFilePath+EmptyFile, "/tmp/whitelister-config/"+EmptyFile)
	if !result1 && err1 != nil {
		t.Errorf("Cannot copy Emptyfile. Error: %v", err1)
	}

	tests := []struct {
		name     string
		args     Git
		wantErr  bool
		errValue error
	}{
		{
			name:     "Config Path",
			args:     Git{AccessToken: "ABC", Config: EmptyFile + ".wrong", URL: "https://example.com"},
			wantErr:  true,
			errValue: errors.New("no such file or directory"),
		},
		{
			name:     "Empty Config",
			args:     Git{AccessToken: "ABC", Config: EmptyFile, URL: "https://example.com"},
			wantErr:  false,
			errValue: nil,
		},
		{
			name:     "Correct Config",
			args:     Git{AccessToken: "ABC", Config: testFile, URL: "https://example.com"},
			wantErr:  false,
			errValue: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args
			config, err := got.readConfig()

			fmt.Println(config)

			if err != nil && tt.wantErr {
				if !strings.Contains(err.Error(), tt.errValue.Error()) {
					t.Errorf("Got = %v, wanted %v", err, tt.errValue.Error())
				}
				return
			}

			if err == nil && tt.wantErr {
				t.Errorf("Wanted = %v, but got no error", err)
				return
			}

			if err == nil && !tt.wantErr {
				if reflect.DeepEqual(config, Config{}) {
					fmt.Println("Empty Config File.")
				}
				for _, permission := range config.IpPermissions {
					if permission.FromPort == nil {
						t.Errorf("permission.FromPort not set")
					}
					if permission.ToPort == nil {
						t.Errorf("permission.ToPort not set")
					}
					if permission.IpProtocol == nil {
						t.Errorf("permission.IpProtocol not set")
					}
					if permission.IpRanges == nil {
						t.Errorf("permission.IpRanges not set")
					}
				}
			}
		})
	}
}
