package lifecycle

import (
	"context"
	"fmt"
	"time"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/Asutorufa/fabricsdk/chaincode/client/clientcommon"
	"github.com/Asutorufa/fabricsdk/chaincode/client/peerclient"
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

func query(
	//peerGrpcOpt GrpcTLSOpt,
	signer msp.SigningIdentity,
	proposal *peer.Proposal,
	peer chaincode.Endpoint,
) (*peer.ProposalResponse, error) {

	signedProposal, err := signProposal(proposal, signer)
	if err != nil {
		return nil, err
	}

	peerClient, err := peerclient.NewPeerClient(
		peer.Address,
		peer.ServerNameOverride,
		clientcommon.WithClientCert2(peer.ClientKey, peer.ClientCrt),
		clientcommon.WithTLS2(peer.Ca),
		clientcommon.WithTimeout(6*time.Second),
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

	if resp == nil {
		return nil, fmt.Errorf("resp is nil")
	}

	if resp.Response.Status != int32(common.Status_SUCCESS) {
		return nil, fmt.Errorf("%d - %s", resp.Response.Status, resp.Response.Message)
	}

	return resp, nil
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

func createProposal(args proto.Message, signer msp.SigningIdentity, function, channel string) (*peer.Proposal, error) {
	argsBytes, err := proto.Marshal(args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal args")
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
		return nil, errors.WithMessage(err, "failed to serialize identity")
	}

	proposal, _, err := protoutil.CreateProposalFromCIS(common.HeaderType_ENDORSER_TRANSACTION, channel, cis, signerSerialized)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create ChaincodeInvocationSpec proposal")
	}

	return proposal, nil
}
