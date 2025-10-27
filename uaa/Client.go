package uaa

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudfoundry-community/go-uaa"
	"github.com/cloudfoundry-community/go-uaa/passwordcredentials"
	"github.com/rabobank/id-broker/cfg"
	"github.com/rabobank/id-broker/util"
	"golang.org/x/oauth2"
)

var uaac *uaa.API

func Initialize() {
	fmt.Println("Initializing UAA client")
	var e error
	if uaac, e = uaa.New(cfg.CfAuthUrl, uaa.WithClientCredentials(cfg.ClientId, cfg.ClientSecret, uaa.JSONWebToken), uaa.WithSkipSSLValidation(true)); e != nil {
		fmt.Printf("Failed to authenticate with UAA: %v\n", e)
		os.Exit(8)
	}
}

func Token(user, password string, audience []string) (*oauth2.Token, error) {
	uaaCredentials := &passwordcredentials.Config{
		ClientID:     cfg.ClientId,
		ClientSecret: cfg.ClientSecret,
		Username:     user,
		Password:     password,
		Endpoint:     oauth2.Endpoint{TokenURL: cfg.CfTokenUrl},
		Scopes:       audience,
	}
	return uaaCredentials.TokenSource(context.Background()).Token()
}

func CreateUser(username string) (*uaa.User, error) {
	password := util.GeneratePassword()

	TRUE := true
	uaaUser := uaa.User{
		Password: password,
		Username: username,
		// test if indeed name and emails are required
		Name: &uaa.UserName{
			FamilyName: username,
			GivenName:  username,
		},
		Emails: []uaa.Email{
			{
				Value:   username,
				Primary: &TRUE,
			},
		},
		Origin: "",
	}

	user, e := uaac.CreateUser(uaaUser)
	if e == nil {
		user.Password = password
	}
	return user, e
}

func DeleteUser(id string) (*uaa.User, error) {
	return uaac.DeleteUser(id)
}
