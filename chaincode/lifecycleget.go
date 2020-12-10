package chaincode

import (
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-protos-go/peer/lifecycle"
)

// GetInstalledPackage
// chainOpt just need PackageID
func GetInstalledPackage(
	chainOpt ChainOpt,
	peerGrpcTLSOpt GrpcTLSOpt,
	mspOpt MSPOpt,
	peerAddress string,
) (*peer.ProposalResponse, error) {
	signer, err := GetSigner(mspOpt.Path, mspOpt.Id)
	args := &lifecycle.GetInstalledChaincodePackageArgs{
		PackageId: chainOpt.PackageID,
	}

	function := "GetInstalledChaincodePackage"

	proposal, err := createProposal(args, signer, function, "")
	if err != nil {

	}

	return query(peerGrpcTLSOpt, signer, proposal, peerAddress)
}

func GetInstalledPackage2(
	chainOpt ChainOpt,
	peerGrpcTLSOpt GrpcTLSOpt2,
	mspOpt MSPOpt,
	peerAddress string,
) (*peer.ProposalResponse, error) {
	grpc, err := GrpcTLSOpt2ToGrpcTLSOpt(peerGrpcTLSOpt)
	if err != nil {
		return nil, err
	}

	return GetInstalledPackage(chainOpt, grpc, mspOpt, peerAddress)
}
