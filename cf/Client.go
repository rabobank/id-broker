package cf

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/cloudfoundry/go-cfclient/v3/config"
	"github.com/rabobank/id-broker/cfg"
	"github.com/rabobank/id-broker/domain"
	"github.com/rabobank/id-broker/http"
)

var (
	cfClient *client.Client
	cfCtx    = context.Background()
)

func Initialize() {
	fmt.Println("Initializing CF")
	if cfConfig, e := config.New(cfg.CfApiUrl, config.ClientCredentials(cfg.ClientId, cfg.ClientSecret), config.UserAgent(fmt.Sprintf("id-broker/%s", cfg.BuildTime))); e != nil {
		log.Fatalf("failed to create new cf config: %s", e)
	} else if cfClient, e = client.New(cfConfig); e != nil {
		log.Fatalf("failed to create new cf client: %s", e)
	}
}

func ValidateApp(appGuid string, ipAddress string) error {
	// fmt.Printf("Validating App %s on %s\n", appGuid, ipAddress)
	app, e := cfClient.Applications.Get(cfCtx, appGuid)
	if e != nil {
		return fmt.Errorf("error getting app %s: %s", appGuid, e)
	}
	if strings.ToLower(app.State) != "started" {
		return fmt.Errorf("app %s is not running", appGuid)
	}

	processes, e := cfClient.Processes.ListForAppAll(cfCtx, appGuid, nil)
	if e != nil {
		return fmt.Errorf("error getting app %s processes: %s", appGuid, e)
	}

	for _, process := range processes {
		statResources := &domain.StatResources{}
		if e = http.Request(cfg.CfApiUrl, "/v3/processes", process.GUID, "stats").WithClient(cfClient.HTTPAuthClient()).GetJson(statResources); e != nil {
			return fmt.Errorf("error getting process %s info: %s", process.GUID, e)
		}

		for _, stats := range statResources.Resources {
			if stats.InstanceInternalIp == ipAddress {
				// fmt.Println("App %s instance found running on %s", appGuid, ipAddress)
				return nil
			}
		}
	}
	return fmt.Errorf("no app %s processes found running on ip %s", appGuid, ipAddress)
}
