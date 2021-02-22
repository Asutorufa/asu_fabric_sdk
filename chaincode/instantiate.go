package chaincode

import (
	"encoding/json"

	"github.com/Asutorufa/fabricsdk/client"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/protoutil"
	"github.com/pkg/errors"
)

func Instantiate(chainOpt ChainOpt, signer msp.SigningIdentity, peerClient []client.PeerClient, ordererClient []client.OrdererClient, cTor string) (*peer.ProposalResponse, error) {

	input := &peer.ChaincodeInput{}
	if err := json.Unmarshal([]byte(cTor), &input); err != nil {
		return nil, errors.Wrap(err, "chaincode argument error")
	}
	input.IsInit = chainOpt.IsInit

	spec := &peer.ChaincodeSpec{
		Type:        peer.ChaincodeSpec_Type(chainOpt.Type),
		ChaincodeId: &peer.ChaincodeID{Path: chainOpt.Path, Name: chainOpt.Name, Version: chainOpt.Version},
		Input:       input,
	}
	deployMent := &peer.ChaincodeDeploymentSpec{ChaincodeSpec: spec}

	creator, err := signer.Serialize()
	prop, _, err := protoutil.CreateDeployProposalFromCDS("", deployMent, creator, []byte{}, []byte{}, []byte{}, []byte{})
}
