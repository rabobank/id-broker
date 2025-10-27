package controllers

import (
	"net/http"

	"github.com/gomatbase/go-we"
	wutil "github.com/gomatbase/go-we/util"
	"github.com/rabobank/id-broker/domain"
	"github.com/rabobank/id-broker/managers"
)

func CreateServiceBinding(w we.ResponseWriter, r we.RequestScope) error {
	serviceInstanceId := r.Var("service_instance_guid")
	serviceBindingId := r.Var("service_binding_guid")

	serviceBinding, e := wutil.ReadJsonBody[domain.ServiceBinding](r)
	if e != nil {
		return wutil.DumpedReplyJson(w, http.StatusBadRequest, domain.BrokerError{Error: "FAILED", Description: e.Error(), InstanceUsable: false, UpdateRepeatable: false})
	}

	if credentials, e := managers.BindIdentityToApp(serviceInstanceId, serviceBindingId, serviceBinding.AppGuid); e == nil {
		response := domain.CreateServiceBindingResponse{Credentials: credentials}
		return wutil.DumpedReplyJson(w, http.StatusCreated, response)
	} else {
		return wutil.DumpedReplyJson(w, http.StatusBadRequest, domain.BrokerError{Error: "FAILED", Description: e.Error(), InstanceUsable: false, UpdateRepeatable: false})
	}

}

func DeleteServiceBinding(w we.ResponseWriter, r we.RequestScope) error {
	serviceInstanceId := r.Var("service_instance_guid")
	serviceBindingId := r.Var("service_binding_guid")

	if e := managers.UnbindIdentity(serviceInstanceId, serviceBindingId); e == nil {
		response := domain.DeleteServiceBindingResponse{Result: "unbind completed"}
		return wutil.DumpedReplyJson(w, http.StatusOK, response)
	} else {
		return wutil.DumpedReplyJson(w, http.StatusBadRequest, domain.BrokerError{Error: "FAILED", Description: e.Error(), InstanceUsable: false, UpdateRepeatable: false})
	}
}
