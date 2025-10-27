package controllers

import (
	"fmt"
	"net/http"

	"github.com/gomatbase/go-we"
	wutil "github.com/gomatbase/go-we/util"
	"github.com/rabobank/id-broker/cfg"
	"github.com/rabobank/id-broker/domain"
	"github.com/rabobank/id-broker/managers"
)

func Catalog(w we.ResponseWriter, r we.RequestScope) error {
	fmt.Printf("get service broker catalog from %s...\n", r.Request().RemoteAddr)
	return wutil.DumpedReplyJson(w, http.StatusOK, cfg.Catalog)
}

func CreateServiceInstance(w we.ResponseWriter, r we.RequestScope) error {
	serviceInstanceId := r.Var("service_instance_guid")

	fmt.Printf("create service instance for %s...\n", serviceInstanceId)
	serviceInstance, e := wutil.ReadJsonBody[domain.ServiceInstance](r)
	if e != nil {
		return wutil.DumpedReplyJson(w, http.StatusBadRequest, domain.BrokerError{Error: "FAILED", Description: e.Error(), InstanceUsable: false, UpdateRepeatable: false})
	}

	if e = managers.CreateIdentity(serviceInstanceId, serviceInstance); e != nil {
		return wutil.DumpedReplyJson(w, http.StatusBadRequest, domain.BrokerError{Error: "FAILED", Description: fmt.Sprintf("service %s creation failed: %v", serviceInstance.ServiceId, e), InstanceUsable: false, UpdateRepeatable: false})
	} else {
		response := domain.CreateServiceInstanceResponse{
			ServiceId: serviceInstance.ServiceId,
			PlanId:    serviceInstance.PlanId,
			LastOperation: &domain.LastOperation{
				State: "succeeded",
			},
		}
		return wutil.DumpedReplyJson(w, http.StatusOK, response)
	}
}

func DeleteServiceInstance(w we.ResponseWriter, r we.RequestScope) error {
	serviceInstanceId := r.Var("service_instance_guid")
	fmt.Printf("delete service instance %s...\n", serviceInstanceId)

	if e := managers.DeleteIdentity(serviceInstanceId); e != nil {
		return wutil.DumpedReplyJson(w, http.StatusBadRequest, domain.BrokerError{Error: "FAILED", Description: fmt.Sprintf("service instance deletion with guid %s failed: %v", serviceInstanceId, e), InstanceUsable: false, UpdateRepeatable: false})
	}

	return wutil.DumpedReplyJson(w, http.StatusGone, domain.DeleteServiceInstanceResponse{Result: fmt.Sprintf("service instance with guid %s deleted", serviceInstanceId)})
}

func UpdateServiceInstance(w we.ResponseWriter, r we.RequestScope) error {
	serviceInstanceId := r.Var("service_instance_guid")
	fmt.Printf("update service instance %s...\n", serviceInstanceId)

	return wutil.DumpedReplyJson(w, http.StatusBadRequest, domain.UpdateServiceInstanceResponse{})
}
