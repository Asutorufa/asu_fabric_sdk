package chaincode

import (
	"context"
	"fmt"
	"log"

	"github.com/Asutorufa/fabricsdk/client"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/protoutil"
)

//ListInstantiated list in use chaincodes
func ListInstantiated(channelID string, mspOpt MSPOpt, peers []Endpoint) (*peer.ProposalResponse, error) {
	return InternalListInstantiated(channelID, mspOpt, GetPeerClients(peers))
}

//InternalListInstantiated list in use chaincodes
func InternalListInstantiated(channelID string, mspOpt MSPOpt, peers []*client.PeerClient) (*peer.ProposalResponse, error) {
	signer, err := GetSigner(mspOpt.Path, mspOpt.ID)
	if err != nil {
		return nil, fmt.Errorf("get signer failed: %v", err)
	}
	creator, err := signer.Serialize()
	if err != nil {
		return nil, fmt.Errorf("get creator failed: %v", err)
	}
	proposal, _, err := protoutil.CreateGetChaincodesProposal(channelID, creator)
	if err != nil {
		return nil, fmt.Errorf("create get chaincodes proposal failed: %v", err)
	}

	signedProposal, err := protoutil.GetSignedProposal(proposal, signer)
	if err != nil {
		return nil, fmt.Errorf("get signed proposal failed: %v", err)
	}
	for pi := range peers {
		endorser, err := peers[pi].Endorser()
		if err != nil {
			log.Printf("get endorser failed: %v\n", err)
			continue
		}

		proposalResponse, err := endorser.ProcessProposal(context.Background(), signedProposal)
		if err != nil {
			log.Printf("process proposal failed: %v\n", err)
		}

		return proposalResponse, nil
	}

	return nil, fmt.Errorf("no peers can process proposal")
}
