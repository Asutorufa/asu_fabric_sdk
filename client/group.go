package client

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/msp/mgmt"
)

//GrpcTLSOpt grpc tls opt(cert is []byte)
type GrpcTLSOpt struct {
	ClientCrt []byte
	ClientKey []byte
	Ca        []byte

	ServerNameOverride string
	Timeout            time.Duration
}

//Endpoint endpoint, such as: peer and orderer
type Endpoint struct {
	Address string
	GrpcTLSOpt
}

//Group peer and orderer group
type Group struct {
	peers    sync.Map
	orderers sync.Map
	signers  sync.Map
}

//NewGroup new clients group
func NewGroup() *Group {
	return &Group{}
}

//AddPeerClient add a peer client
func (g *Group) AddPeerClient(d Endpoint) error {
	c, err := NewClient(
		d.Address,
		d.ServerNameOverride,
		WithClientCert(d.ClientKey, d.ClientCrt),
		WithTLS(d.Ca),
	)
	if err != nil {
		return fmt.Errorf("new client failed: %v", err)
	}

	g.peers.Store(d.Address, c)
	return nil
}

//GetPeerClients get all peers' clients
func (g *Group) GetPeerClients() []*PeerClient {
	var c []*PeerClient

	g.peers.Range(func(key, value interface{}) bool {
		x, ok := value.(*Client)
		if !ok {
			return false
		}

		c = append(c, &PeerClient{*x})
		return false
	})

	return c
}

//GetPeerClient get one peer client
func (g *Group) GetPeerClient(address string) *PeerClient {
	p, _ := g.peers.Load(address)
	if p == nil {
		return nil
	}

	pp, ok := p.(*Client)
	if !ok {
		return nil
	}

	return &PeerClient{*pp}
}

//DeletePeerClient delete a peer client
func (g *Group) DeletePeerClient(address string) {
	g.peers.Delete(address)
}

//AddOrdererClient add a orderer client
func (g *Group) AddOrdererClient(d Endpoint) error {
	c, err := NewClient(
		d.Address,
		d.ServerNameOverride,
		WithClientCert(d.ClientKey, d.ClientCrt),
		WithTLS(d.Ca),
	)
	if err != nil {
		return fmt.Errorf("new client failed: %v", err)
	}

	g.orderers.Store(d.Address, c)
	return nil
}

//GetOrderersClients get all orderers' clients
func (g *Group) GetOrderersClients() []*OrdererClient {
	var c []*OrdererClient

	g.orderers.Range(func(key, value interface{}) bool {
		x, ok := value.(*Client)
		if !ok {
			return false
		}

		c = append(c, &OrdererClient{*x})
		return false
	})

	return c
}

//GetOrdererClient get one orderer client
func (g *Group) GetOrdererClient(address string) *OrdererClient {
	p, _ := g.orderers.Load(address)
	if p == nil {
		return nil
	}

	pp, ok := p.(*Client)
	if !ok {
		return nil
	}

	return &OrdererClient{*pp}
}

//DeleteOrdererClient delete a orderer client
func (g *Group) DeleteOrdererClient(address string) {
	g.orderers.Delete(address)
}

//EndorserProposal endorse proposal
func (g *Group) EndorserProposal(endorserAddress []string, sp *peer.SignedProposal) *peer.ProposalResponse {
	endorserMap := make(map[string]interface{})

	for ea := range endorserAddress {
		endorserMap[endorserAddress[ea]] = nil
	}

	var proposalResponse *peer.ProposalResponse

	g.peers.Range(func(key, value interface{}) bool {
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

//AddSigner add a msp
func (g *Group) addSigner(mspID, mspPath string) error {
	signer, err := getSigner(mspPath, mspID)
	if err != nil {
		return fmt.Errorf("get signer failed: %v", err)
	}

	g.signers.Store(mspID, signer)
	return nil
}

//GetSigner get map signing
// TODO can't use fabric msp load
func (g *Group) getSigner(mspID string) *msp.SigningIdentity {
	v, _ := g.signers.Load(mspID)
	if v == nil {
		return nil
	}

	x, ok := v.(msp.SigningIdentity)
	if !ok {
		return nil
	}

	return &x
}

//DeleteSigner delete a msp
// TODO can't use fabric msp load
func (g *Group) deleteSigner(mspID string) {
	g.signers.Delete(mspID)
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

// GetSigner initialize msp
func getSigner(mspPath, mspID string) (msp.SigningIdentity, error) {
	err := mgmt.LoadLocalMsp(
		mspPath, // core.yaml -> peer_mspConfigPath
		factory.GetDefaultOpts(),
		mspID, // peer_localMspId
		// msp.ProviderTypeToString(msp.FABRIC), // peer_localMspType, DEFAULT: SW
	)
	if err != nil {
		return nil, fmt.Errorf("load local msp failed: %v", err)
	}
	return mgmt.GetLocalMSP(factory.GetDefault()).GetDefaultSigningIdentity()
}
