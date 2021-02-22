package client

import "sync"

//
type ClientGroup struct {
	peers    clientCache
	orderers clientCache
}

func (c *ClientGroup) exec() {

}

type clientCache struct {
	stack sync.Map
}

func (c *clientCache) add(x *Client) {
	c.stack.Store(x.address, x)
}

func (c *clientCache) delete(address string) {
	c.stack.Delete(address)
}

func (c *clientCache) getSyncMap() *sync.Map {
	return &c.stack
}
