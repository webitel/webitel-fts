package consul

import (
	"github.com/webitel/engine/discovery" // TODO
	"time"
)

const serviceTTL = time.Second * 30
const deregisterTTL = time.Second * 60

type Cluster struct {
	consulAddr string
	name       string
	discovery  discovery.ServiceDiscovery
}

func NewCluster(name string, consulAddr string) *Cluster {
	return &Cluster{
		name:       name,
		consulAddr: consulAddr,
	}
}

func (c *Cluster) Start(serviceId string, host string, port int) error {

	sd, err := discovery.NewServiceDiscovery(serviceId, c.consulAddr, func() (b bool, appError error) {
		return true, nil
	})
	if err != nil {
		return err
	}
	c.discovery = sd

	err = sd.RegisterService(c.name, host, port, serviceTTL, deregisterTTL)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cluster) Stop() {
	c.discovery.Shutdown()
}
