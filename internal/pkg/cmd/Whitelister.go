package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/stakater/Whitelister/internal/pkg/config"
	"github.com/stakater/Whitelister/internal/pkg/controller"
	"github.com/stakater/Whitelister/pkg/kube"
)

//NewWhitelisterCommand to start and run Whitelister
func NewWhitelisterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Whitelister",
		Short: "A tool which manages AWS security groups to allow access to nodes and developers",
		Run:   startWhitelister,
	}
	return cmd
}

func startWhitelister(cmd *cobra.Command, args []string) {
	logrus.Infof("Starting Whitelister")
	// create the clientset
	clientset, err := kube.GetClient()
	if err != nil {
		logrus.Panic(err)
	}
	if clientset != nil {
		logrus.Infof("GOT CLIENT SET")
	} else {
		logrus.Panicf("Kube Client set not found.")
	}

	// get the Controller config file
	config := config.GetConfiguration()

	controller, err := controller.NewController(clientset, config)
	if err != nil {
		logrus.Errorf("Error occured while creating controller. Reason: %s", err.Error())
		return
	}

	go controller.Run()

	// Wait forever
	select {}
}
