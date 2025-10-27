package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rabobank/id-broker/cfg"
	"github.com/rabobank/id-broker/controllers"

	// "github.com/gomatbase/go-log"
	"github.com/gomatbase/go-we"
	"github.com/gomatbase/go-we/security"
)

func StartServer() {
	openIdProvider := security.OpenIdIdentityProvider(cfg.OpenIdUrl).
		Client(cfg.ClientId, cfg.ClientSecret).Scope("cloud_controller.read", "openid").
		Build()
	bearerAuthenticationProvider := security.BearerAuthenticationProvider().
		Introspector(openIdProvider.TokenIntrospector()).Build()
	basicAuthenticationProvider := security.BasicAuthenticationProvider(security.User{
		Username: cfg.BrokerUser,
		Password: cfg.BrokerPassword,
	}).Realm("id-broker - Cloud Foundry Identity Service Broker").Build()
	MtlsAuthenticationProvider := security.MtlsAuthenticationProvider().WithCertificateIdExtractor(instanceCertificateValidator).Build()

	allowedUsers := security.Either(security.Scope("cloud_controller.admin"), security.AuthorizationFunc(isDeveloper))

	securityFilter := security.Filter(security.DefaultAuthenticatedAccess).Authentication(basicAuthenticationProvider).
		Path("/health").Anonymous().
		Path("/token/**").Authentication(MtlsAuthenticationProvider).Authorize(security.AuthorizationFunc(checkBoundApp)).
		Path("/api/**").Authentication(bearerAuthenticationProvider).Authorize(allowedUsers).
		Build()

	// the go-router presents only the leaf certificate. The instanceIdentityCA is regenerated with every cf update.
	// Without the intermediate certificate, if using reversed proxy mtls, there's no reliable way to verify the certificate
	// itself, we can only validate that it is present, meaning that mTls was used.

	engine := we.New()
	engine.AddFilter(we.FilterFunction(DebugFilter))
	engine.AddFilter(security.ReverseProxyMtlsFilter().IgnoringVerificationErrors().Build())
	engine.AddFilter(securityFilter)

	// open broker api
	engine.HandleMethod(http.MethodGet, "/v2/catalog", controllers.Catalog)
	engine.HandleMethod(http.MethodPut, "/v2/service_instances/{service_instance_guid}", controllers.CreateServiceInstance)
	engine.HandleMethod(http.MethodPatch, "/v2/service_instances/{service_instance_guid}", controllers.UpdateServiceInstance)
	engine.HandleMethod(http.MethodDelete, "/v2/service_instances/{service_instance_guid}", controllers.DeleteServiceInstance)
	engine.HandleMethod(http.MethodPut, "/v2/service_instances/{service_instance_guid}/service_bindings/{service_binding_guid}", controllers.CreateServiceBinding)
	engine.HandleMethod(http.MethodDelete, "/v2/service_instances/{service_instance_guid}/service_bindings/{service_binding_guid}", controllers.DeleteServiceBinding)

	// token api
	engine.HandleMethod(http.MethodGet, "/token/{service_instance_guid}/{service_binding_guid}", controllers.GetToken)

	// support api
	engine.HandleMethod(http.MethodGet, "/api/info/{service_instance_guid}", controllers.GetTokenInfo)

	log.Fatal(engine.Listen(fmt.Sprintf(":%d", cfg.ListenPort)))
}
