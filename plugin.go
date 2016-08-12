package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/rancher/go-rancher/client"
)

type (
	Rancher struct {
		Url        string `json:"url"`
		AccessKey  string `json:"access_key"`
		SecretKey  string `json:"secret_key"`
		Service    string `json:"service"`
		Image      string `json:"docker_image"`
		StartFirst bool   `json:"start_first"`
		Confirm    bool   `json:"confirm"`
		Timeout    int64  `json:"timeout"`
		Params     string `json:"params"`
	}

	ErrorMessage struct {
		message string
	}

	Plugin struct {
		Rancher      Rancher
		ErrorMessage ErrorMessage
	}
)

func (p ErrorMessage) Error() string {
	return fmt.Sprintf("Error: %s", p.message)
}

func (p Plugin) Exec() error {
	if len(p.Rancher.Url) == 0 || len(p.Rancher.AccessKey) == 0 || len(p.Rancher.SecretKey) == 0 || len(p.Rancher.Service) == 0 {
		return ErrorMessage{"Please provide the following arguments: Url, AccessKey, SecretKey and Service"}
	}

	if !strings.HasPrefix(p.Rancher.Image, "docker:") {
		p.Rancher.Image = fmt.Sprintf("docker:%s", p.Rancher.Image)
	}

	var wantedService, wantedStack string
	if strings.Contains(p.Rancher.Service, "/") {
		parts := strings.SplitN(p.Rancher.Service, "/", 2)
		wantedStack = parts[0]
		wantedService = parts[1]
	} else {
		wantedService = p.Rancher.Service
	}

	rancher, err := client.NewRancherClient(&client.ClientOpts{
		Url:       p.Rancher.Url,
		AccessKey: p.Rancher.AccessKey,
		SecretKey: p.Rancher.SecretKey,
	})

	if err != nil {
		return ErrorMessage{fmt.Sprintf("Failed to create rancher client: %s\n", err)}
	}

	// Search for stack
	var stackId string
	if wantedStack != "" {
		environments, err := rancher.Environment.List(&client.ListOpts{})
		if err != nil {
			return ErrorMessage{fmt.Sprintf("Failed to list rancher environments: %s\n", err)}
		}

		for _, env := range environments.Data {
			if env.Name == wantedStack {
				stackId = env.Id
			}
		}

		if stackId == "" {
			return ErrorMessage{fmt.Sprintf("Unable to find stack %s\n", wantedStack)}
		}
	}

	// Search for service
	services, err := rancher.Service.List(&client.ListOpts{})
	if err != nil {
		return ErrorMessage{fmt.Sprintf("Failed to list rancher services: %s\n", err)}
	}

	found := false
	var service client.Service
	for _, svc := range services.Data {
		if svc.Name == wantedService && ((wantedStack != "" && svc.EnvironmentId == stackId) || wantedStack == "") {
			service = svc
			found = true
		}
	}

	if !found {
		return ErrorMessage{fmt.Sprintf("Unable to find service %s\n", p.Rancher.Service)}
	}

	service.LaunchConfig.ImageUuid = p.Rancher.Image
	upgrade := &client.ServiceUpgrade{}
	upgrade.InServiceStrategy = &client.InServiceUpgradeStrategy{
		LaunchConfig:           service.LaunchConfig,
		SecondaryLaunchConfigs: service.SecondaryLaunchConfigs,
		StartFirst:             p.Rancher.StartFirst,
	}
	upgrade.ToServiceStrategy = &client.ToServiceUpgradeStrategy{}

	_, err = rancher.Service.ActionUpgrade(&service, upgrade)
	if err != nil {
		return ErrorMessage{fmt.Sprintf("Unable to upgrade service %s: %s\n", p.Rancher.Service, err)}
	}

	fmt.Printf("Upgraded %s to %s\n", p.Rancher.Service, p.Rancher.Image)

	if p.Rancher.Confirm {
		srv, err := retry(func() (interface{}, error) {
			s, e := rancher.Service.ById(service.Id)
			if e != nil {
				return nil, e
			}
			if s.State != "upgraded" {
				return nil, ErrorMessage{fmt.Sprintf("Service not upgraded: %s\n", s.State)}
			}
			return s, nil
		}, time.Duration(p.Rancher.Timeout)*time.Second, 3*time.Second)

		if err != nil {
			return ErrorMessage{fmt.Sprintf("Error waiting for service upgrade to complete: %s\n", err)}
		}

		_, err = rancher.Service.ActionFinishupgrade(srv.(*client.Service))
		if err != nil {
			return ErrorMessage{fmt.Sprintf("Unable to finish upgrade %s: %s\n", p.Rancher.Service, err)}
		}
		fmt.Printf("Finished upgrade %s\n", p.Rancher.Service)
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
