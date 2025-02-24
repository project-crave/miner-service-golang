package configuration

import (
	craveConfiguration "crave/shared/configuration"
)

type Dependency struct {
	HubGrpc *craveConfiguration.Api
}

type Variable struct {
	Dependency *Dependency
	Neo4jMiner *craveConfiguration.Database
	Api        *craveConfiguration.Api
	GrpcApi    *craveConfiguration.Api
}

func NewVariable() *Variable {
	return &Variable{
		Neo4jMiner: &craveConfiguration.Database{
			Uri:      "neo4j://127.0.0.1:7686",
			Username: "neo4j",
			Password: "your_password",
		},
		Dependency: &Dependency{
			HubGrpc: &craveConfiguration.Api{
				Ip:   "localhost",
				Port: 3002,
			},
		},
		Api: &craveConfiguration.Api{
			Ip:   "localhost",
			Port: 3000,
		},
		GrpcApi: &craveConfiguration.Api{
			Ip:   "localhost",
			Port: 3001,
		},
	}
}
