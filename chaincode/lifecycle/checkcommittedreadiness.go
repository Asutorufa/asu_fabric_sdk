package lifecycle

import (
	"encoding/json"
	"fmt"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-protos-go/peer/lifecycle"
	"github.com/hyperledger/fabric/common/policydsl"
	"github.com/hyperledger/fabric/protoutil"
	"github.com/pkg/errors"
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

	var ccp *peer.CollectionConfigPackage
	if chainOpt.CollectionConfig != "" {
		ccp, _, err = getCollectionConfigFromBytes([]byte(chainOpt.CollectionConfig))
		if err != nil {
			return nil, err
		}
	}
	args := &lifecycle.CheckCommitReadinessArgs{
		Sequence:            chainOpt.Sequence,
		Name:                chainOpt.Name,
		Version:             chainOpt.Name,
		EndorsementPlugin:   chainOpt.EndorsementPlugin,
		ValidationPlugin:    chainOpt.ValidationPlugin,
		ValidationParameter: protoutil.MarshalOrPanic(applicationPolicy),
		Collections:         ccp,
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

type collectionConfigJson struct {
	Name              string `json:"name"`
	Policy            string `json:"policy"`
	RequiredPeerCount *int32 `json:"requiredPeerCount"`
	MaxPeerCount      *int32 `json:"maxPeerCount"`
	BlockToLive       uint64 `json:"blockToLive"`
	MemberOnlyRead    bool   `json:"memberOnlyRead"`
	MemberOnlyWrite   bool   `json:"memberOnlyWrite"`
	EndorsementPolicy *struct {
		SignaturePolicy     string `json:"signaturePolicy"`
		ChannelConfigPolicy string `json:"channelConfigPolicy"`
	} `json:"endorsementPolicy"`
}

// getCollectionConfig retrieves the collection configuration
// from the supplied byte array; the byte array must contain a
// json-formatted array of collectionConfigJson elements
func getCollectionConfigFromBytes(cconfBytes []byte) (*peer.CollectionConfigPackage, []byte, error) {
	cconf := &[]collectionConfigJson{}
	err := json.Unmarshal(cconfBytes, cconf)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not parse the collection configuration")
	}

	ccarray := make([]*peer.CollectionConfig, 0, len(*cconf))
	for _, cconfitem := range *cconf {
		p, err := policydsl.FromString(cconfitem.Policy)
		if err != nil {
			return nil, nil, errors.WithMessagef(err, "invalid policy %s", cconfitem.Policy)
		}

		cpc := &peer.CollectionPolicyConfig{
			Payload: &peer.CollectionPolicyConfig_SignaturePolicy{
				SignaturePolicy: p,
			},
		}

		var ep *peer.ApplicationPolicy
		if cconfitem.EndorsementPolicy != nil {
			signaturePolicy := cconfitem.EndorsementPolicy.SignaturePolicy
			channelConfigPolicy := cconfitem.EndorsementPolicy.ChannelConfigPolicy
			if (signaturePolicy != "" && channelConfigPolicy != "") || (signaturePolicy == "" && channelConfigPolicy == "") {
				return nil, nil, fmt.Errorf("incorrect policy")
			}

			if signaturePolicy != "" {
				poli, err := policydsl.FromString(signaturePolicy)
				if err != nil {
					return nil, nil, err
				}
				ep = &peer.ApplicationPolicy{
					Type: &peer.ApplicationPolicy_SignaturePolicy{
						SignaturePolicy: poli,
					},
				}
			} else {
				ep = &peer.ApplicationPolicy{
					Type: &peer.ApplicationPolicy_ChannelConfigPolicyReference{
						ChannelConfigPolicyReference: channelConfigPolicy,
					},
				}
			}
		}

		// Set default requiredPeerCount and MaxPeerCount if not specified in json
		requiredPeerCount := int32(0)
		maxPeerCount := int32(1)
		if cconfitem.RequiredPeerCount != nil {
			requiredPeerCount = *cconfitem.RequiredPeerCount
		}
		if cconfitem.MaxPeerCount != nil {
			maxPeerCount = *cconfitem.MaxPeerCount
		}

		cc := &peer.CollectionConfig{
			Payload: &peer.CollectionConfig_StaticCollectionConfig{
				StaticCollectionConfig: &peer.StaticCollectionConfig{
					Name:              cconfitem.Name,
					MemberOrgsPolicy:  cpc,
					RequiredPeerCount: requiredPeerCount,
					MaximumPeerCount:  maxPeerCount,
					BlockToLive:       cconfitem.BlockToLive,
					MemberOnlyRead:    cconfitem.MemberOnlyRead,
					MemberOnlyWrite:   cconfitem.MemberOnlyWrite,
					EndorsementPolicy: ep,
				},
			},
		}

		ccarray = append(ccarray, cc)
	}

	ccp := &peer.CollectionConfigPackage{Config: ccarray}
	ccpBytes, err := proto.Marshal(ccp)
	return ccp, ccpBytes, err
}
