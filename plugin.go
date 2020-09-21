package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	client "github.com/rancher/go-rancher/v2"
)

type Plugin struct {
	URL                 string
	Key                 string
	Secret              string
	Service             string
	SidekickDockerImage []string
	DockerImage         string
	StartFirst          bool
	Confirm             bool
	Timeout             int
	IntervalMillis      int64
	BatchSize           int64
	YamlVerified        bool
	Environment         string
}

func (p *Plugin) Exec() error {
	log.Info("Drone Rancher Plugin built")

	if p.URL == "" || p.Key == "" || p.Secret == "" || p.Service == "" {
		return errors.New("Eek: Must have url, key, secret, and service definied")
	}

	var wantedService, wantedStack string
	if strings.Contains(p.Service, "/") {
		parts := strings.SplitN(p.Service, "/", 2)
		wantedStack = parts[0]
		wantedService = parts[1]
	} else {
		wantedService = p.Service
	}

	rancher, err := client.NewRancherClient(&client.ClientOpts{
		Url:       p.URL,
		AccessKey: p.Key,
		SecretKey: p.Secret,
	})
	if err != nil {
		return fmt.Errorf("Failed to create rancher client: %s", err)
	}

	// Prepare service filters for service listing
	serviceFilters := map[string]interface{}{"name": wantedService}

	if len(p.Environment) >= 1 {
		environments, err := rancher.Account.List(&client.ListOpts{Filters: map[string]interface{}{"name": p.Environment}})

		if err != nil {
			return fmt.Errorf("Failed to find given rancher environment: %s", err)
		}
		if len(environments.Data) <= 0 {
			return fmt.Errorf("Unable to find environment %s", p.Environment)
		}

		// If found add environmentID to serviceFilters
		serviceFilters["accountId"] = environments.Data[0].Id
	}

	// Query stacks with filter name=wantedStack
	if wantedStack != "" {
		stacks, err := rancher.Stack.List(&client.ListOpts{Filters: map[string]interface{}{"name": wantedStack}})

		// If environment is defined re-query the API with the accountId
		if len(p.Environment) >= 1 {
			stacks, err = rancher.Stack.List(&client.ListOpts{Filters: map[string]interface{}{"accountId": serviceFilters["accountId"], "name": wantedStack}})
		}

		if err != nil {
			return fmt.Errorf("Failed to list rancher environments: %s", err)
		}
		if len(stacks.Data) <= 0 {
			return fmt.Errorf("Unable to find stack %s", wantedStack)
		}
		// If found add stackID to serviceFilters
		serviceFilters["stackId"] = stacks.Data[0].Id

	}

	// Query services with prepared filters
	services, err := rancher.Service.List(&client.ListOpts{Filters: serviceFilters})
	if err != nil {
		return fmt.Errorf("Failed to list rancher services: %s", err)
	}
	if len(services.Data) <= 0 {
		return fmt.Errorf("Unable to find service %s", p.Service)
	}
	service := services.Data[0]
	// Service is found, proceed with upgrade

	// We want to exit if there is no docker image updates
	if p.DockerImage == "" && len(p.SidekickDockerImage) <= 0 {
		return fmt.Errorf("Nothing to upgrade")
	}

	// Only change value if it's not null.
	if p.DockerImage != "" {
		// Add prefix when missing to meet Rancher API requirement
		service.LaunchConfig.ImageUuid = prepareDockerPrefix(p.DockerImage)
	}

	// Iterate over provided sidekick flags
	for _, sidekick := range p.SidekickDockerImage {
		// Split flag in two from "--sidekick nginx nginx:latest"
		parts := strings.SplitN(sidekick, " ", 2)
		wantedSidekick := parts[0]
		wantedImage := parts[1]
		for i, s := range service.SecondaryLaunchConfigs {
			if wantedSidekick == s.Name {
				service.SecondaryLaunchConfigs[i].ImageUuid = prepareDockerPrefix(wantedImage)
			}
		}
	}
	upgrade := &client.ServiceUpgrade{}
	upgrade.InServiceStrategy = &client.InServiceUpgradeStrategy{
		LaunchConfig:           service.LaunchConfig,
		SecondaryLaunchConfigs: service.SecondaryLaunchConfigs,
		StartFirst:             p.StartFirst,
		IntervalMillis:         p.IntervalMillis,
		BatchSize:              p.BatchSize,
	}
	upgrade.ToServiceStrategy = &client.ToServiceUpgradeStrategy{}
	_, err = rancher.Service.ActionUpgrade(&service, upgrade)
	if err != nil {
		return fmt.Errorf("Unable to upgrade service %s: %s", p.Service, err)
	}

	if p.Confirm {
		srv, err := retry(func() (interface{}, error) {
			s, e := rancher.Service.ById(service.Id)
			if e != nil {
				return nil, e
			}
			if s.State != "upgraded" {
				return nil, fmt.Errorf("Service not upgraded: %s", s.State)
			}
			return s, nil
		}, time.Duration(p.Timeout)*time.Second, 3*time.Second)

		if err != nil {
			return fmt.Errorf("Error waiting for service upgrade to complete: %s", err)
		}

		_, err = rancher.Service.ActionFinishupgrade(srv.(*client.Service))
		if err != nil {
			return fmt.Errorf("Unable to finish upgrade %s: %s", p.Service, err)
		}
		log.Infof("Finished upgrade %s", p.Service)
	}

	log.Infof("Upgraded %s to %s", p.Service, p.DockerImage)

	return nil
}

type retryFunc func() (interface{}, error)

func prepareDockerPrefix(image string) string {
	if !strings.HasPrefix(image, "docker:") {
		image = fmt.Sprintf("docker:%s", image)
	}
	return image
}

func retry(f retryFunc, timeout time.Duration, interval time.Duration) (interface{}, error) {
	finish := time.After(timeout)
	for {
		result, err := f()
		if err == nil {
			return result, nil
		}
		select {
		case <-finish:
			return nil, err
		case <-time.After(interval):
		}
	}
}
