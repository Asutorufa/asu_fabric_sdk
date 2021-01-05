package lifecycle

import (
	"fmt"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/hyperledger/fabric-protos-go/peer"
	lb "github.com/hyperledger/fabric-protos-go/peer/lifecycle"
	"github.com/hyperledger/fabric/common/policydsl"
)

func ApproveForMyOrg(
	chainOpt chaincode.ChainOpt,
	mspOpt chaincode.MSPOpt,
	channelID string,
	peers []chaincode.Endpoint,
) (*peer.ProposalResponse, error) {
	signer, err := chaincode.GetSigner(mspOpt.Path, mspOpt.Id)
	if err != nil {
		return nil, fmt.Errorf("get signer [mspPath:%s, mspID:%s] error -> %v", mspOpt.Path, mspOpt.Id, err)
	}

	var ccsrc *lb.ChaincodeSource
	if chainOpt.PackageID == "" {
		ccsrc = &lb.ChaincodeSource{
			Type: &lb.ChaincodeSource_LocalPackage{
				LocalPackage: &lb.ChaincodeSource_Local{
					PackageId: "",
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

	var collections *peer.CollectionConfigPackage

	for i := range chainOpt.CollectionsConfig {
		var ep *peer.ApplicationPolicy
		if chainOpt.CollectionsConfig[i].SignaturePolicy != "" &&
			chainOpt.CollectionsConfig[i].ChannelConfigPolicy != "" {
			return nil, fmt.Errorf("must spcify only one policy both SignaturePolicy and ChannelConfigPolicy")
		}
		if chainOpt.CollectionsConfig[i].SignaturePolicy != "" {
			p, err := policydsl.FromString(chainOpt.CollectionsConfig[i].SignaturePolicy)
			if err != nil {
				return nil, fmt.Errorf("format policy error -> %v", err)
			}

			ep = &peer.ApplicationPolicy{
				Type: &peer.ApplicationPolicy_SignaturePolicy{
					SignaturePolicy: p,
				},
			}
		}
		if chainOpt.CollectionsConfig[i].ChannelConfigPolicy != "" {
			ep = &peer.ApplicationPolicy{
				Type: &peer.ApplicationPolicy_ChannelConfigPolicyReference{
					ChannelConfigPolicyReference: chainOpt.CollectionsConfig[i].ChannelConfigPolicy,
				},
			}
		}
		p, err := policydsl.FromString(chainOpt.CollectionsConfig[i].Policy)
		if err != nil {
			return nil, fmt.Errorf("policy string error -> %v", err)
		}

		cc := &peer.CollectionConfig{
			Payload: &peer.CollectionConfig_StaticCollectionConfig{
				StaticCollectionConfig: &peer.StaticCollectionConfig{
					Name: chainOpt.CollectionsConfig[i].Name,
					MemberOrgsPolicy: &peer.CollectionPolicyConfig{
						Payload: &peer.CollectionPolicyConfig_SignaturePolicy{
							SignaturePolicy: p,
						},
					},
					RequiredPeerCount: chainOpt.CollectionsConfig[i].RequiredPeerCount,
					MaximumPeerCount:  chainOpt.CollectionsConfig[i].MaxPeerCount,
					BlockToLive:       chainOpt.CollectionsConfig[i].BlockToLive,
					MemberOnlyRead:    chainOpt.CollectionsConfig[i].MemberOnlyRead,
					MemberOnlyWrite:   chainOpt.CollectionsConfig[i].MemberOnlyWrite,
					EndorsementPolicy: ep,
				},
			},
		}

		collections.Config = append(collections.Config, cc)
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

	proposal, err := createProposal(args, signer, approveFuncName, channelID)
	if err != nil {
		return nil, fmt.Errorf("crate proposal error -> %v", err)
	}

	return query(signer, proposal, peers)
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
