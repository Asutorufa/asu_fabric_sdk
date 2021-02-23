package lifecycle

import (
	"fmt"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/hyperledger/fabric-protos-go/peer"
	lb "github.com/hyperledger/fabric-protos-go/peer/lifecycle"
)

// ApproveForMyOrg approve for my org
// chainOpt -> need: Name,Version,Sequence optional: others
func ApproveForMyOrg(
	chainOpt chaincode.ChainOpt,
	mspOpt chaincode.MSPOpt,
	channelID string,
	peers []chaincode.Endpoint,
	orderers []chaincode.Endpoint,
) (*peer.ProposalResponse, error) {
	signer, err := chaincode.GetSigner(mspOpt.Path, mspOpt.ID)
	if err != nil {
		return nil, fmt.Errorf("get signer [mspPath:%s, mspID:%s] error -> %v", mspOpt.Path, mspOpt.ID, err)
	}

	var ccsrc *lb.ChaincodeSource
	if chainOpt.PackageID != "" {
		ccsrc = &lb.ChaincodeSource{
			Type: &lb.ChaincodeSource_LocalPackage{
				LocalPackage: &lb.ChaincodeSource_Local{
					PackageId: chainOpt.PackageID,
				},
			},
		}
	} else {
		ccsrc = &lb.ChaincodeSource{
			Type: &lb.ChaincodeSource_Unavailable_{
				Unavailable: &lb.ChaincodeSource_Unavailable{},
			},
		}
	}

	collections, err := chaincode.ConvertCollectionConfig(chainOpt.CollectionsConfig)
	if err != nil {
		return nil, fmt.Errorf("convert collections config failed: %v", err)
	}

	args := &lb.ApproveChaincodeDefinitionForMyOrgArgs{
		Name:                chainOpt.Name,
		Version:             chainOpt.Version,
		Sequence:            chainOpt.Sequence,
		EndorsementPlugin:   chainOpt.EndorsementPlugin,
		ValidationPlugin:    chainOpt.ValidationPlugin,
		ValidationParameter: chainOpt.ValidationParameter,
		InitRequired:        chainOpt.IsInit,
		Collections:         collections,
		Source:              ccsrc,
	}

	proposal, txID, err := createProposal(args, signer, approveFuncName, channelID)
	if err != nil {
		return nil, fmt.Errorf("crate proposal error -> %v", err)
	}

	return invoke(signer, proposal, peers, orderers, channelID, txID)
}

// ApproveForMyOrg2 to ApproveForMyOrg
func ApproveForMyOrg2(
	chainOpt chaincode.ChainOpt,
	mspOpt chaincode.MSPOpt,
	channelID string,
	peers []chaincode.EndpointWithPath,
	orderers []chaincode.EndpointWithPath,
) (*peer.ProposalResponse, error) {
	p, err := chaincode.ParseEndpointsWithPath(peers)
	if err != nil {
		return nil, fmt.Errorf("peers' endpoint2s to endpoints error -> %v", err)
	}
	o, err := chaincode.ParseEndpointsWithPath(orderers)
	if err != nil {
		return nil, fmt.Errorf("orderers' endpoint2s to endpoint error -> %v", err)
	}
	return ApproveForMyOrg(chainOpt, mspOpt, channelID, p, o)
}

/**
private data collection https://hyperledger-fabric.readthedocs.io/en/release-2.2/private_data_tutorial.html
// collections_config.json

[
   {
   "name": "assetCollection",
   "policy": "OR('Org1MSP.member', 'Org2MSP.member')",
   "requiredPeerCount": 1,
   "maxPeerCount": 1,
   "blockToLive":1000000,
   "memberOnlyRead": true,
   "memberOnlyWrite": true
   },
   {
   "name": "Org1MSPPrivateCollection",
   "policy": "OR('Org1MSP.member')",
   "requiredPeerCount": 0,
   "maxPeerCount": 1,
   "blockToLive":3,
   "memberOnlyRead": true,
   "memberOnlyWrite": false,
   "endorsementPolicy": {
       "signaturePolicy": "OR('Org1MSP.member')"
   }
   },
   {
   "name": "Org2MSPPrivateCollection",
   "policy": "OR('Org2MSP.member')",
   "requiredPeerCount": 0,
   "maxPeerCount": 1,
   "blockToLive":3,
   "memberOnlyRead": true,
   "memberOnlyWrite": false,
   "endorsementPolicy": {
       "signaturePolicy": "OR('Org2MSP.member')"
   }
   }
]
*/
