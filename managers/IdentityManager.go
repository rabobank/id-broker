package managers

import (
	"fmt"
	"log"

	"github.com/rabobank/credhub-client"
	"github.com/rabobank/id-broker/cfg"
	"github.com/rabobank/id-broker/domain"
	"github.com/rabobank/id-broker/uaa"
)

const (
	IdbIdentityCredhubName        = "/idb/%s"
	IdbIdentityBindingCredhubName = "/idb/%s/%s"
	IdbUrlFormat                  = "https://%s/token/%s/%s"
)

var credhubClient credhub.Client

func Initialize() {
	var e error
	if credhubClient, e = credhub.New(&credhub.Options{
		Client: cfg.ClientId,
		Secret: cfg.ClientSecret,
	}); e != nil {
		log.Fatalf("Unable to initialize credhub client: %v", e)
	}
}

func CreateIdentity(serviceInstanceId string, _ *domain.ServiceInstance) error {
	fmt.Println("Creating Identity for service instance", serviceInstanceId)
	idName := fmt.Sprintf(IdbIdentityCredhubName, serviceInstanceId)

	if _, e := credhubClient.GetByName(idName); e == nil {
		return fmt.Errorf("service instance id %s already exists", serviceInstanceId)
	}

	user, e := uaa.CreateUser(fmt.Sprintf("idb-%s", serviceInstanceId))
	if e != nil {
		return e
	}

	idInfo := domain.Info{
		Username: user.Username,
		Password: user.Password,
		Subject:  user.ID,
		Issuer:   cfg.CfTokenUrl,
		Audience: []string{"openid"},
	}

	if _, e = credhubClient.SetByName("json", idName, idInfo); e != nil {
		if _, e := uaa.DeleteUser(user.ID); e != nil {
			fmt.Printf("Unable to delete (rollback) user %s: %v\n", user.ID, e)
		}
		return e
	}

	return nil
}

func DeleteIdentity(serviceInstanceId string) error {
	fmt.Println("Deleting Identity for service instance", serviceInstanceId)
	idName := fmt.Sprintf(IdbIdentityCredhubName, serviceInstanceId)

	if credentials, e := credhubClient.GetJsonCredentialByName(idName); e != nil {
		fmt.Printf("deleting identity for service instance %s has failed: %v\n", serviceInstanceId, e)
		return e
	} else if _, e = uaa.DeleteUser(credentials.Value["subject"].(string)); e != nil {
		fmt.Printf("deleting identity for service instance %s has failed: %v\n", serviceInstanceId, e)
		return e
	}
	return nil
}

func BindIdentityToApp(serviceInstanceId string, bindingId string, appGuid string) (*domain.Credentials, error) {
	fmt.Println("Binding Identity", serviceInstanceId, "to app", appGuid)
	bindingName := fmt.Sprintf(IdbIdentityBindingCredhubName, serviceInstanceId, bindingId)

	if _, e := credhubClient.SetJsonByName(bindingName, map[string]any{"app-guid": appGuid}); e != nil {
		return nil, e
	}

	return &domain.Credentials{IdUrl: fmt.Sprintf(IdbUrlFormat, cfg.IdbUrl, serviceInstanceId, bindingId)}, nil
}

func UnbindIdentity(serviceInstanceId string, bindingId string) error {
	fmt.Println("Unbinding Identity", serviceInstanceId, "with id", bindingId)
	bindingName := fmt.Sprintf(IdbIdentityBindingCredhubName, serviceInstanceId, bindingId)

	if e := credhubClient.DeleteByName(bindingName); e != nil {
		return e
	}

	return nil
}

func ValidateBinding(serviceInstanceId string, bindingId string, appGuid string) error {
	fmt.Printf("Validating Binding %s of service instance %s with app %s", bindingId, serviceInstanceId, appGuid)
	bindingName := fmt.Sprintf(IdbIdentityBindingCredhubName, serviceInstanceId, bindingId)
	fmt.Println("Checking Binding", bindingName)
	credentials, e := credhubClient.GetJsonByName(bindingName)
	fmt.Println("result", credentials, e)
	if e != nil {
		fmt.Println("Error getting binding entry", e)
		return e
	}
	fmt.Println(credentials)
	if appGuid != credentials["app-guid"].(string) {
		return fmt.Errorf("app guid does not match binding guid")
	}
	return nil
}

func GenerateToken(serviceInstanceId string) (string, error) {
	fmt.Println("Generating token for service instance", serviceInstanceId)
	idName := fmt.Sprintf(IdbIdentityCredhubName, serviceInstanceId)

	if credentials, e := credhubClient.GetJsonByName(idName); e != nil {
		fmt.Println(e)
		return "", e
	} else {
		info := &domain.Info{}
		info.FromMap(credentials)

		token, e := uaa.Token(info.Username, info.Password, info.Audience)
		fmt.Println("result", token, e)
		if e != nil {
			return "", e
		}
		return token.AccessToken, nil
	}
}

func GetIdInfo(serviceInstanceId string) (*domain.Info, error) {
	fmt.Println("Getting token info for service instance", serviceInstanceId)
	idName := fmt.Sprintf(IdbIdentityCredhubName, serviceInstanceId)

	if credentials, e := credhubClient.GetJsonByName(idName); e != nil {
		return nil, e
	} else {
		info := &domain.Info{}
		info.FromMap(credentials)

		// let's clear the user credentials, if present
		info.Username = ""
		info.Password = ""
		return info, nil
	}
}
