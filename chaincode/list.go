package chaincode

import (
	"context"
	"fabricSDK/chaincode/client/clientcommon"
	"fabricSDK/chaincode/client/peerclient"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/protoutil"
)

func list(
	peerGrpcOpt GrpcTLSOpt,
	signer msp.SigningIdentity,
	peerAddress []string,
	proposal *peer.Proposal,
) (*peer.ProposalResponse, error) {
	signedProposal, err := protoutil.GetSignedProposal(proposal, signer)
	if err != nil {
		return nil, err
	}

	for index := range peerAddress {
		peerClient, err := peerclient.NewPeerClient(
			peerAddress[index],
			peerGrpcOpt.ServerNameOverride,
			clientcommon.WithClientCert2(peerGrpcOpt.ClientKey, peerGrpcOpt.ClientCrt),
			clientcommon.WithTLS2(peerGrpcOpt.Ca),
			clientcommon.WithTimeout(6*time.Second),
		)
		if err != nil {
			log.Println(err)
			continue
		}
		endorser, err := peerClient.Endorser()
		if err != nil {
			log.Println(err)
			continue
		}

		resp, err := endorser.ProcessProposal(context.Background(), signedProposal)
		if err != nil {
			return nil, err
		}

		if resp == nil {
			log.Printf("peer %s resp is nil", peerAddress[index])
			continue
		}

		if resp.Response.Status != int32(common.Status_SUCCESS) {
			log.Printf("%d - %s", resp.Response.Status, resp.Response.Message)
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("create peer client failed")
}

func GetInstalled(
	peerGrpcOpt GrpcTLSOpt,
	mspOpt MSPOpt,
	peerAddress []string,
) (*peer.ProposalResponse, error) {
	signer, err := GetSigner(mspOpt.Path, mspOpt.Id)
	if err != nil {
		return nil, err
	}

	creator, err := signer.Serialize()
	if err != nil {
		return nil, err
	}
	proposal, _, err := protoutil.CreateGetInstalledChaincodesProposal(creator)
	return list(peerGrpcOpt, signer, peerAddress, proposal)
}

func GetInstalled2(
	peerGrpcOpt GrpcTLSOpt2,
	mspOpt MSPOpt,
	peerAddress []string,
) (*peer.ProposalResponse, error) {
	grpc, err := GrpcTLSOpt2ToGrpcTLSOpt(peerGrpcOpt)
	if err != nil {
		return nil, err
	}

	return GetInstalled(grpc, mspOpt, peerAddress)
}

func GetInstantiated(
	peerGrpcOpt GrpcTLSOpt,
	mspOpt MSPOpt,
	peerAddress []string,
	channelID string,
) (*peer.ProposalResponse, error) {
	signer, err := GetSigner(mspOpt.Path, mspOpt.Id)
	if err != nil {
		return nil, err
	}

	creator, err := signer.Serialize()
	if err != nil {
		return nil, err
	}

	proposal, _, err := protoutil.CreateGetChaincodesProposal(channelID, creator)
	return list(peerGrpcOpt, signer, peerAddress, proposal)
}

func GetInstantiated2(
	peerGrpcOpt GrpcTLSOpt2,
	mspOpt MSPOpt,
	peerAddress []string,
	channelID string,
) (*peer.ProposalResponse, error) {
	grpc, err := GrpcTLSOpt2ToGrpcTLSOpt(peerGrpcOpt)
	if err != nil {
		return nil, err
	}

	return GetInstantiated(grpc, mspOpt, peerAddress, channelID)
}
