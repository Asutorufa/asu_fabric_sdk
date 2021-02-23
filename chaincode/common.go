package chaincode

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	mb "github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/bccsp"
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

func newMSPConfig() {
	mspConfig := &mb.FabricMSPConfig{
		Name:                 "", // msp id
		RootCerts:            [][]byte{},
		IntermediateCerts:    [][]byte{},
		Admins:               [][]byte{},
		RevocationList:       [][]byte{},
		TlsRootCerts:         [][]byte{},
		TlsIntermediateCerts: [][]byte{},
		OrganizationalUnitIdentifiers: []*mb.FabricOUIdentifier{
			{
				OrganizationalUnitIdentifier: "",
				Certificate:                  []byte{},
			},
		},
		FabricNodeOus: &mb.FabricNodeOUs{},
		CryptoConfig: &mb.FabricCryptoConfig{ // Set FabricCryptoConfig
			SignatureHashFamily:            bccsp.SHA2,
			IdentityIdentifierHashFunction: bccsp.SHA256,
		},
		SigningIdentity: &mb.SigningIdentityInfo{
			PublicSigner:  []byte{}, // 节点证书(signCerts)
			PrivateSigner: nil,
		},
	}

	data, err := proto.Marshal(mspConfig)
	if err != nil {
		log.Println(err)
		return
	}

	_ = &mb.MSPConfig{
		Type:   int32(msp.FABRIC),
		Config: data,
	}
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
		ChaincodeSpec: getChaincodeSpec(path, name, isInit, version, args, Type),
	}
}

//ChainOpt chaincode about options for functions
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

//PrivateDataCollectionConfig private data collection config
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

func ConvertCollectionConfig(CollectionsConfig []PrivateDataCollectionConfig) (*peer.CollectionConfigPackage, error) {
	var collections *peer.CollectionConfigPackage

	for i := range CollectionsConfig {
		var ep *peer.ApplicationPolicy
		if CollectionsConfig[i].SignaturePolicy != "" &&
			CollectionsConfig[i].ChannelConfigPolicy != "" {
			return nil, fmt.Errorf("must spcify only one policy both SignaturePolicy and ChannelConfigPolicy")
		}
		if CollectionsConfig[i].SignaturePolicy != "" {
			p, err := policydsl.FromString(CollectionsConfig[i].SignaturePolicy)
			if err != nil {
				return nil, fmt.Errorf("format policy error -> %v", err)
			}

			ep = &peer.ApplicationPolicy{
				Type: &peer.ApplicationPolicy_SignaturePolicy{
					SignaturePolicy: p,
				},
			}
		}
		if CollectionsConfig[i].ChannelConfigPolicy != "" {
			ep = &peer.ApplicationPolicy{
				Type: &peer.ApplicationPolicy_ChannelConfigPolicyReference{
					ChannelConfigPolicyReference: CollectionsConfig[i].ChannelConfigPolicy,
				},
			}
		}
		p, err := policydsl.FromString(CollectionsConfig[i].Policy)
		if err != nil {
			return nil, fmt.Errorf("policy string error -> %v", err)
		}

		cc := &peer.CollectionConfig{
			Payload: &peer.CollectionConfig_StaticCollectionConfig{
				StaticCollectionConfig: &peer.StaticCollectionConfig{
					Name: CollectionsConfig[i].Name,
					MemberOrgsPolicy: &peer.CollectionPolicyConfig{
						Payload: &peer.CollectionPolicyConfig_SignaturePolicy{
							SignaturePolicy: p,
						},
					},
					RequiredPeerCount: CollectionsConfig[i].RequiredPeerCount,
					MaximumPeerCount:  CollectionsConfig[i].MaxPeerCount,
					BlockToLive:       CollectionsConfig[i].BlockToLive,
					MemberOnlyRead:    CollectionsConfig[i].MemberOnlyRead,
					MemberOnlyWrite:   CollectionsConfig[i].MemberOnlyWrite,
					EndorsementPolicy: ep,
				},
			},
		}

		collections.Config = append(collections.Config, cc)
	}

	return collections, nil
}

//EndorsementPolicy endorser policy
type EndorsementPolicy struct {
	ChannelConfigPolicy string
	SignaturePolicy     string
}

//GrpcTLSOptWithPath grpc tls opt(cert is path)
type GrpcTLSOptWithPath struct {
	ClientCrtPath string // for client auth req
	ClientKeyPath string // for client auth req
	CaPath        string // tls ca cert

	ServerNameOverride string
	Timeout            time.Duration
}

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

//EndpointWithPath endpoint, such as: peer and orderer(cert is path)
type EndpointWithPath struct {
	Address string
	GrpcTLSOptWithPath
}

//ParseEndpointWithPath EndpointWithPath type to Endpoint
func ParseEndpointWithPath(p EndpointWithPath) (Endpoint, error) {
	opt, err := ParseGrpcTLSOptWithPath(p.GrpcTLSOptWithPath)
	if err != nil {
		return Endpoint{}, err
	}

	return Endpoint{
		Address:    p.Address,
		GrpcTLSOpt: opt,
	}, nil
}

//ParseEndpointsWithPath EndpointWithPath type array to Endpoint array
func ParseEndpointsWithPath(p []EndpointWithPath) ([]Endpoint, error) {
	var res []Endpoint
	for index := range p {
		tmp, err := ParseEndpointWithPath(p[index])
		if err != nil {
			return []Endpoint{}, fmt.Errorf("convert error -> %v", err)
		}

		res = append(res, tmp)
	}
	return res, nil
}

//ParseGrpcTLSOptWithPath GrpcTLSOptWithPath type to GrpcTLSOpt
func ParseGrpcTLSOptWithPath(g GrpcTLSOptWithPath) (gg GrpcTLSOpt, err error) {
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

//MSPOpt msp about options
type MSPOpt struct {
	Path string
	ID   string
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

// var (
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
// )

//collectionConfigJSON private data config json
type collectionConfigJSON struct {
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
	cconf := &[]collectionConfigJSON{}
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
