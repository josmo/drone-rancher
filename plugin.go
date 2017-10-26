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
	URL            string
	Key            string
	Secret         string
	Service        string
	DockerImage    string
	StartFirst     bool
	Confirm        bool
	Timeout        int
	IntervalMillis int64
	BatchSize      int64
	YamlVerified   bool
}

func (p *Plugin) Exec() error {
	log.Info("Drone Rancher Plugin built")

	if p.URL == "" || p.Key == "" || p.Secret == "" {
		return errors.New("Eek: Must have url, key, secret, and service definied")
	}

	if !strings.HasPrefix(p.DockerImage, "docker:") {
		p.DockerImage = fmt.Sprintf("docker:%s", p.DockerImage)
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

	var stackID string
	if wantedStack != "" {
		environments, err := rancher.Stack.List(&client.ListOpts{})
		if err != nil {
			return fmt.Errorf("Failed to list rancher environments: %s", err)
		}
		for _, env := range environments.Data {
			if env.Name == wantedStack {
				stackID = env.Id
			}
		}
		if stackID == "" {
			return fmt.Errorf("Unable to find stack %s", wantedStack)
		}
	}

	// Get the initial set of services
	services, err := rancher.Service.List(&client.ListOpts{})
	// Check for an error
	if err != nil {
		return fmt.Errorf("Failed to list rancher services: %s", err)
	}

	var service client.Service
	found := false
	//TODO: find a more elegant and clear way to iterate through the service paging
	for {

		// Iterate the current services
		for _, svc := range services.Data {
			if svc.Name == wantedService && ((wantedStack != "" && svc.StackId == stackID) || wantedStack == "") {
				service = svc
				found = true
				break
			}
		}
		if found {
			break
		}

		// Get the next set of services (paginate)
		if !found {
			services, err = services.Next()
			if err != nil {
				return fmt.Errorf("Failed to list rancher services: %s", err)
			}
			if services == nil {
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("Unable to find service %s", p.Service)
	}

	service.LaunchConfig.ImageUuid = p.DockerImage
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

	log.Infof("Upgraded %s to %s", p.Service, p.DockerImage)
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
	return nil
}

type retryFunc func() (interface{}, error)

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
