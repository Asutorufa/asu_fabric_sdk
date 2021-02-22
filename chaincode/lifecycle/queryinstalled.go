package lifecycle

import (
	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-protos-go/peer/lifecycle"
)

// QueryInstalled2 query installed chaincode
func QueryInstalled2(
	mspOpt chaincode.MSPOpt,
	peer []chaincode.EndpointWithPath,
) (*peer.ProposalResponse, error) {
	ep, err := chaincode.ParseEndpointsWithPath(peer)
	if err != nil {
		return nil, err
	}
	return QueryInstalled(mspOpt, ep)
}

// QueryInstalled query installed chaincode
func QueryInstalled(
	mspOpt chaincode.MSPOpt,
	peer []chaincode.Endpoint,
) (*peer.ProposalResponse, error) {
	signer, err := chaincode.GetSigner(mspOpt.Path, mspOpt.ID)
	if err != nil {
		return nil, err
	}

	proposal, _, err := createProposal(&lifecycle.QueryInstalledChaincodeArgs{}, signer, "QueryInstalledChaincodes", "")
	if err != nil {
		return nil, err
	}

	return query(signer, proposal, peer)
}
