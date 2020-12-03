package chaincode

import (
	"context"
	"crypto/tls"
	"errors"
	"fabricSDK/chaincode/orderclient"
	"fabricSDK/chaincode/peerclient"
	"fmt"
	"io/ioutil"
	"math"
	"sync"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/msp/mgmt"
	"github.com/hyperledger/fabric/protoutil"
)

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

func GetSigner() (msp.SigningIdentity, error) {
	err := mgmt.LoadLocalMspWithType(
		"/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/peers/peer-0-baas98/msp", // core.yaml -> peer_mspConfigPath
		factory.GetDefaultOpts(),
		"baas98",                             // peer_localMspId
		msp.ProviderTypeToString(msp.FABRIC), // peer_localMspType, DEFAULT: SW
	)
	if err != nil {
		return nil, err
	}
	return mgmt.GetLocalMSP(factory.GetDefault()).GetDefaultSigningIdentity()
}

func GetBroadcastClient(order *orderclient.OrdererClient) (orderer.AtomicBroadcast_BroadcastClient, error) {
	return order.Broadcast()
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

func Query2(
	chaincode ChainOpt,
	peerGrpcOpt GrpcTLSOpt2,
	args [][]byte,
	channelID string,
	peerAddress []string,
) (*peer.ProposalResponse, error) {
	grpc := GrpcTLSOpt{
		ServerNameOverride: peerGrpcOpt.ServerNameOverride,
		Timeout:            peerGrpcOpt.Timeout,
	}
	var err error
	switch {
	case peerGrpcOpt.CaPath != "":
		grpc.Ca, err = ioutil.ReadFile(peerGrpcOpt.CaPath)
		if err != nil {
			return nil, err
		}
		fallthrough
	case peerGrpcOpt.ClientKeyPath != "":
		grpc.ClientKey, err = ioutil.ReadFile(peerGrpcOpt.ClientKeyPath)
		if err != nil {
			return nil, err
		}
		fallthrough
	case peerGrpcOpt.ClientCrtPath != "":
		grpc.ClientCrt, err = ioutil.ReadFile(peerGrpcOpt.ClientCrtPath)
		if err != nil {
			return nil, err
		}
	}
	return Query(
		chaincode,
		grpc,
		args,
		channelID,
		"",
		peerAddress,
	)
}

func Query(
	chaincode ChainOpt,
	peerGrpcOpt GrpcTLSOpt,
	args [][]byte,
	channelID string,
	txID string,
	peerAddress []string,
) (*peer.ProposalResponse, error) {
	invocation := getChaincodeInvocationSpec(
		chaincode.Path,
		chaincode.Name,
		chaincode.IsInit,
		chaincode.Version,
		args,
	)
	signer, err := GetSigner()
	if err != nil {
		return nil, fmt.Errorf("GetSigner() -> %v", err)
	}
	creator, err := signer.Serialize()
	if err != nil {
		return nil, fmt.Errorf("signer.Serialize() -> %v", err)
	}

	prop, txid, err := protoutil.CreateChaincodeProposalWithTxIDAndTransient(
		common.HeaderType_ENDORSER_TRANSACTION,
		channelID,
		invocation,
		creator,
		txID,
		map[string][]byte{},
	)
	if err != nil {
		return nil, fmt.Errorf("protoutil.CreateChaincodeProposalWithTxIDAndTransient() -> %v", err)
	}

	signedProp, err := protoutil.GetSignedProposal(prop, signer)
	if err != nil {
		return nil, fmt.Errorf("protoutil.GetSignedProposal() -> %v", err)
	}
	var endorserClients []peer.EndorserClient
	for index := range peerAddress {
		peerClient, err := peerclient.NewPeerClient(
			peerAddress[index],
			peerGrpcOpt.ServerNameOverride,
			peerclient.WithTLS2(peerGrpcOpt.Ca),
			peerclient.WithClientCert2(peerGrpcOpt.ClientKey, peerGrpcOpt.ClientCrt),
			peerclient.WithTimeout(peerGrpcOpt.Timeout),
		)
		if err != nil {
			return nil, fmt.Errorf("NewPeerClient() -> %v", err)
		}
		endorserClient, err := peerClient.Endorser()
		if err != nil {
			return nil, fmt.Errorf("peerClient.Endorser() -> %v", err)
		}
		endorserClients = append(endorserClients, endorserClient)
	}
	responses, err := processProposals(endorserClients, signedProp)
	if err != nil {
		return nil, fmt.Errorf("processProposals() -> %v", err)
	}
	fmt.Printf("txid: %s\n", txid)
	return responses[0], nil
}

func Invoke2(
	chaincode ChainOpt,
	peerGrpcOpt GrpcTLSOpt2,
	ordererGrpcOpt GrpcTLSOpt2,
	args [][]byte, // [][]byte{[]byte("function"),[]byte("a"),[]byte("b")}, first array is function name
	channelID string,
	peerAddress []string,
	ordererAddress string,
) (*peer.ProposalResponse, error) {
	peerGrpc := GrpcTLSOpt{
		ServerNameOverride: peerGrpcOpt.ServerNameOverride,
		Timeout:            peerGrpcOpt.Timeout,
	}
	ordererGrpc := GrpcTLSOpt{
		ServerNameOverride: ordererGrpcOpt.ServerNameOverride,
		Timeout:            ordererGrpcOpt.Timeout,
	}
	var err error
	switch {
	case peerGrpcOpt.CaPath != "":
		peerGrpc.Ca, err = ioutil.ReadFile(peerGrpcOpt.CaPath)
		if err != nil {
			return nil, err
		}
		fallthrough
	case peerGrpcOpt.ClientKeyPath != "":
		peerGrpc.ClientKey, err = ioutil.ReadFile(peerGrpcOpt.ClientKeyPath)
		if err != nil {
			return nil, err
		}
		fallthrough
	case peerGrpcOpt.ClientCrtPath != "":
		peerGrpc.ClientCrt, err = ioutil.ReadFile(peerGrpcOpt.ClientCrtPath)
		if err != nil {
			return nil, err
		}
		fallthrough
	case ordererGrpcOpt.CaPath != "":
		ordererGrpc.Ca, err = ioutil.ReadFile(ordererGrpcOpt.CaPath)
		if err != nil {
			return nil, err
		}
		fallthrough
	case ordererGrpcOpt.ClientKeyPath != "":
		ordererGrpc.ClientKey, err = ioutil.ReadFile(ordererGrpcOpt.ClientKeyPath)
		if err != nil {
			return nil, err
		}
		fallthrough
	case ordererGrpcOpt.ClientCrtPath != "":
		ordererGrpc.ClientCrt, err = ioutil.ReadFile(ordererGrpcOpt.ClientCrtPath)
		if err != nil {
			return nil, err
		}
	}
	return Invoke(
		chaincode,
		peerGrpc,
		ordererGrpc,
		args,
		channelID,
		"",
		peerAddress,
		ordererAddress,
	)
}

func Invoke(
	chaincode ChainOpt,
	peerGrpcOpt GrpcTLSOpt,
	ordererGrpcOpt GrpcTLSOpt,
	args [][]byte, // [][]byte{[]byte("function"),[]byte("a"),[]byte("b")}, first array is function name
	channelID string,
	txID string,
	peerAddress []string,
	ordererAddress string,
) (*peer.ProposalResponse, error) {

	invocation := getChaincodeInvocationSpec(
		chaincode.Path,
		chaincode.Name,
		chaincode.IsInit,
		chaincode.Version,
		args,
	)
	signer, err := GetSigner()
	if err != nil {
		return nil, err
	}
	creator, err := signer.Serialize()
	if err != nil {
		return nil, err
	}

	//tMap := map[string][]byte{
	//	"cert": []byte("transient"),
	//}

	prop, txid, err := protoutil.CreateChaincodeProposalWithTxIDAndTransient(
		common.HeaderType_ENDORSER_TRANSACTION,
		channelID,
		invocation,
		creator,
		txID,
		map[string][]byte{},
	)
	if err != nil {
		return nil, err
	}

	signedProp, err := protoutil.GetSignedProposal(prop, signer)
	if err != nil {
		return nil, err
	}

	var endorserClients []peer.EndorserClient
	var deliverClients []peer.DeliverClient
	var certificate tls.Certificate
	for index := range peerAddress {
		peerClient, err := peerclient.NewPeerClient(
			peerAddress[index],
			peerGrpcOpt.ServerNameOverride,
			peerclient.WithClientCert2(peerGrpcOpt.ClientKey, peerGrpcOpt.ClientCrt),
			peerclient.WithTLS2(peerGrpcOpt.Ca),
			peerclient.WithTimeout(peerGrpcOpt.Timeout),
		)
		if err != nil {
			return nil, err
		}
		certificate = peerClient.Certificate()
		endorserClient, err := peerClient.Endorser()
		if err != nil {
			return nil, err
		}
		endorserClients = append(endorserClients, endorserClient)

		deliverClient, err := peerClient.PeerDeliver()
		if err != nil {
			return nil, err
		}
		deliverClients = append(deliverClients, deliverClient)
	}
	responses, err := processProposals(endorserClients, signedProp)
	if err != nil {
		return nil, err
	}
	fmt.Printf("txid: %s\n", txid)
	resp := responses[0]
	if resp == nil {
		return resp, nil
	}

	if resp.Response.Status >= shim.ERRORTHRESHOLD {
		return resp, nil
	}

	env, err := protoutil.CreateSignedTx(prop, signer, responses...)
	if err != nil {
		return resp, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dg := NewDeliverGroup(
		deliverClients,
		peerAddress,
		signer,
		certificate,
		channelID,
		txid,
	)

	err = dg.Connect(ctx)
	if err != nil {
		return nil, err
	}

	order, err := orderclient.NewOrdererClient(
		ordererAddress,
		ordererGrpcOpt.ServerNameOverride,
		orderclient.WithClientCert2(ordererGrpcOpt.ClientKey, ordererGrpcOpt.ClientCrt),
		orderclient.WithTLS2(ordererGrpcOpt.Ca),
		orderclient.WithTimeout(ordererGrpcOpt.Timeout),
	)
	if err != nil {
		return nil, err
	}
	ordererClient, err := GetBroadcastClient(order)
	if err != nil {
		return nil, err
	}
	err = ordererClient.Send(env)
	if err != nil {
		return resp, err
	}

	if dg != nil && ctx != nil {
		err = dg.Wait(ctx)
		if err != nil {
			return nil, fmt.Errorf("dg.Wait() -> %v", err)
		}
	}
	return resp, nil
}

// DeliverGroup holds all of the information needed to connect
// to a set of peers to wait for the interested txid to be
// committed to the ledgers of all peers. This functionality
// is currently implemented via the peer's DeliverFiltered service.
// An error from any of the peers/deliver clients will result in
// the invoke command returning an error. Only the first error that
// occurs will be set
type DeliverGroup struct {
	Clients     []*DeliverClient
	Certificate tls.Certificate
	ChannelID   string
	TxID        string
	Signer      msp.SigningIdentity
	mutex       sync.Mutex
	Error       error
	wg          sync.WaitGroup
}

// DeliverClient holds the client/connection related to a specific
// peer. The address is included for logging purposes
type DeliverClient struct {
	Client     peer.DeliverClient
	Connection peer.Deliver_DeliverClient
	Address    string
}

func NewDeliverGroup(
	deliverClients []peer.DeliverClient,
	peerAddresses []string,
	signer msp.SigningIdentity,
	certificate tls.Certificate,
	channelID string,
	txid string,
) *DeliverGroup {
	clients := make([]*DeliverClient, len(deliverClients))
	for i, client := range deliverClients {
		address := peerAddresses[i]
		//if address == "" {
		//	address = viper.GetString("peer.address")
		//}
		dc := &DeliverClient{
			Client:  client,
			Address: address,
		}
		clients[i] = dc
	}

	dg := &DeliverGroup{
		Clients:     clients,
		Certificate: certificate,
		ChannelID:   channelID,
		TxID:        txid,
		Signer:      signer,
	}

	return dg
}

// Connect waits for all deliver clients in the group to connect to
// the peer's deliver service, receive an error, or for the context
// to timeout. An error will be returned whenever even a single
// deliver client fails to connect to its peer
func (dg *DeliverGroup) Connect(ctx context.Context) error {
	dg.wg.Add(len(dg.Clients))
	for _, client := range dg.Clients {
		go dg.ClientConnect(ctx, client)
	}
	readyCh := make(chan struct{})
	go dg.WaitForWG(readyCh)

	select {
	case <-readyCh:
		if dg.Error != nil {
			err := fmt.Errorf("%v failed to connect to deliver on all peers", dg.Error)
			return err
		}
	case <-ctx.Done():
		err := errors.New("timed out waiting for connection to deliver on all peers")
		return err
	}

	return nil
}

// ClientConnect sends a deliver seek info envelope using the
// provided deliver client, setting the deliverGroup's Error
// field upon any error
func (dg *DeliverGroup) ClientConnect(ctx context.Context, dc *DeliverClient) {
	defer dg.wg.Done()
	df, err := dc.Client.DeliverFiltered(ctx)
	if err != nil {
		//err = errors.WithMessagef(err, "error connecting to deliver filtered at %s", dc.Address)
		dg.setError(err)
		return
	}
	defer df.CloseSend()
	dc.Connection = df

	envelope := createDeliverEnvelope(dg.ChannelID, dg.Certificate, dg.Signer)
	err = df.Send(envelope)
	if err != nil {
		//err = errors.WithMessagef(err, "error sending deliver seek info envelope to %s", dc.Address)
		dg.setError(err)
		return
	}
}

// Wait waits for all deliver client connections in the group to
// either receive a block with the txid, an error, or for the
// context to timeout
func (dg *DeliverGroup) Wait(ctx context.Context) error {
	if len(dg.Clients) == 0 {
		return nil
	}

	dg.wg.Add(len(dg.Clients))
	for _, client := range dg.Clients {
		go dg.ClientWait(client)
	}
	readyCh := make(chan struct{})
	go dg.WaitForWG(readyCh)

	select {
	case <-readyCh:
		if dg.Error != nil {
			return dg.Error
		}
	case <-ctx.Done():
		err := errors.New("timed out waiting for txid on all peers")
		return err
	}

	return nil
}

// ClientWait waits for the specified deliver client to receive
// a block event with the requested txid
func (dg *DeliverGroup) ClientWait(dc *DeliverClient) {
	defer dg.wg.Done()
	for {
		resp, err := dc.Connection.Recv()
		if err != nil {
			//err = errors.WithMessagef(err, "error receiving from deliver filtered at %s", dc.Address)
			dg.setError(err)
			return
		}
		switch r := resp.Type.(type) {
		case *peer.DeliverResponse_FilteredBlock:
			filteredTransactions := r.FilteredBlock.FilteredTransactions
			for _, tx := range filteredTransactions {
				if tx.Txid == dg.TxID {
					//logger.Infof("txid [%s] committed with status (%s) at %s", dg.TxID, tx.TxValidationCode, dc.Address)
					if tx.TxValidationCode != peer.TxValidationCode_VALID {
						//err = errors.Errorf("transaction invalidated with status (%s)", tx.TxValidationCode)
						dg.setError(err)
					}
					return
				}
			}
		case *peer.DeliverResponse_Status:
			//err = errors.Errorf("deliver completed with status (%s) before txid received", r.Status)
			dg.setError(err)
			return
		default:
			//err = errors.Errorf("received unexpected response type (%T) from %s", r, dc.Address)
			dg.setError(err)
			return
		}
	}
}

// WaitForWG waits for the deliverGroup's wait group and closes
// the channel when ready
func (dg *DeliverGroup) WaitForWG(readyCh chan struct{}) {
	dg.wg.Wait()
	close(readyCh)
}

// setError serializes an error for the deliverGroup
func (dg *DeliverGroup) setError(err error) {
	dg.mutex.Lock()
	dg.Error = err
	dg.mutex.Unlock()
}

func createDeliverEnvelope(
	channelID string,
	certificate tls.Certificate,
	signer msp.SigningIdentity,
) *common.Envelope {
	var tlsCertHash []byte
	// check for client certificate and create hash if present
	if len(certificate.Certificate) > 0 {
		tlsCertHash = util.ComputeSHA256(certificate.Certificate[0])
	}

	start := &orderer.SeekPosition{
		Type: &orderer.SeekPosition_Newest{
			Newest: &orderer.SeekNewest{},
		},
	}

	stop := &orderer.SeekPosition{
		Type: &orderer.SeekPosition_Specified{
			Specified: &orderer.SeekSpecified{
				Number: math.MaxUint64,
			},
		},
	}

	seekInfo := &orderer.SeekInfo{
		Start:    start,
		Stop:     stop,
		Behavior: orderer.SeekInfo_BLOCK_UNTIL_READY,
	}

	env, err := protoutil.CreateSignedEnvelopeWithTLSBinding(
		common.HeaderType_DELIVER_SEEK_INFO,
		channelID,
		signer,
		seekInfo,
		int32(0),
		uint64(0),
		tlsCertHash,
	)
	if err != nil {
		//logger.Errorf("Error signing envelope: %s", err)
		return nil
	}

	return env
}
