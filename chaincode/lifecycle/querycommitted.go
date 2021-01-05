package lifecycle

import (
	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-protos-go/peer/lifecycle"
)

func QueryCommitted(
	chainOpt chaincode.ChainOpt,
	//peerGrpcTLSOpt GrpcTLSOpt,
	mspOpt chaincode.MSPOpt,
	channelID string,
	peer chaincode.Endpoint,
) (*peer.ProposalResponse, error) {
	var function string
	var args proto.Message
	if chainOpt.Name != "" {
		function = "QueryChaincodeDefinition"
		args = &lifecycle.QueryChaincodeDefinitionArgs{
			Name: chainOpt.Name,
		}
	} else {
		function = "QueryChaincodeDefinitions"
		args = &lifecycle.QueryChaincodeDefinitionsArgs{}
	}

	signer, err := chaincode.GetSigner(mspOpt.Path, mspOpt.Id)
	if err != nil {
		return nil, err
	}

	proposal, err := createProposal(args, signer, function, channelID)
	if err != nil {
		return nil, err
	}

	return query(signer, proposal, peer)
}

func QueryCommitted2(
	chainOpt chaincode.ChainOpt,
	mspOpt chaincode.MSPOpt,
	channelID string,
	peer chaincode.Endpoint2,
) (*peer.ProposalResponse, error) {
	ep, err := chaincode.Endpoint2ToEndpoint(peer)
	if err != nil {
		return nil, err
	}
	return QueryCommitted(chainOpt, mspOpt, channelID, ep)
}
