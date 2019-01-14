package controller

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/stakater/Whitelister/internal/pkg/config"
	"github.com/stakater/Whitelister/internal/pkg/ipProviders"
	"github.com/stakater/Whitelister/internal/pkg/providers"
	"github.com/stakater/Whitelister/internal/pkg/tasks"
	clientset "k8s.io/client-go/kubernetes"
)

// Controller Jamadar Controller to check for left over items
type Controller struct {
	clientset   clientset.Interface
	config      config.Config
	ipProviders []ipProviders.IpProvider
	provider    providers.Provider
}

// NewController for initializing the Controller
func NewController(clientset clientset.Interface, config config.Config) (*Controller, error) {
	controller := &Controller{
		clientset: clientset,
		config:    config,
	}

	controller.ipProviders = ipProviders.PopulateFromConfig(config.IpProviders)
	if len(controller.ipProviders) == 0 {
		return nil, errors.New("No Ip Provider specified")
	}
	controller.provider = providers.PopulateFromConfig(config.Provider)
	if controller.provider == nil {
		return nil, errors.New("No Provider specified")
	}
	return controller, nil
}

//Run function for controller which handles the logic
func (c *Controller) Run() {
	for {
		c.handleTasks()
		timeInterval := c.config.SyncInterval
		duration, err := time.ParseDuration(timeInterval)
		if err != nil {
			logrus.Infof("Error Parsing Time Interval: %v", err)
			return
		}
		time.Sleep(duration)
	}
}

func (c *Controller) handleTasks() {
	task := tasks.NewTask(c.clientset, c.ipProviders, c.provider, c.config)
	task.PerformTasks()
}
