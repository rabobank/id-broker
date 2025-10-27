package cfg

type Link struct {
	Href string            `json:"href"`
	Meta map[string]string `json:"meta,omitempty"`
}

type CfInfo struct {
	Links struct {
		Self              *Link `json:"self,omitempty"`
		CloudControllerV3 *Link `json:"cloud_controller_v3,omitempty"`
		NetworkPolicyV0   *Link `json:"network_policy_v0,omitempty"`
		NetworkPolicyV1   *Link `json:"network_policy_v1,omitempty"`
		Login             *Link `json:"login,omitempty"`
		Uaa               *Link `json:"uaa,omitempty"`
		Credhub           *Link `json:"credhub,omitempty"`
		Routing           *Link `json:"routing,omitempty"`
		Logging           *Link `json:"logging,omitempty"`
		LogCache          *Link `json:"log_cache,omitempty"`
		LogStream         *Link `json:"log_stream,omitempty"`
		AppSsh            *Link `json:"app_ssh,omitempty"`
		CloudControllerV2 *Link `json:"cloud_controller_v2,omitempty"`
	} `json:"links"`
}
