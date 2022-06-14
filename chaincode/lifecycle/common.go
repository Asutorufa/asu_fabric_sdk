package lifecycle

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/Asutorufa/fabricsdk/client"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/protoutil"
)

const (
	lifecycleName                = "_lifecycle"
	approveFuncName              = "ApproveChaincodeDefinitionForMyOrg"
	commitFuncName               = "CommitChaincodeDefinition"
	checkCommitReadinessFuncName = "CheckCommitReadiness"
)

func peerInvoke(
	peers []chaincode.Endpoint,
	signedProposal *peer.SignedProposal,
) ([]*peer.ProposalResponse, error) {
	var resps []*peer.ProposalResponse
	for _, peer := range peers {
		peerClient, err := client.NewPeerClientSelf(
			peer.Address,
			peer.ServerNameOverride,
			client.WithClientCert(peer.ClientKey, peer.ClientCrt),
			client.WithTLS(peer.Ca),
			client.WithTimeout(6*time.Second),
		)
		if err != nil {
			return nil, err
		}

		endorserClient, err := peerClient.Endorser()
		if err != nil {
			return nil, err
		}

		resp, err := endorserClient.ProcessProposal(context.Background(), signedProposal)
		if err != nil {
			return nil, err
		}

		resps = append(resps, resp)
	}

	if len(resps) == 0 {
		return nil, fmt.Errorf("all peers response is empty")
	}

	resp := resps[0]

	if resp == nil {
		return nil, fmt.Errorf("resp is nil")
	}

	if resp.Response.Status != int32(common.Status_SUCCESS) {
		return nil, fmt.Errorf("%d - %s", resp.Response.Status, resp.Response.Message)
	}

	return resps, nil
}

func query(signer msp.SigningIdentity, proposal *peer.Proposal,
	peers []chaincode.Endpoint) (*peer.ProposalResponse, error) {
	signedProposal, err := signProposal(proposal, signer)
	if err != nil {
		return nil, err
	}

	var resps []*peer.ProposalResponse
	for _, peer := range peers {
		peerClient, err := client.NewPeerClientSelf(
			peer.Address,
			peer.ServerNameOverride,
			client.WithClientCert(peer.ClientKey, peer.ClientCrt),
			client.WithTLS(peer.Ca),
			client.WithTimeout(6*time.Second),
		)
		if err != nil {
			return nil, err
		}

		endorserClient, err := peerClient.Endorser()
		if err != nil {
			return nil, err
		}

		resp, err := endorserClient.ProcessProposal(context.Background(), signedProposal)
		if err != nil {
			return nil, err
		}

		resps = append(resps, resp)
	}

	if len(resps) == 0 {
		return nil, fmt.Errorf("all peers response is empty")
	}

	resp := resps[0]

	if resp == nil {
		return nil, fmt.Errorf("resp is nil")
	}

	if resp.Response.Status != int32(common.Status_SUCCESS) {
		return nil, fmt.Errorf("%d - %s", resp.Response.Status, resp.Response.Message)
	}

	return resp, nil
}

func queryAll(signer msp.SigningIdentity, proposal *peer.Proposal,
	peers []chaincode.Endpoint) ([]*peer.ProposalResponse, error) {
	peerClients := chaincode.GetPeerClients(peers)
	if len(peerClients) == 0 {
		return nil, fmt.Errorf("no peer can be connect[peerClients' size is 0]")
	}
	defer chaincode.CloseClients(peerClients)

	return internalQueryAll(signer, proposal, peerClients)
}

func internalQueryAll(signer msp.SigningIdentity, proposal *peer.Proposal,
	peers []*client.PeerClient) ([]*peer.ProposalResponse, error) {
	signedProposal, err := signProposal(proposal, signer)
	if err != nil {
		return nil, err
	}

	var resps []*peer.ProposalResponse
	for _, peer := range peers {
		endorserClient, err := peer.Endorser()
		if err != nil {
			return nil, err
		}

		resp, err := endorserClient.ProcessProposal(context.Background(), signedProposal)
		if err != nil {
			return nil, err
		}

		resps = append(resps, resp)
	}

	if len(resps) == 0 {
		return nil, fmt.Errorf("all peers response is empty")
	}

	resp := resps[0]

	if resp == nil {
		return nil, fmt.Errorf("resp is nil")
	}

	if resp.Response.Status != int32(common.Status_SUCCESS) {
		return nil, fmt.Errorf("%d - %s", resp.Response.Status, resp.Response.Message)
	}

	return resps, nil
}

func invoke(signer msp.SigningIdentity, proposal *peer.Proposal,
	peers []chaincode.Endpoint, orderers []chaincode.Endpoint,
	channelID string, txID string) (*peer.ProposalResponse, error) {
	peerClients := chaincode.GetPeerClients(peers)
	if len(peerClients) == 0 {
		return nil, fmt.Errorf("peer clients' is 0")
	}
	defer chaincode.CloseClients(peerClients)

	ordererClients := chaincode.GetOrdererClients(orderers)
	if len(ordererClients) == 0 {
		return nil, fmt.Errorf("orderer clients' is 0")
	}
	defer chaincode.CloseClients(ordererClients)

	return internalInvoke(signer, proposal, peerClients, ordererClients, channelID, txID)
}

func internalInvoke(signer msp.SigningIdentity, proposal *peer.Proposal, peers []*client.PeerClient,
	orderers []*client.OrdererClient, channelID string, txID string) (*peer.ProposalResponse, error) {
	resp, err := internalQueryAll(signer, proposal, peers)
	if err != nil {
		return nil, fmt.Errorf("invoke from peers error -> %v", err)
	}

	env, err := protoutil.CreateSignedTx(proposal, signer, resp...)
	if err != nil {
		return nil, fmt.Errorf("failed to create signed transaction -> %v", err)
	}

	//
	//            orderers
	//
	var endorserClients []peer.EndorserClient
	var deliverClients []peer.DeliverClient
	var certificate tls.Certificate
	for pi := range peers {
		certificate = peers[pi].Certificate()
		endorserClient, err := peers[pi].Endorser()
		if err != nil {
			return nil, err
		}
		endorserClients = append(endorserClients, endorserClient)

		deliverClient, err := peers[pi].PeerDeliver()
		if err != nil {
			return nil, err
		}
		deliverClients = append(deliverClients, deliverClient)
	}

	dg := chaincode.NewDeliverGroup(
		deliverClients,
		signer,
		certificate,
		channelID,
		txID,
	)

	for _, orderer := range orderers {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		err = dg.Connect(ctx)
		if err != nil {
			// return nil, err
			log.Println(err)
			continue
		}

		ordererClient, err := orderer.Broadcast()
		if err != nil {
			// return nil, err
			log.Println(err)
			continue
		}
		err = ordererClient.Send(env)
		if err != nil {
			// return nil, err
			log.Println(err)
			continue
		}

		if dg != nil && ctx != nil {
			err = dg.Wait(ctx)
			if err != nil {
				// return nil, fmt.Errorf("dg.Wait() -> %v", err)
				log.Printf("dg.Wait() -> %v\n", err)
				continue
			}
		}
		return resp[0], nil
	}

	return nil, fmt.Errorf("failed send envelop to all orderers")
}

func signProposal(proposal *peer.Proposal, signer msp.SigningIdentity) (*peer.SignedProposal, error) {
	// check for nil argument
	if proposal == nil {
		return nil, errors.New("proposal cannot be nil")
	}

	if signer == nil {
		return nil, errors.New("signer cannot be nil")
	}

	proposalBytes, err := proto.Marshal(proposal)
	if err != nil {
		return nil, fmt.Errorf("error marshaling proposal: %w", err)
	}

	signature, err := signer.Sign(proposalBytes)
	if err != nil {
		return nil, err
	}

	return &peer.SignedProposal{
		ProposalBytes: proposalBytes,
		Signature:     signature,
	}, nil
}

func createProposal(
	args proto.Message,
	signer msp.SigningIdentity,
	function, channel string,
) (*peer.Proposal, string, error) {
	argsBytes, err := proto.Marshal(args)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal args: %w", err)
	}
	ccInput := &peer.ChaincodeInput{Args: [][]byte{[]byte(function), argsBytes}}

	cis := &peer.ChaincodeInvocationSpec{
		ChaincodeSpec: &peer.ChaincodeSpec{
			ChaincodeId: &peer.ChaincodeID{Name: lifecycleName},
			Input:       ccInput,
		},
	}

	signerSerialized, err := signer.Serialize()
	if err != nil {
		return nil, "", fmt.Errorf("failed to serialize identity: %w", err)
	}

	proposal, txID, err := protoutil.CreateProposalFromCIS(common.HeaderType_ENDORSER_TRANSACTION, channel, cis, signerSerialized)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create ChaincodeInvocationSpec proposal: %w", err)
	}

	return proposal, txID, nil
}
