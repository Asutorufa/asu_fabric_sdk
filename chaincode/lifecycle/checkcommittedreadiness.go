package lifecycle

import (
	"fmt"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-protos-go/peer/lifecycle"
	"github.com/hyperledger/fabric/common/policydsl"
	"github.com/hyperledger/fabric/protoutil"
)

// CheckCommittedReadiness check committed read readiness
// chainOpt need: Name, Sequence, Policy, optional: others
func CheckCommittedReadiness(
	chainOpt chaincode.ChainOpt,
	mspOpt chaincode.MSPOpt,
	channelID string,
	pEER []chaincode.Endpoint,
) (*peer.ProposalResponse, error) {
	signaturePolicyEnvelope, err := policydsl.FromString(chainOpt.Policy)
	if err != nil {
		return nil, err
	}

	applicationPolicy := &peer.ApplicationPolicy{
		Type: &peer.ApplicationPolicy_SignaturePolicy{
			SignaturePolicy: signaturePolicyEnvelope,
		},
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

	// var ccp *peer.CollectionConfigPackage
	// if chainOpt.CollectionConfig != "" {
	// 	ccp, _, err = getCollectionConfigFromBytes([]byte(chainOpt.CollectionConfig))
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	args := &lifecycle.CheckCommitReadinessArgs{
		Sequence:            chainOpt.Sequence,
		Name:                chainOpt.Name,
		Version:             chainOpt.Name,
		EndorsementPlugin:   chainOpt.EndorsementPlugin,
		ValidationPlugin:    chainOpt.ValidationPlugin,
		ValidationParameter: protoutil.MarshalOrPanic(applicationPolicy),
		Collections:         collections,
		InitRequired:        chainOpt.IsInit,
	}

	signer, err := chaincode.GetSigner(mspOpt.Path, mspOpt.Id)
	if err != nil {
		return nil, err
	}

	proposal, _, err := createProposal(args, signer, checkCommitReadinessFuncName, channelID)
	if err != nil {
		return nil, err
	}

	resp, err := query(signer, proposal, pEER)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CheckCommittedReadiness2 to CheckCommittedReadiness
func CheckCommittedReadiness2(
	chainOpt chaincode.ChainOpt,
	mspOpt chaincode.MSPOpt,
	channelID string,
	peer []chaincode.Endpoint2,
) (*peer.ProposalResponse, error) {
	ep, err := chaincode.Endpoint2sToEndpoints(peer)
	if err != nil {
		return nil, err
	}
	return CheckCommittedReadiness(chainOpt, mspOpt, channelID, ep)
}

//[
//{
//"name": "collectionMarbles",
//"policy": "OR('Org1MSP.member', 'Org2MSP.member')",
//"requiredPeerCount": 0,
//"maxPeerCount": 3,
//"blockToLive":1000000,
//"memberOnlyRead": true,
//"memberOnlyWrite": true
//},
//{
//"name": "collectionMarblePrivateDetails",
//"policy": "OR('Org1MSP.member')",
//"requiredPeerCount": 0,
//"maxPeerCount": 3,
//"blockToLive":3,
//"memberOnlyRead": true,
//"memberOnlyWrite":true,
//"endorsementPolicy": {
//"signaturePolicy": "OR('Org1MSP.member')"
//}
//}
//]
