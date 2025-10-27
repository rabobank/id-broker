package controllers

import (
	"fmt"
	"net/http"

	"github.com/gomatbase/go-we"
	"github.com/gomatbase/go-we/events"
	wutil "github.com/gomatbase/go-we/util"
	"github.com/rabobank/id-broker/managers"
)

func GetToken(writer we.ResponseWriter, r we.RequestScope) error {
	serviceInstanceId := r.Var("service_instance_guid")
	serviceBindingId := r.Var("service_binding_guid")
	fmt.Println("GetToken", serviceInstanceId, serviceBindingId)

	token, e := managers.GenerateToken(serviceInstanceId)
	if e != nil {
		fmt.Printf("failed to get token for instance %s: %v\n", serviceInstanceId, e)
		return events.ForbiddenError
	}
	return wutil.ReplyJson(writer, http.StatusOK, map[string]string{"accessToken": token})
}
