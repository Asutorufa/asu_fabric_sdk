package chaincode

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/Asutorufa/fabricsdk/client"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/common/policydsl"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/msp/mgmt"
)

// GetSigner initialize msp
func GetSigner(mspPath, mspID string) (msp.SigningIdentity, error) {
	err := mgmt.LoadLocalMspWithType(
		mspPath, // core.yaml -> peer_mspConfigPath
		factory.GetDefaultOpts(),
		mspID,                                // peer_localMspId
		msp.ProviderTypeToString(msp.FABRIC), // peer_localMspType, DEFAULT: SW
	)
	if err != nil {
		return nil, err
	}
	return mgmt.GetLocalMSP(factory.GetDefault()).GetDefaultSigningIdentity()
}

// getChaincodeSpec
// path Chaincode Path
// name Chaincode Name
// version Chaincode Version
// isInit
// args Invoke or Query arguments
func getChaincodeSpec(
	path string,
	name string,
	isInit bool,
	version string,
	args [][]byte,
	Type peer.ChaincodeSpec_Type,
) *peer.ChaincodeSpec {
	return &peer.ChaincodeSpec{
		Type: Type, // <- from fabric-protos-go
		ChaincodeId: &peer.ChaincodeID{
			Path:    path,
			Name:    name,
			Version: version,
		},
		Input: &peer.ChaincodeInput{
			Args:        args,
			Decorations: map[string][]byte{},
			IsInit:      isInit,
		},
	}
}

func getChaincodeInvocationSpec(
	path string,
	name string,
	isInit bool,
	version string,
	Type peer.ChaincodeSpec_Type,
	args [][]byte) *peer.ChaincodeInvocationSpec {
	return &peer.ChaincodeInvocationSpec{
		ChaincodeSpec: getChaincodeSpec(
			path,
			name,
			isInit,
			version,
			args,
			Type,
		),
	}
}

type ChainOpt struct {
	Path                string
	Name                string
	Label               string
	IsInit              bool
	Version             string
	PackageID           string
	Sequence            int64
	EndorsementPlugin   string
	ValidationPlugin    string
	ValidationParameter []byte
	Policy              string
	// CollectionConfig    string
	CollectionsConfig []PrivateDataCollectionConfig
	// 详见: https://hyperledger-fabric.readthedocs.io/en/release-2.2/private_data_tutorial.html
	Type peer.ChaincodeSpec_Type
}

type PrivateDataCollectionConfig struct {
	Name              string
	Policy            string
	RequiredPeerCount int32
	MaxPeerCount      int32
	BlockToLive       uint64
	MemberOnlyRead    bool
	MemberOnlyWrite   bool
	EndorsementPolicy
}

type EndorsementPolicy struct {
	ChannelConfigPolicy string
	SignaturePolicy     string
}

type GrpcTLSOpt2 struct {
	ClientCrtPath string
	ClientKeyPath string
	CaPath        string

	ServerNameOverride string
	Timeout            time.Duration
}

type GrpcTLSOpt struct {
	ClientCrt []byte
	ClientKey []byte
	Ca        []byte

	ServerNameOverride string

	Timeout time.Duration
}

type Endpoint struct {
	Address string
	GrpcTLSOpt
}

type Endpoint2 struct {
	Address string
	GrpcTLSOpt2
}

func Endpoint2ToEndpoint(p Endpoint2) (Endpoint, error) {
	opt, err := GrpcTLSOpt2ToGrpcTLSOpt(p.GrpcTLSOpt2)
	if err != nil {
		return Endpoint{}, err
	}

	return Endpoint{
		Address:    p.Address,
		GrpcTLSOpt: opt,
	}, nil
}

func Endpoint2sToEndpoints(p []Endpoint2) ([]Endpoint, error) {
	var res []Endpoint
	for index := range p {
		tmp, err := Endpoint2ToEndpoint(p[index])
		if err != nil {
			return []Endpoint{}, fmt.Errorf("convert error -> %v", err)
		}

		res = append(res, tmp)
	}
	return res, nil
}

func GrpcTLSOpt2ToGrpcTLSOpt(g GrpcTLSOpt2) (gg GrpcTLSOpt, err error) {
	switch {
	case g.ClientCrtPath != "":
		gg.ClientCrt, err = ioutil.ReadFile(g.ClientCrtPath)
		if err != nil {
			return
		}
		fallthrough
	case g.ClientKeyPath != "":
		gg.ClientKey, err = ioutil.ReadFile(g.ClientKeyPath)
		if err != nil {
			return
		}
		fallthrough
	case g.CaPath != "":
		gg.Ca, err = ioutil.ReadFile(g.CaPath)
		if err != nil {
			return
		}
	}
	gg.ServerNameOverride = g.ServerNameOverride
	gg.Timeout = g.Timeout
	return
}

type MSPOpt struct {
	Path string
	Id   string
}

// processProposals sends a signed proposal to a set of peers, and gathers all the responses.
func processProposals(endorserClients []peer.EndorserClient, signedProposal *peer.SignedProposal) ([]*peer.ProposalResponse, error) {
	responsesCh := make(chan *peer.ProposalResponse, len(endorserClients))
	errorCh := make(chan error, len(endorserClients))
	wg := sync.WaitGroup{}
	for _, endorser := range endorserClients {
		wg.Add(1)
		go func(endorser peer.EndorserClient) {
			defer wg.Done()
			proposalResp, err := endorser.ProcessProposal(context.Background(), signedProposal)
			if err != nil {
				errorCh <- err
				return
			}
			responsesCh <- proposalResp
		}(endorser)
	}
	wg.Wait()
	close(responsesCh)
	close(errorCh)
	for err := range errorCh {
		return nil, err
	}
	var responses []*peer.ProposalResponse
	for response := range responsesCh {
		responses = append(responses, response)
	}
	return responses, nil
}

func NewDeliverClient(peer *client.PeerClient) (peer.DeliverClient, error) {
	return peer.PeerDeliver()
}

func NewEndorserClient(client *client.PeerClient) (peer.EndorserClient, error) {
	return client.Endorser()
}

var (
//spec *peer.ChaincodeSpec
//cID  string
//txID string
//signer identity.SignerSerializer
//certificate     tls.Certificate
//endorserClients []peer.EndorserClient
//deliverClients  []peer.DeliverClient

//bc common.BroadCastClient
//option string

// caFile string // <- orderer_tls_rootcert_file
// keyFile string // <- orderer_tls_clientKey_file
// certFile string // <- orderer_tls_clientCert_file
// orderingEndpoint string // <- orderer_address
// ordererTLSHostnameOverride // <- orderer_tls_serverhostoverride
// tlsEnabled bool // <- orderer_tls_enabled
// clientAuth bool // <- orderer_tls_clientAuthRequired
// connTimeout time.Duration // <- orderer_client_connTimeout
// tlsHandshakeTimeShift time.Duration // <- orderer_tls_handshakeTimeShift
)

func GetEndorserClient(client *client.PeerClient) (peer.EndorserClient, error) {
	return client.Endorser()
}

func GetDeliverClient(peer *client.PeerClient) (peer.DeliverClient, error) {
	return peer.PeerDeliver()
}

func GetCertificate(peer *client.PeerClient) tls.Certificate {
	return peer.Certificate()
}

func GetBroadcastClient(order *client.OrdererClient) (orderer.AtomicBroadcast_BroadcastClient, error) {
	return order.Broadcast()
}

type collectionConfigJson struct {
	Name              string `json:"name"`
	Policy            string `json:"policy"`
	RequiredPeerCount *int32 `json:"requiredPeerCount"`
	MaxPeerCount      *int32 `json:"maxPeerCount"`
	BlockToLive       uint64 `json:"blockToLive"`
	MemberOnlyRead    bool   `json:"memberOnlyRead"`
	MemberOnlyWrite   bool   `json:"memberOnlyWrite"`
	EndorsementPolicy *struct {
		SignaturePolicy     string `json:"signaturePolicy"`
		ChannelConfigPolicy string `json:"channelConfigPolicy"`
	} `json:"endorsementPolicy"`
}

// getCollectionConfig retrieves the collection configuration
// from the supplied byte array; the byte array must contain a
// json-formatted array of collectionConfigJson elements
func getCollectionConfigFromBytes(cconfBytes []byte) (*peer.CollectionConfigPackage, []byte, error) {
	cconf := &[]collectionConfigJson{}
	err := json.Unmarshal(cconfBytes, cconf)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not parse the collection configuration")
	}

	ccarray := make([]*peer.CollectionConfig, 0, len(*cconf))
	for _, cconfitem := range *cconf {
		p, err := policydsl.FromString(cconfitem.Policy)
		if err != nil {
			return nil, nil, errors.WithMessagef(err, "invalid policy %s", cconfitem.Policy)
		}

		cpc := &peer.CollectionPolicyConfig{
			Payload: &peer.CollectionPolicyConfig_SignaturePolicy{
				SignaturePolicy: p,
			},
		}

		var ep *peer.ApplicationPolicy
		if cconfitem.EndorsementPolicy != nil {
			signaturePolicy := cconfitem.EndorsementPolicy.SignaturePolicy
			channelConfigPolicy := cconfitem.EndorsementPolicy.ChannelConfigPolicy
			if (signaturePolicy != "" && channelConfigPolicy != "") || (signaturePolicy == "" && channelConfigPolicy == "") {
				return nil, nil, fmt.Errorf("incorrect policy")
			}

			if signaturePolicy != "" {
				poli, err := policydsl.FromString(signaturePolicy)
				if err != nil {
					return nil, nil, err
				}
				ep = &peer.ApplicationPolicy{
					Type: &peer.ApplicationPolicy_SignaturePolicy{
						SignaturePolicy: poli,
					},
				}
			} else {
				ep = &peer.ApplicationPolicy{
					Type: &peer.ApplicationPolicy_ChannelConfigPolicyReference{
						ChannelConfigPolicyReference: channelConfigPolicy,
					},
				}
			}
		}

		// Set default requiredPeerCount and MaxPeerCount if not specified in json
		requiredPeerCount := int32(0)
		maxPeerCount := int32(1)
		if cconfitem.RequiredPeerCount != nil {
			requiredPeerCount = *cconfitem.RequiredPeerCount
		}
		if cconfitem.MaxPeerCount != nil {
			maxPeerCount = *cconfitem.MaxPeerCount
		}

		cc := &peer.CollectionConfig{
			Payload: &peer.CollectionConfig_StaticCollectionConfig{
				StaticCollectionConfig: &peer.StaticCollectionConfig{
					Name:              cconfitem.Name,
					MemberOrgsPolicy:  cpc,
					RequiredPeerCount: requiredPeerCount,
					MaximumPeerCount:  maxPeerCount,
					BlockToLive:       cconfitem.BlockToLive,
					MemberOnlyRead:    cconfitem.MemberOnlyRead,
					MemberOnlyWrite:   cconfitem.MemberOnlyWrite,
					EndorsementPolicy: ep,
				},
			},
		}

		ccarray = append(ccarray, cc)
	}

	ccp := &peer.CollectionConfigPackage{Config: ccarray}
	ccpBytes, err := proto.Marshal(ccp)
	return ccp, ccpBytes, err
}
