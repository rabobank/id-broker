package controllers

import (
	"fmt"
	"net/http"

	"github.com/gomatbase/go-we"
	"github.com/gomatbase/go-we/events"
	wutil "github.com/gomatbase/go-we/util"
	"github.com/rabobank/id-broker/managers"
)

func GetTokenInfo(writer we.ResponseWriter, r we.RequestScope) error {
	serviceInstanceId := r.Var("service_instance_guid")
	fmt.Println("GetTokenInfo", serviceInstanceId)

	info, e := managers.GetIdInfo(serviceInstanceId)
	if e != nil {
		fmt.Printf("failed to get info for service instance %s: %v\n", serviceInstanceId, e)
		return events.BadRequestError
	}
	return wutil.ReplyJson(writer, http.StatusOK, info)
}
