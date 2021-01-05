package lifecycle

import (
	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/peer"
	lifecycleb "github.com/hyperledger/fabric-protos-go/peer/lifecycle"
)

// QueryApproved query approved chaincode
// chainOpt just need Name , Sequence default last or specific number
// peerGrpcOpt Timeout is necessary
// channelID fabric channel name
// peerAddress peer address
func QueryApproved(
	chainOpt chaincode.ChainOpt,
	//peerGrpcOpt GrpcTLSOpt,
	mspOpt chaincode.MSPOpt,
	channelID string,
	peer []chaincode.Endpoint,
) (*peer.ProposalResponse, error) {
	var args proto.Message

	function := "QueryApprovedChaincodeDefinition"
	args = &lifecycleb.QueryApprovedChaincodeDefinitionArgs{
		Name:     chainOpt.Name,
		Sequence: chainOpt.Sequence,
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

// QueryApproved2 query approved chaincode
// opt2 peer Grpc tls setting by path
// others -> QueryApproved
func QueryApproved2(
	opt chaincode.ChainOpt,
	//opt2 GrpcTLSOpt2,
	mspOpt chaincode.MSPOpt,
	channelID string,
	peer []chaincode.Endpoint2,
) (*peer.ProposalResponse, error) {
	//grpc, err := GrpcTLSOpt2ToGrpcTLSOpt(opt2)
	//if err != nil {
	//	return nil, err
	//}

	ep, err := chaincode.Endpoint2sToEndpoints(peer)
	if err != nil {
		return nil, err
	}
	return QueryApproved(opt, mspOpt, channelID, ep)
}
