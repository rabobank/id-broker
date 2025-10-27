package cfg

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/gomatbase/csn"
)

type extract struct {
	name        string
	placeholder *string
}

func extractIntValue(variable *int, environmentVariable string) error {
	if value := os.Getenv(environmentVariable); value != "" {
		var err error
		*variable, err = strconv.Atoi(value)
		if err != nil {
			return csn.Error(fmt.Sprintf("failed reading envvar %s, err: %s", environmentVariable, err))
		}
	}

	return nil
}

func extractBoolValue(variable *bool, environmentVariable string) error {
	*variable = os.Getenv(environmentVariable) == "true"
	return nil
}

func extractStringValue(variable *string, environmentVariable string) error {
	value := os.Getenv(environmentVariable)
	if value != "" {
		*variable = value
	} else if *variable == "" {
		return csn.Error(fmt.Sprintf("environment variable %s is not set", environmentVariable))
	}
	return nil
}

// extractCredhubSecrets - Get the credentials from credhub (VCAP_SERVICES envvar) if available.
func extractCredhubSecrets(variables ...extract) {
	fmt.Println("getting credentials from credhub...")
	if appEnv, err := cfenv.Current(); err == nil {
		services, err := appEnv.Services.WithLabel("credhub")
		if err == nil {
			if len(services) != 1 {
				fmt.Printf("expected exactly one bound credhub service instance, but found %d\n", len(services))
				fmt.Println("Falling back to environment variables.")
			} else {
				for _, variable := range variables {
					if value, found := services[0].Credentials[variable.name]; found {
						*variable.placeholder = fmt.Sprint(value)
					} else {
						fmt.Printf("credhub variable %s is missing. Falling back to environment variables.\n", variable.name)
					}
				}
			}
		} else {
			fmt.Printf("failed getting services from cf env: %s\n", err)
			fmt.Println("credhub not available")
		}
	} else {
		fmt.Printf("failed to get the current cf env: %s\n", err)
		fmt.Println("credhub not available")
	}
}
