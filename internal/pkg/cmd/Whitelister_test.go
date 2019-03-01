package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestGetWhiteListerCommand(t *testing.T) {

	tests := []struct {
		name    string
		want    *cobra.Command
		wantErr bool
	}{
		{
			name: "Get Cobra Command",
			want: &cobra.Command{
				Use:   "Whitelister",
				Short: "A tool which manages AWS security groups to allow access to nodes and developers",
				Run:   startWhitelister,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWhitelisterCommand()

			if got.Use != tt.want.Use || got.Short != tt.want.Short || got.Run == nil {
				t.Errorf("NewWhitelisterCommand() = %v, \n want = %v", got, tt.want)
			}
		})
	}
}
