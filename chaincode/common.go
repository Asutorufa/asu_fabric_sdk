package chaincode

import (
	"context"
	"crypto/tls"
	"fabricSDK/chaincode/client/orderclient"
	"fabricSDK/chaincode/client/peerclient"
	"sync"
	"time"

	"github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/msp/mgmt"
)

func GetSigner(mspPath, mspId string) (msp.SigningIdentity, error) {
	err := mgmt.LoadLocalMspWithType(
		mspPath, // core.yaml -> peer_mspConfigPath
		factory.GetDefaultOpts(),
		mspId,                                // peer_localMspId
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
) *peer.ChaincodeSpec {
	return &peer.ChaincodeSpec{
		Type: peer.ChaincodeSpec_GOLANG, // <- from fabric-protos-go
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
	args [][]byte) *peer.ChaincodeInvocationSpec {
	return &peer.ChaincodeInvocationSpec{
		ChaincodeSpec: getChaincodeSpec(
			path,
			name,
			isInit,
			version,
			args,
		),
	}
}

type ChainOpt struct {
	Path    string
	Name    string
	IsInit  bool
	Version string
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

func NewDeliverClient(peer *peerclient.PeerClient) (peer.DeliverClient, error) {
	return peer.PeerDeliver()
}

func NewEndorserClient(client *peerclient.PeerClient) (peer.EndorserClient, error) {
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

func GetEndorserClient(client *peerclient.PeerClient) (peer.EndorserClient, error) {
	return client.Endorser()
}

func GetDeliverClient(peer *peerclient.PeerClient) (peer.DeliverClient, error) {
	return peer.PeerDeliver()
}

func GetCertificate(peer *peerclient.PeerClient) tls.Certificate {
	return peer.Certificate()
}

func GetBroadcastClient(order *orderclient.OrdererClient) (orderer.AtomicBroadcast_BroadcastClient, error) {
	return order.Broadcast()
}
