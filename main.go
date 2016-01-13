package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/drone/drone-go/plugin"
	"github.com/rancher/go-rancher/client"
)

type Rancher struct {
	Url        string `json:"url"`
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	Service    string `json:"service"`
	Image      string `json:"docker_image"`
	StartFirst bool   `json:"start_first"`
	Confirm    bool   `json:"confirm"`
}

func main() {
	vargs := Rancher{StartFirst: true}

	plugin.Param("vargs", &vargs)
	err := plugin.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(vargs.Url) == 0 || len(vargs.AccessKey) == 0 || len(vargs.SecretKey) == 0 || len(vargs.Service) == 0 {
		return
	}

	if !strings.HasPrefix(vargs.Image, "docker:") {
		vargs.Image = fmt.Sprintf("docker:%s", vargs.Image)
	}

	var wantedService, wantedStack string
	if strings.Contains(vargs.Service, "/") {
		parts := strings.SplitN(vargs.Service, "/", 2)
		wantedStack = parts[0]
		wantedService = parts[1]
	} else {
		wantedService = vargs.Service
	}

	rancher, err := client.NewRancherClient(&client.ClientOpts{
		Url:       vargs.Url,
		AccessKey: vargs.AccessKey,
		SecretKey: vargs.SecretKey,
	})

	if err != nil {
		fmt.Printf("Failed to create rancher client: %s\n", err)
		os.Exit(1)
	}

	var stackId string
	if wantedStack != "" {
		environments, err := rancher.Environment.List(&client.ListOpts{})
		if err != nil {
			fmt.Printf("Failed to list rancher environments: %s\n", err)
			os.Exit(1)
		}

		for _, env := range environments.Data {
			if env.Name == wantedStack {
				stackId = env.Id
			}
		}

		if stackId == "" {
			fmt.Printf("Unable to find stack %s\n", wantedStack)
			os.Exit(1)
		}
	}

	services, err := rancher.Service.List(&client.ListOpts{})
	if err != nil {
		fmt.Printf("Failed to list rancher services: %s\n", err)
		os.Exit(1)
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
		fmt.Printf("Unable to find service %s\n", vargs.Service)
		os.Exit(1)
	}

	service.LaunchConfig.ImageUuid = vargs.Image
	upgrade := &client.ServiceUpgrade{}
	upgrade.InServiceStrategy = &client.InServiceUpgradeStrategy{
		LaunchConfig:           service.LaunchConfig,
		SecondaryLaunchConfigs: service.SecondaryLaunchConfigs,
		StartFirst:             vargs.StartFirst,
	}
	upgrade.ToServiceStrategy = &client.ToServiceUpgradeStrategy{}

	_, err = rancher.Service.ActionUpgrade(&service, upgrade)
	if err != nil {
		fmt.Printf("Unable to upgrade service %s\n", vargs.Service)
		os.Exit(1)
	}

	fmt.Printf("Upgraded %s to %s\n", vargs.Service, vargs.Image)

	if vargs.Confirm {
		err = retryFunc(30*time.Second, func() error {
			s, e := rancher.Service.ById(service.Id)
			if e != nil {
				return e
			}
			if s.State != "upgraded" {
				return fmt.Errorf("Service not upgraded: %s", s.State)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Error waiting for service upgrade to complete: %s", err)
			os.Exit(1)
		}

		s, err := rancher.Service.ById(service.Id)
		_, err = rancher.Service.ActionFinishupgrade(s)
		if err != nil {
			fmt.Printf("Unable to finish upgrade %s: %s\n", vargs.Service, err)
			os.Exit(1)
		}
		fmt.Printf("Finished upgrade %s\n", vargs.Service)
	}
}

func retryFunc(timeout time.Duration, f func() error) error {
	finish := time.After(timeout)
	for {
		err := f()
		if err == nil {
			return nil
		}
		log.Printf("Retryable error: %v", err)

		select {
		case <-finish:
			return err
		case <-time.After(3 * time.Second):
		}
	}
}
