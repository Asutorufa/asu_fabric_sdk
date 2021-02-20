package lifecycle

import (
	"context"
	"crypto/tls"
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
	"github.com/pkg/errors"
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
		peerClient, err := client.NewPeerClient(
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

func query(
	signer msp.SigningIdentity,
	proposal *peer.Proposal,
	peers []chaincode.Endpoint,
) (*peer.ProposalResponse, error) {
	signedProposal, err := signProposal(proposal, signer)
	if err != nil {
		return nil, err
	}

	var resps []*peer.ProposalResponse
	for _, peer := range peers {
		peerClient, err := client.NewPeerClient(
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

func queryAll(
	signer msp.SigningIdentity,
	proposal *peer.Proposal,
	peers []chaincode.Endpoint,
) ([]*peer.ProposalResponse, error) {
	signedProposal, err := signProposal(proposal, signer)
	if err != nil {
		return nil, err
	}

	var resps []*peer.ProposalResponse
	for _, peer := range peers {
		peerClient, err := client.NewPeerClient(
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

func invoke(
	signer msp.SigningIdentity,
	proposal *peer.Proposal,
	peers []chaincode.Endpoint,
	orderers []chaincode.Endpoint,
	channelID string,
	txID string,
) (*peer.ProposalResponse, error) {
	resp, err := queryAll(signer, proposal, peers)
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
	var peerClients []*client.PeerClient
	var endorserClients []peer.EndorserClient
	var deliverClients []peer.DeliverClient
	var certificate tls.Certificate
	for index := range peers {
		peerClient, err := client.NewPeerClient(
			peers[index].Address,
			peers[index].GrpcTLSOpt.ServerNameOverride,
			client.WithClientCert(peers[index].GrpcTLSOpt.ClientKey, peers[index].GrpcTLSOpt.ClientCrt),
			client.WithTLS(peers[index].GrpcTLSOpt.Ca),
			client.WithTimeout(peers[index].GrpcTLSOpt.Timeout),
		)
		if err != nil {
			return nil, err
		}
		peerClients = append(peerClients, peerClient)
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
	defer func() {
		for index := range peerClients {
			_ = peerClients[index].Close()
		}
	}()
	var eps []string
	for index := range peers {
		eps = append(eps, peers[index].Address)
	}

	dg := chaincode.NewDeliverGroup(
		deliverClients,
		eps,
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

		order, err := client.NewOrdererClient(
			orderer.Address,
			orderer.GrpcTLSOpt.ServerNameOverride,
			client.WithClientCert(orderer.GrpcTLSOpt.ClientKey, orderer.GrpcTLSOpt.ClientCrt),
			client.WithTLS(orderer.GrpcTLSOpt.Ca),
			client.WithTimeout(orderer.GrpcTLSOpt.Timeout),
		)
		if err != nil {
			// return nil, err
			log.Println(err)
			continue
		}
		defer order.Close()
		ordererClient, err := order.Broadcast()
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
		return nil, errors.Wrap(err, "error marshaling proposal")
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
		return nil, "", errors.Wrap(err, "failed to marshal args")
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
		return nil, "", errors.WithMessage(err, "failed to serialize identity")
	}

	proposal, txID, err := protoutil.CreateProposalFromCIS(common.HeaderType_ENDORSER_TRANSACTION, channel, cis, signerSerialized)
	if err != nil {
		return nil, "", errors.WithMessage(err, "failed to create ChaincodeInvocationSpec proposal")
	}

	return proposal, txID, nil
}
