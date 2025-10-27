package server

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gomatbase/go-we"
	"github.com/gomatbase/go-we/events"
	"github.com/gomatbase/go-we/security"
	"github.com/rabobank/id-broker/cf"
	"github.com/rabobank/id-broker/cfg"
	"github.com/rabobank/id-broker/http"
	"github.com/rabobank/id-broker/managers"
)

const (
	PermissionsUrl = "%s/v3/service_instances/%s/permissions"
)

type CfServiceInstancePermissions struct {
	Read   bool `json:"read"`
	Manage bool `json:"manage"`
}

func instanceCertificateValidator(certificate *x509.Certificate) (*security.User, error) {
	// extract app guid
	var appGuid string
	for _, ou := range certificate.Subject.OrganizationalUnit {
		if strings.HasPrefix(ou, "app:") {
			appGuid = strings.TrimPrefix(ou, "app:")
			break
		}
	}
	if appGuid == "" {
		return nil, events.ForbiddenError
	}

	user := &security.User{
		Username: appGuid,
	}

	if len(certificate.IPAddresses) != 1 {
		return nil, events.ForbiddenError
	}

	user.Data = certificate.IPAddresses[0].String()

	if e := cf.ValidateApp(appGuid, user.Data.(string)); e != nil {
		return nil, events.ForbiddenError
	}

	user.Active = true

	return user, nil
}

func checkBoundApp(user *security.User, r we.RequestScope) bool {
	// fmt.Println("Checking BoundApp", user, r)
	if user == nil {
		return false
	}

	serviceInstanceId := r.Var("service_instance_guid")
	serviceBindingId := r.Var("service_binding_guid")
	fmt.Println("Checking BoundApp", serviceInstanceId, serviceBindingId)
	if e := managers.ValidateBinding(serviceInstanceId, serviceBindingId, user.Username); e != nil {
		return false
	}

	return true
}

func isDeveloper(user *security.User, r we.RequestScope) bool {
	serviceInstanceId := r.Var("service_instance_guid")
	if user == nil {
		return false
	}

	bearerToken := r.Request().Header.Get("Authorization")
	if len(bearerToken) == 0 {
		// call is not authenticated with a bearer token, the user should have the token in the metadata
		if token, isTokenData := user.Data.(*security.TokenData); !isTokenData {
			// can't get a token to validate
			return false
		} else {
			bearerToken = "Bearer " + token.Raw
		}
	}

	body, e := http.Request(fmt.Sprintf(PermissionsUrl, cfg.CfApiUrl, serviceInstanceId)).WithAuthorization(bearerToken).Get()
	if e != nil {
		// log it
	} else {
		var permissions CfServiceInstancePermissions
		if e = json.Unmarshal(body, &permissions); e != nil {
			fmt.Printf("Unable to check user %s permissions for service %s: %v\n", user.Username, serviceInstanceId, e)
		} else if permissions.Manage {
			return true
		} else {
			fmt.Printf("[AUTH] User %s has no permissions to manage requested service %s\n", user.Username, serviceInstanceId)
		}
	}
	return false
}
