package lifecycle

import (
	"fmt"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-protos-go/peer/lifecycle"
	"github.com/hyperledger/fabric/common/policydsl"
)

// Commit commit a chaincode
// chainOpt need: name,version,sequence optional: others
func Commit(
	chainOpt chaincode.ChainOpt,
	mspOpt chaincode.MSPOpt,
	channelID string,
	peers []chaincode.Endpoint,
	orderers []chaincode.Endpoint,
) (*peer.ProposalResponse, error) {
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

	signer, err := chaincode.GetSigner(mspOpt.Path, mspOpt.Id)
	if err != nil {
		return nil, fmt.Errorf("get signer error -> %v", err)
	}

	args := &lifecycle.CommitChaincodeDefinitionArgs{
		Name:                chainOpt.Name,
		Version:             chainOpt.Version,
		Sequence:            chainOpt.Sequence,
		EndorsementPlugin:   chainOpt.EndorsementPlugin,
		ValidationPlugin:    chainOpt.ValidationPlugin,
		ValidationParameter: chainOpt.ValidationParameter,
		InitRequired:        chainOpt.IsInit,
		Collections:         collections,
	}

	proposal, txID, err := createProposal(args, signer, commitFuncName, channelID)
	if err != nil {
		return nil, fmt.Errorf("create proposal error -> %v", err)
	}

	return invoke(signer, proposal, peers, orderers, channelID, txID)
}

// Commit2 to Commit
func Commit2(
	chainOpt chaincode.ChainOpt,
	mspOpt chaincode.MSPOpt,
	channelID string,
	peers []chaincode.Endpoint2,
	orderers []chaincode.Endpoint2,
) (*peer.ProposalResponse, error) {
	p, err := chaincode.Endpoint2sToEndpoints(peers)
	if err != nil {
		return nil, fmt.Errorf("peers' endpoint2s to endpoints error -> %v", err)
	}
	o, err := chaincode.Endpoint2sToEndpoints(orderers)
	if err != nil {
		return nil, fmt.Errorf("orderers' endpoint2s to endpoint error -> %v", err)
	}

	return Commit(chainOpt, mspOpt, channelID, p, o)
}
