package chaincode

import (
	"context"
	"errors"
	"fmt"

	"github.com/Asutorufa/fabricsdk/client"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/protoutil"
)

// Query2 .
func Query2(chaincode ChainOpt, mspOpt MSPOpt, args [][]byte, privateData map[string][]byte,
	channelID string, peers []EndpointWithPath) (*peer.ProposalResponse, error) {
	var peers2 []Endpoint

	for index := range peers {
		ep, err := ParseEndpointWithPath(peers[index])
		if err != nil {
			return nil, err
		}
		peers2 = append(peers2, ep)
	}

	return Query(chaincode, mspOpt, args, privateData, channelID, peers2)
}

// Query query from chaincode
// chaincode Path,Name,IsInit,Version,Type are necessary
// peerGrpcOpt Timeout is necessary
// mspOpt necessary
// args necessary
// privateData not necessary
// channelID necessary
// peerAddress necessary
func Query(chaincode ChainOpt, mspOpt MSPOpt, args [][]byte, privateData map[string][]byte,
	channelID string, peers []Endpoint) (*peer.ProposalResponse, error) {
	proposalResponse, err := query(chaincode, mspOpt, args, privateData, channelID, peers)
	if err != nil {
		return nil, err
	}

	resp := proposalResponse[0]

	if resp == nil {
		return nil, errors.New("received nil proposal response")
	}

	if resp.Response == nil {
		return nil, errors.New("received proposal response with nil response")
	}

	if resp.Response.Status != int32(common.Status_SUCCESS) {
		return nil, fmt.Errorf("query failed with status: %d - %s", resp.Response.Status, resp.Response.Message)
	}

	return resp, nil
}

func query(chaincode ChainOpt, mspOpt MSPOpt, args [][]byte, privateData map[string][]byte,
	channelID string, peers []Endpoint) ([]*peer.ProposalResponse, error) {
	peerClients := GetPeerClients(peers)
	if len(peerClients) == 0 {
		return nil, fmt.Errorf("peer clients' number is 0")
	}
	defer CloseClients(peerClients)
	return internalQuery(chaincode, mspOpt, args, privateData, channelID, peerClients)
}

func internalQuery(chaincode ChainOpt, mspOpt MSPOpt, args [][]byte,
	privateData map[string][]byte, channelID string,
	peers []*client.PeerClient) ([]*peer.ProposalResponse, error) {
	invocation := getChaincodeInvocationSpec(
		chaincode.Path,
		chaincode.Name,
		chaincode.IsInit,
		chaincode.Version,
		peer.ChaincodeSpec_GOLANG,
		args,
	)
	signer, err := GetSigner(mspOpt.Path, mspOpt.ID)
	if err != nil {
		return nil, fmt.Errorf("GetSigner() -> %v", err)
	}
	creator, err := signer.Serialize()
	if err != nil {
		return nil, fmt.Errorf("signer.Serialize() -> %v", err)
	}

	prop, txid, err := protoutil.CreateChaincodeProposalWithTxIDAndTransient(
		common.HeaderType_ENDORSER_TRANSACTION,
		channelID,
		invocation,
		creator,
		"",
		privateData, // <- 因为链码提案被存储在区块链上，
		// 不要把私有数据包含在链码提案中也是非常重要的。
		//在链码提案中有一个特殊的字段 transient，
		//可以用它把私有数据来从客户端（或者链码将用来生成私有数据的数据）传递给节点上的链码调用。
		//链码可以通过调用 GetTransient() API 来获取 transient 字段。
		//这个 transient 字段会从通道交易中被排除。
	)
	if err != nil {
		return nil, fmt.Errorf("protoutil.CreateChaincodeProposalWithTxIDAndTransient() -> %v", err)
	}
	fmt.Printf("txid: %s\n", txid)

	signedProp, err := protoutil.GetSignedProposal(prop, signer)
	if err != nil {
		return nil, fmt.Errorf("protoutil.GetSignedProposal() -> %v", err)
	}

	var proposalResponse []*peer.ProposalResponse
	for pi := range peers {

		endorserClient, err := peers[pi].Endorser()
		if err != nil {
			return nil, fmt.Errorf("peerClient.Endorser() -> %v", err)
		}

		resp, err := endorserClient.ProcessProposal(context.Background(), signedProp)
		if err != nil {
			fmt.Printf("get endorser from peer client failed: %v", err)
			continue
		}
		proposalResponse = append(proposalResponse, resp)
	}

	if len(proposalResponse) == 0 {
		return nil, errors.New("all peers process proposal failed")
	}

	return proposalResponse, nil

}
