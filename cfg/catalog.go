package cfg

import (
	"github.com/rabobank/id-broker/domain"
)

const (
	IdbApiUrlAttribute = "idb_api"
)

var Catalog = domain.Catalog{
	Services: []domain.Service{
		{
			Name:        "cf.identity",
			Id:          "ae3b8714-d9d7-4570-915f-1deecd97b610",
			Description: "Cf Identity instance",
			Bindable:    true,
			Metadata: map[string]any{
				"shareable":        true,
				"longDescription":  "Service broker that creates cf (uaa) users for which tokens can be requested to access external resources when uaa is setup as backend IDP",
				"displayName":      "cf.identity",
				"documentationUrl": "https://confluence.dev.rabobank.nl/display/IT4IT/Cloud+Foundry+Identity",
			},
			MaxPollInterval: 7200,
			PlanUpdateable:  false,
			Plans: []domain.ServicePlan{
				{
					Name:        "standard",
					Id:          "2b21351e-b99a-4f5c-b79b-188487c8acb9",
					Description: "Standard CF ID",
					Metadata: map[string]any{
						"cost":             0,
						IdbApiUrlAttribute: "",
					},
				},
			},
		},
	},
}
