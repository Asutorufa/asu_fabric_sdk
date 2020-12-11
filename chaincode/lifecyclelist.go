package chaincode

import (
	"context"
	"encoding/json"
	"fabricSDK/chaincode/client/clientcommon"
	"fabricSDK/chaincode/client/peerclient"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/common/policydsl"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-protos-go/peer/lifecycle"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/protoutil"
	"github.com/pkg/errors"
)

const (
	lifecycleName                = "_lifecycle"
	approveFuncName              = "ApproveChaincodeDefinitionForMyOrg"
	commitFuncName               = "CommitChaincodeDefinition"
	checkCommitReadinessFuncName = "CheckCommitReadiness"
)

func QueryInstalled2(
	peerGrpcOpt GrpcTLSOpt2,
	mspOpt MSPOpt,
	peerAddress string,
) (*peer.ProposalResponse, error) {
	grpc, err := GrpcTLSOpt2ToGrpcTLSOpt(peerGrpcOpt)
	if err != nil {
		return nil, err
	}
	return QueryInstalled(grpc, mspOpt, peerAddress)
}

func QueryInstalled(
	peerGrpcOpt GrpcTLSOpt,
	mspOpt MSPOpt,
	peerAddress string,
) (*peer.ProposalResponse, error) {
	signer, err := GetSigner(mspOpt.Path, mspOpt.Id)
	if err != nil {
		return nil, err
	}

	proposal, err := createProposal(&lifecycle.QueryInstalledChaincodeArgs{}, signer, "QueryInstalledChaincodes", "")
	if err != nil {
		return nil, err
	}

	return query(peerGrpcOpt, signer, proposal, peerAddress)
}

func query(
	peerGrpcOpt GrpcTLSOpt,
	signer msp.SigningIdentity,
	proposal *peer.Proposal,
	peerAddress string,
) (*peer.ProposalResponse, error) {

	signedProposal, err := signProposal(proposal, signer)
	if err != nil {
		return nil, err
	}

	peerClient, err := peerclient.NewPeerClient(
		peerAddress,
		peerGrpcOpt.ServerNameOverride,
		clientcommon.WithClientCert2(peerGrpcOpt.ClientKey, peerGrpcOpt.ClientCrt),
		clientcommon.WithTLS2(peerGrpcOpt.Ca),
		clientcommon.WithTimeout(6*time.Second),
	)
	if err != nil {
		return nil, err
	}

	endorserClient, err := peerClient.Endorser()
	if err != nil {
		return nil, err
	}

	resp, err := endorserClient.ProcessProposal(context.Background(), signedProposal)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("resp is nil")
	}

	if resp.Response.Status != int32(common.Status_SUCCESS) {
		return nil, fmt.Errorf("%d - %s", resp.Response.Status, resp.Response.Message)
	}

	return resp, nil
}

func signProposal(proposal *peer.Proposal, signer msp.SigningIdentity) (*peer.SignedProposal, error) {
	// check for nil argument
	if proposal == nil {
		return nil, errors.New("proposal cannot be nil")
	}

	if signer == nil {
		return nil, errors.New("signer cannot be nil")
	}

	proposalBytes, err := proto.Marshal(proposal)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling proposal")
	}

	signature, err := signer.Sign(proposalBytes)
	if err != nil {
		return nil, err
	}

	return &peer.SignedProposal{
		ProposalBytes: proposalBytes,
		Signature:     signature,
	}, nil
}

// QueryApproved
// chainOpt just need Name , Sequence default last or specific number
// peerGrpcOpt Timeout is necessary
// channelID fabric channel name
// peerAddress peer address
func QueryApproved(
	chainOpt ChainOpt,
	peerGrpcOpt GrpcTLSOpt,
	mspOpt MSPOpt,
	channelID string,
	peerAddress string,
) (*peer.ProposalResponse, error) {
	var args proto.Message

	function := "QueryApprovedChaincodeDefinition"
	args = &lifecycle.QueryApprovedChaincodeDefinitionArgs{
		Name:     chainOpt.Name,
		Sequence: chainOpt.Sequence,
	}

	signer, err := GetSigner(mspOpt.Path, mspOpt.Id)
	if err != nil {
		return nil, err
	}

	proposal, err := createProposal(args, signer, function, channelID)
	if err != nil {
		return nil, err
	}

	return query(peerGrpcOpt, signer, proposal, peerAddress)
}

// QueryApproved2
// opt2 peer Grpc tls setting by path
// others -> QueryApproved
func QueryApproved2(
	opt ChainOpt,
	opt2 GrpcTLSOpt2,
	mspOpt MSPOpt,
	channelID string,
	peerAddress string,
) (*peer.ProposalResponse, error) {
	grpc, err := GrpcTLSOpt2ToGrpcTLSOpt(opt2)
	if err != nil {
		return nil, err
	}
	return QueryApproved(opt, grpc, mspOpt, channelID, peerAddress)
}

func createProposal(args proto.Message, signer msp.SigningIdentity, function, channel string) (*peer.Proposal, error) {
	argsBytes, err := proto.Marshal(args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal args")
	}
	ccInput := &peer.ChaincodeInput{Args: [][]byte{[]byte(function), argsBytes}}

	cis := &peer.ChaincodeInvocationSpec{
		ChaincodeSpec: &peer.ChaincodeSpec{
			ChaincodeId: &peer.ChaincodeID{Name: lifecycleName},
			Input:       ccInput,
		},
	}

	signerSerialized, err := signer.Serialize()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to serialize identity")
	}

	proposal, _, err := protoutil.CreateProposalFromCIS(common.HeaderType_ENDORSER_TRANSACTION, channel, cis, signerSerialized)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create ChaincodeInvocationSpec proposal")
	}

	return proposal, nil
}

func QueryCommitted(
	chainOpt ChainOpt,
	peerGrpcTLSOpt GrpcTLSOpt,
	mspOpt MSPOpt,
	channelID string,
	peerAddress string,
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

	signer, err := GetSigner(mspOpt.Path, mspOpt.Id)
	if err != nil {
		return nil, err
	}

	proposal, err := createProposal(args, signer, function, channelID)
	if err != nil {
		return nil, err
	}

	return query(peerGrpcTLSOpt, signer, proposal, peerAddress)
}

func QueryCommitted2(
	chainOpt ChainOpt,
	peerGrpcTLSOpt GrpcTLSOpt2,
	mspOpt MSPOpt,
	channelID string,
	peerAddress string,
) (*peer.ProposalResponse, error) {
	grpc, err := GrpcTLSOpt2ToGrpcTLSOpt(peerGrpcTLSOpt)
	if err != nil {
		return nil, err
	}

	return QueryCommitted(chainOpt, grpc, mspOpt, channelID, peerAddress)
}

func CheckCommittedReadiness2(
	chainOpt ChainOpt,
	peerGrpcTLSOpt GrpcTLSOpt2,
	mspOpt MSPOpt,
	channelID string,
	peerAddress string,
) (*peer.ProposalResponse, error) {
	grpc, err := GrpcTLSOpt2ToGrpcTLSOpt(peerGrpcTLSOpt)
	if err != nil {
		return nil, err
	}

	return CheckCommittedReadiness(chainOpt, grpc, mspOpt, channelID, peerAddress)
}

func CheckCommittedReadiness(
	chainOpt ChainOpt,
	peerGrpcTLSOpt GrpcTLSOpt,
	mspOpt MSPOpt,
	channelID string,
	peerAddress string,
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

	ccp, _, err := getCollectionConfigFromBytes([]byte(chainOpt.CollectionConfig))
	if err != nil {
		return nil, err
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

	signer, err := GetSigner(mspOpt.Path, mspOpt.Id)
	if err != nil {
		return nil, err
	}

	proposal, err := createProposal(args, signer, checkCommitReadinessFuncName, channelID)
	if err != nil {
		return nil, err
	}

	resp, err := query(peerGrpcTLSOpt, signer, proposal, peerAddress)
	if err != nil {
		return nil, err
	}

	return resp, nil
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
