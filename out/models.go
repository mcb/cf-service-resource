package out

import "github.com/mcb/cf-service-resource"

type Request struct {
	Source resource.Source `json:"source"`
	Params Params          `json:"params"`
}

type Params struct {
	Service          string    `json:"service"`
	Plan             string    `json:"plan"`
	InstanceName     string    `json:"instance_name"`
	ManifestPath     string    `json:"manifest"`
	CurrentAppName   string    `json:"current_app_name"`
	SkipBinding      bool    	 `json:"skip_binding"`
	SkipRestage      bool      `json:"skip_restage"`
}

type Response struct {
	Version  resource.Version        `json:"version"`
	Metadata []resource.MetadataPair `json:"metadata"`
}
