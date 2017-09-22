package main

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/client"
	"strings"
	"time"
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
		//need to check for other tags
		//if !strings.HasSuffix(p.DockerImage, "docker:") {
			p.DockerImage = fmt.Sprintf("%s:%s", p.DockerImage, p.Tags)
		//}
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
		return errors.New(fmt.Sprintf("Failed to create rancher client: %s\n :(", err))
	}

	var stackId string
	if wantedStack != "" {
		environments, err := rancher.Environment.List(&client.ListOpts{})
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to list rancher environments: %s\n", err))
		}
		for _, env := range environments.Data {
			if env.Name == wantedStack {
				stackId = env.Id
			}
		}
		if stackId == "" {
			return errors.New(fmt.Sprintf("Unable to find stack %s\n", wantedStack))
		}
	}

	// Get the initial set of services
	services, err := rancher.Service.List(&client.ListOpts{})
	// Check for an error
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to list rancher services: %s\n", err))
	}

	var service client.Service
	found := false
	//TODO: find a more elegant and clear way to iterate through the service paging
	for {

		// Iterate the current services
		for _, svc := range services.Data {
			if svc.Name == wantedService && ((wantedStack != "" && svc.EnvironmentId == stackId) || wantedStack == "") {
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
				return errors.New(fmt.Sprintf("Failed to list rancher services: %s\n", err))
			}
			if services == nil {
				break
			}
		}
	}

	if !found {
		return errors.New(fmt.Sprintf("Unable to find service %s\n", p.Service))
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
		return errors.New(fmt.Sprintf("Unable to upgrade service %s: %s\n", p.Service, err))
	}

	log.Info(fmt.Sprintf("Upgraded %s to %s\n", p.Service, p.DockerImage))
	if p.Confirm {
		srv, err := retry(func() (interface{}, error) {
			s, e := rancher.Service.ById(service.Id)
			if e != nil {
				return nil, e
			}
			if s.State != "upgraded" {
				return nil, errors.New(fmt.Sprintf("Service not upgraded: %s", s.State))
			}
			return s, nil
		}, time.Duration(p.Timeout)*time.Second, 3*time.Second)

		if err != nil {
			return errors.New(fmt.Sprintf("Error waiting for service upgrade to complete: %s", err))
		}

		_, err = rancher.Service.ActionFinishupgrade(srv.(*client.Service))
		if err != nil {
			return errors.New(fmt.Sprintf("Unable to finish upgrade %s: %s\n", p.Service, err))
		}
		log.Info(fmt.Printf("Finished upgrade %s\n", p.Service))
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
