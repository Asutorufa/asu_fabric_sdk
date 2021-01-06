package lifecycle

import (
	"fmt"
	"io/ioutil"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-protos-go/peer/lifecycle"
)

// Install install a chaincode
// chainOpt need: path optional: others
func Install(
	chainOpt chaincode.ChainOpt,
	mspOpt chaincode.MSPOpt,
	peers []chaincode.Endpoint,
) (*peer.ProposalResponse, error) {
	pkgBytes, err := ioutil.ReadFile(chainOpt.Path)
	if err != nil {
		return nil, fmt.Errorf("read chaincode package from [%s] error -> %v", chainOpt.Path, err)
	}

	signer, err := chaincode.GetSigner(mspOpt.Path, mspOpt.Id)
	if err != nil {
		return nil, fmt.Errorf("get signer error -> %v", err)
	}

	args := &lifecycle.InstallChaincodeArgs{
		ChaincodeInstallPackage: pkgBytes,
	}

	proposal, _, err := createProposal(args, signer, "InstallChaincode", "")
	if err != nil {
		return nil, fmt.Errorf("create proposal error -> %v", err)
	}

	return query(signer, proposal, peers)
}

// Install2 to Install
func Install2(
	chainOpt chaincode.ChainOpt,
	mspOpt chaincode.MSPOpt,
	peers []chaincode.Endpoint2,
) (*peer.ProposalResponse, error) {
	p, err := chaincode.Endpoint2sToEndpoints(peers)
	if err != nil {
		return nil, fmt.Errorf("endpoint2s to endpoint error -> %v", err)
	}

	return Install(chainOpt, mspOpt, p)
}
