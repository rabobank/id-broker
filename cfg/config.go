package cfg

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/gomatbase/csn"
	"github.com/rabobank/id-broker/http"
)

// default values
const (
	DefaultHttpTimeout = 10
	DefaultListenPort  = 8080

	DefaultCredhubUrl = "https://credhub.service.cf.internal:8844"
)

// Variable Names
const (
	DebugVariable       = "DEBUG"
	UsernameVariable    = "USERNAME"
	PasswordVariable    = "PASSWORD"
	ListenPortVariable  = "LISTEN_PORT"
	HttpTimeoutVariable = "HTTP_TIMEOUT"
	CredhubUrlVariable  = "CREDHUB_URL"

	UaaUrlVariable          = "UAA_URL"
	UaaClientIdVariable     = "UAA_CLIENT_ID"
	UaaClientSecretVariable = "UAA_CLIENT_SECRET"
)

// Variables to identify the build
var (
	Version    string
	BuildTime  string
	CommitHash string
)

var (
	ListenPort = DefaultListenPort
	Debug      = false

	BrokerUser     string
	BrokerPassword string

	HttpTimeout = DefaultHttpTimeout

	CfApiUrl   string
	CfAuthUrl  string
	CfTokenUrl string

	OpenIdUrl    string
	ClientId     string
	ClientSecret string
	CredhubURL   = DefaultCredhubUrl

	IdbUrl string
)

func Initialize() {
	errors := csn.Errors()

	cfEnv, e := cfenv.Current()
	if e != nil {
		errors.Add(e)
	} else {
		CfApiUrl = cfEnv.CFAPI
		// Extract the uaa auth endpoint
		var cfInfo CfInfo
		if e = http.Request(cfEnv.CFAPI).GetJson(&cfInfo); e != nil {
			errors.AddErrorMessage(fmt.Sprintf("unable to extract CF info : %v", e))
		} else if cfInfo.Links.Uaa == nil || len(cfInfo.Links.Uaa.Href) == 0 {
			errors.AddErrorMessage("Cf Info holds no UAA URL")
		} else {
			CfAuthUrl = cfInfo.Links.Uaa.Href
			OpenIdUrl = cfInfo.Links.Uaa.Href
			if CfTokenUrl, e = url.JoinPath(cfInfo.Links.Uaa.Href, "oauth", "token"); e != nil {
				errors.AddErrorMessage(fmt.Sprintf("invalid token url : %v", e.Error()))
			}
		}
		if len(cfEnv.ApplicationURIs) == 0 {
			errors.AddErrorMessage("Cf Info holds no application URIs")
		} else if len(cfEnv.ApplicationURIs) > 1 {
			fmt.Println("Cf Info holds multiple application URIs, using the first one", cfEnv.ApplicationURIs[0])
		}

		IdbUrl = cfEnv.ApplicationURIs[0]
		Catalog.Services[0].Plans[0].Metadata.(map[string]any)[IdbApiUrlAttribute] = fmt.Sprintf("https://%s", IdbUrl)
	}

	if v, present := os.LookupEnv("VCAP_PLATFORM_OPTIONS"); present {
		var options map[string]any
		var ok bool
		if e := json.Unmarshal([]byte(v), &options); e != nil {
			fmt.Println("Unable to unmarshal VCAP_PLATFORM_OPTIONS:", e)
		} else if v, present := options["credhub-uri"]; !present {
			fmt.Println("VCAP_PLATFORM_OPTIONS does not have credhub-uri value")
		} else if CredhubURL, ok = v.(string); !ok {
			fmt.Println("VCAP_PLATFORM_OPTIONS credhub-uri value is not a string")
			CredhubURL = DefaultCredhubUrl
		}
	}
	errors.Add(extractStringValue(&CredhubURL, CredhubUrlVariable))

	// Try to get variables from credhub if available, and then falling back or overwriting from the environment
	extractCredhubSecrets(
		extract{PasswordVariable, &BrokerPassword},
		extract{UaaClientSecretVariable, &ClientSecret},
	)

	// Set Debug flag
	errors.Add(extractBoolValue(&Debug, DebugVariable))

	// Configure broker
	errors.Add(extractStringValue(&BrokerUser, UsernameVariable))
	errors.Add(extractStringValue(&BrokerPassword, PasswordVariable))
	errors.Add(extractIntValue(&ListenPort, ListenPortVariable))
	errors.Add(extractIntValue(&HttpTimeout, HttpTimeoutVariable))

	// configure oauth2 client
	errors.Add(extractStringValue(&OpenIdUrl, UaaUrlVariable))
	errors.Add(extractStringValue(&ClientId, UaaClientIdVariable))
	errors.Add(extractStringValue(&ClientSecret, UaaClientSecretVariable))

	if errors.Count() > 0 {
		fmt.Println("id-broker not fully configured....")
		fmt.Println(errors)
		os.Exit(8)
	}

	fmt.Println("Debug :", Debug)
	fmt.Println("Cf API :", cfEnv.CFAPI)
	fmt.Println("OpenIdUrl :", OpenIdUrl)
	fmt.Println("Cf Auth Url :", CfAuthUrl)
	fmt.Println("Cf Token Url :", CfTokenUrl)
	fmt.Println("ClientId :", ClientId)
	fmt.Println("ClientSecret :", len(ClientSecret) != 0)
	fmt.Println("Credhub URL :", CredhubURL)
	fmt.Println("Listening on port :", ListenPort)
	fmt.Println("Http request timeout :", HttpTimeout)

}
