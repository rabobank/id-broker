package domain

type ServiceInstance struct {
	ServiceId        string   `json:"service_id"`
	PlanId           string   `json:"plan_id"`
	OrganizationGuid string   `json:"organization_guid"`
	SpaceGuid        string   `json:"space_guid"`
	Context          *Context `json:"context"`
	Parameters       any      `json:"parameters,omitempty"`
}

type CreateServiceInstanceResponse struct {
	ServiceId     string         `json:"service_id"`
	PlanId        string         `json:"plan_id"`
	DashboardUrl  string         `json:"dashboard_url"`
	LastOperation *LastOperation `json:"last_operation,omitempty"`
}

type UpdateServiceInstanceResponse struct {
	DashBoardUrl string `json:"dashboard_url"`
	Operation    string `json:"operation"`
}

type ServiceInstanceResponse struct {
	ServiceId    string `json:"service_id"`
	PlanId       string `json:"plan_id"`
	DashboardUrl string `json:"dashboard_url"`
	Parameters   any    `json:"parameters"`
}

type DeleteServiceInstanceResponse struct {
	Result string `json:"result,omitempty"`
}
