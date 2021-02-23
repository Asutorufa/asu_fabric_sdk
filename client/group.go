package client

import (
	"context"
	"log"
	"sync"

	"github.com/hyperledger/fabric-protos-go/peer"
)

//Group peer and orderer group
type Group struct {
	peers    clientCache
	orderers clientCache
}

func (g *Group) exec() {

}

//EndorserProposal endorse proposal
func (g *Group) EndorserProposal(endorserAddress []string, sp *peer.SignedProposal) *peer.ProposalResponse {
	endorserMap := make(map[string]interface{})

	for ea := range endorserAddress {
		endorserMap[endorserAddress[ea]] = nil
	}

	var proposalResponse *peer.ProposalResponse

	g.peers.getSyncMap().Range(func(key, value interface{}) bool {
		keyS, ok := key.(string)
		if !ok {
			return false
		}

		if _, ok := endorserMap[keyS]; !ok {
			return false
		}

		vc, ok := value.(*Client)
		if !ok {
			return false
		}

		p := &PeerClient{
			Client: *vc,
		}

		endorser, err := p.Endorser()
		if err != nil {
			log.Printf("get endorser failed: %v\n", err)
			return false
		}

		proposalResponse, err = endorser.ProcessProposal(context.Background(), sp)
		if err != nil {
			log.Printf("endorser process proposal failed: %v\n", err)
			return false
		}

		return false
	})

	return proposalResponse
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
