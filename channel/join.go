package channel

import (
	"context"
	"fmt"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/Asutorufa/fabricsdk/chaincode/client/clientcommon"
	"github.com/Asutorufa/fabricsdk/chaincode/client/peerclient"
	pcommon "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/core/scc/cscc"
	"github.com/hyperledger/fabric/protoutil"
)

func getJoinCCSPec(genesisBlock []byte) *peer.ChaincodeSpec {
	return &peer.ChaincodeSpec{
		Type: peer.ChaincodeSpec_GOLANG,
		ChaincodeId: &peer.ChaincodeID{
			Name: "cscc",
		},
		Input: &peer.ChaincodeInput{
			Args: [][]byte{
				[]byte(cscc.JoinChain),
				genesisBlock,
			},
		},
	}
}

//Join join a channel
func Join(mspOpt chaincode.MSPOpt, peers chaincode.Endpoint, genesisBlock []byte) (*peer.ProposalResponse, error) {
	return join(mspOpt, peers, getJoinCCSPec(genesisBlock))
}

func join(mspOpt chaincode.MSPOpt, peers chaincode.Endpoint, ccSpec *peer.ChaincodeSpec) (*peer.ProposalResponse, error) {
	signer, err := chaincode.GetSigner(mspOpt.Path, mspOpt.Id)
	if err != nil {
		return nil, fmt.Errorf("get signer error -> %v", err)
	}

	creator, err := signer.Serialize()
	if err != nil {
		return nil, fmt.Errorf("signer Serialize error -> %v", err)
	}

	prop, _, err := protoutil.CreateProposalFromCIS(
		pcommon.HeaderType_CONFIG,
		"",
		&peer.ChaincodeInvocationSpec{
			ChaincodeSpec: ccSpec,
		},
		creator,
	)
	if err != nil {
		return nil, fmt.Errorf("create proposal error -> %v", err)
	}

	signedProp, err := protoutil.GetSignedProposal(prop, signer)
	if err != nil {
		return nil, fmt.Errorf("signed proposal error -> %v", err)
	}

	peerClient, err := peerclient.NewPeerClient(
		peers.Address,
		peers.ServerNameOverride,
		clientcommon.WithTLS(peers.Ca),
		clientcommon.WithClientCert(peers.ClientKey, peers.ClientCrt),
		clientcommon.WithTimeout(peers.Timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("get new peer [%s] client error -> %v", peers.Address, err)
	}

	endorser, err := peerClient.Endorser()
	if err != nil {
		return nil, fmt.Errorf("get endorser error -> %v", err)
	}

	proposalResp, err := endorser.ProcessProposal(context.Background(), signedProp)
	if err != nil {
		return nil, fmt.Errorf("endorser process proposal error -> %v", err)
	}
	if proposalResp == nil {
		return nil, fmt.Errorf("nil proposal response")
	}
	if proposalResp.Response.Status != 0 && proposalResp.Response.Status != 200 {
		return nil, fmt.Errorf("bad proposal response %d: %s", proposalResp.Response.Status, proposalResp.Response.Message)
	}
	return proposalResp, nil
}
