package config

import (
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	MongoConfigs struct {
		Hosts    []string      `config:"hosts;required"`
		Ports    []string      `config:"ports;required"`
		User     string        `config:"user;required"`
		Password string        `config:"password;required"`
		Timeout  time.Duration `config:"timeout;default=1s"`
	}

	MismatchHostsAndPortsError struct {
		HostCount int
		PortCount int
	}
)

func (e *MismatchHostsAndPortsError) Error() string {
	return fmt.Sprintf("hosts and ports quantity doesn't match! %d Hosts and %d Ports", e.HostCount, e.PortCount)
}

func (m *MongoConfigs) Url() (string, error) {
	hosts := len(m.Hosts)
	ports := len(m.Ports)
	if hosts != ports {
		return "", &MismatchHostsAndPortsError{hosts, ports}
	}

	hostsAndPorts := make([]string, 0, hosts)
	for index, host := range m.Hosts {
		hostsAndPorts = append(hostsAndPorts, fmt.Sprintf("%s:%s", host, m.Ports[index]))
	}

	return fmt.Sprintf("mongodb://%s", strings.Join(hostsAndPorts, ",")), nil
}

func (m *MongoConfigs) Credentials() options.Credential {
	return options.Credential{
		Username: m.User,
		Password: m.Password,
	}
}
