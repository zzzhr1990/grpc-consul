package server

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

// ConsulRegisterConfig unregister service from consul
type ConsulRegisterConfig struct {
	ConsulAddress string
	ServerID      string
	Name          string
	ServicePort   int
	ServiceIP     string
	EnableCheck   bool
}

type ConsulResult struct {
	Agent    *api.Agent
	ServerID string
}

func NewConsulResult(agent *api.Agent, serverID string) *ConsulResult {
	return &ConsulResult{
		Agent:    agent,
		ServerID: serverID,
	}
}

func (c *ConsulResult) ShutdownAgent() error {
	//return c.Agent
	return c.Agent.ServiceDeregister(c.ServerID)
}

// RegistToConsul register service to consul
func RegisterToConsul(registerConfig *ConsulRegisterConfig) (*ConsulResult, error) {
	// Create a new consul client
	cfg := api.DefaultConfig()
	cfg.Address = registerConfig.ConsulAddress
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create consul client failed")
	}

	// Create a new agent
	agent := client.Agent()
	// agent.ServiceDeregister()

	// Register service with consul
	reg := &api.AgentServiceRegistration{
		ID:      registerConfig.ServerID,
		Name:    registerConfig.Name,
		Port:    registerConfig.ServicePort,
		Address: registerConfig.ServiceIP,
	}

	if registerConfig.EnableCheck {
		reg.Check = &api.AgentServiceCheck{
			// TTL:                            "15s",
			Interval:                       "10s",
			GRPC:                           fmt.Sprintf("%v:%v", registerConfig.ServiceIP, registerConfig.ServicePort),
			DeregisterCriticalServiceAfter: "30s",
		}
	}

	err = agent.ServiceRegister(reg)
	if err != nil {
		return nil, errors.Wrap(err, "register service failed")
	}
	return NewConsulResult(agent, registerConfig.ServerID), nil
}
