package chaincode

import (
	"fabricSDK/chaincode/client/clientcommon"
	"fabricSDK/chaincode/client/peerclient"
	"fmt"
	"io/ioutil"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/protoutil"
)

func Query2(
	chaincode ChainOpt,
	peerGrpcOpt GrpcTLSOpt2,
	mspOpt MSPOpt,
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
		mspOpt,
		args,
		channelID,
		"",
		peerAddress,
	)
}

func Query(
	chaincode ChainOpt,
	peerGrpcOpt GrpcTLSOpt,
	mspOpt MSPOpt,
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
	signer, err := GetSigner(mspOpt.Path, mspOpt.Id)
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

	var peerClients []*peerclient.PeerClient
	var endorserClients []peer.EndorserClient
	for index := range peerAddress {
		peerClient, err := peerclient.NewPeerClient(
			peerAddress[index],
			peerGrpcOpt.ServerNameOverride,
			clientcommon.WithTLS2(peerGrpcOpt.Ca),
			clientcommon.WithClientCert2(peerGrpcOpt.ClientKey, peerGrpcOpt.ClientCrt),
			clientcommon.WithTimeout(peerGrpcOpt.Timeout),
		)
		if err != nil {
			return nil, fmt.Errorf("NewPeerClient() -> %v", err)
		}
		peerClients = append(peerClients, peerClient)

		endorserClient, err := peerClient.Endorser()
		if err != nil {
			return nil, fmt.Errorf("peerClient.Endorser() -> %v", err)
		}
		endorserClients = append(endorserClients, endorserClient)
	}
	defer func() {
		for index := range peerClients {
			peerClients[index].Close()
		}
	}()

	responses, err := processProposals(endorserClients, signedProp)
	if err != nil {
		return nil, fmt.Errorf("processProposals() -> %v", err)
	}
	fmt.Printf("txid: %s\n", txid)
	return responses[0], nil
}
