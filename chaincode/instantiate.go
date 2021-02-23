package chaincode

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Asutorufa/fabricsdk/client"
	"github.com/golang/protobuf/proto"
	protcommon "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/policydsl"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/protoutil"
	"github.com/pkg/errors"
)

//Instantiate init a chaincode
func Instantiate(channelID string, cTor string,
	chainOpt ChainOpt, signer msp.SigningIdentity,
	peerClient client.PeerClient, ordererClients []client.OrdererClient) error {
	env, err := getSingedTx(channelID, cTor, chainOpt, signer, peerClient)
	if err != nil {
		return fmt.Errorf("get signed tx failed: %v", err)
	}
	if env == nil {
		return nil
	}

	for oi := range ordererClients {
		broadcast, err := ordererClients[oi].Broadcast()
		if err != nil {
			log.Printf("get broadcast failed: %v\n", err)
			continue
		}

		err = broadcast.Send(env)
		if err != nil {
			log.Printf("broadcast send envelop failed: %v", err)
			continue
		}

		return nil
	}
	return errors.New("broadcast transaction failed")
}

func getSingedTx(channelID string, cTor string,
	chainOpt ChainOpt, signer msp.SigningIdentity,
	peerClient client.PeerClient) (*protcommon.Envelope, error) {
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
	if err != nil {
		return nil, fmt.Errorf("get creator failed: %v", err)
	}

	var policyMarshalled []byte
	if chainOpt.Policy != "" {
		p, err := policydsl.FromString(chainOpt.Policy)
		if err != nil {
			return nil, errors.Errorf("invalid policy %s", chainOpt.Policy)
		}
		policyMarshalled = protoutil.MarshalOrPanic(p)
	}

	if chainOpt.EndorsementPlugin != "" {
		log.Printf("Using escc %s\n", chainOpt.EndorsementPlugin)
	} else {
		log.Println("Using default escc")
		chainOpt.EndorsementPlugin = "escc"
	}

	if chainOpt.ValidationPlugin != "" {
		log.Printf("Using vscc %s\n", chainOpt.ValidationPlugin)
	} else {
		log.Println("Using default vscc")
		chainOpt.ValidationPlugin = "vscc"
	}

	collections, err := ConvertCollectionConfig(chainOpt.CollectionsConfig)
	if err != nil {
		return nil, fmt.Errorf("convert collection failed: %v", err)
	}

	collectionsByte, err := proto.Marshal(collections)
	if err != nil {
		return nil, fmt.Errorf("marshal collections config failed: %v", err)
	}

	prop, _, err := protoutil.CreateDeployProposalFromCDS(
		channelID, deployMent, creator, policyMarshalled,
		[]byte(chainOpt.EndorsementPlugin), []byte(chainOpt.ValidationPlugin), collectionsByte)
	if err != nil {
		return nil, fmt.Errorf("create deploy proposal failed: %v", err)
	}

	signedProposal, err := protoutil.GetSignedProposal(prop, signer)
	if err != nil {
		return nil, fmt.Errorf("sign proposal failed: %v", err)
	}

	endorser, err := peerClient.Endorser()
	if err != nil {
		return nil, fmt.Errorf("get endorser failed: %v", err)
	}

	resp, err := endorser.ProcessProposal(context.Background(), signedProposal)
	if err != nil {
		return nil, fmt.Errorf("process proposal failed: %v", err)
	}

	if resp != nil {
		env, err := protoutil.CreateSignedTx(prop, signer, resp)
		if err != nil {
			return nil, fmt.Errorf("create signed tx failed: %v", err)
		}

		return env, nil
	}

	return nil, nil
}
