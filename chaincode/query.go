package chaincode

import (
	"fmt"

	"github.com/Asutorufa/fabricsdk/chaincode/client/clientcommon"
	"github.com/Asutorufa/fabricsdk/chaincode/client/peerclient"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/protoutil"
)

func Query2(
	chaincode ChainOpt,
	mspOpt MSPOpt,
	args [][]byte,
	privateData map[string][]byte,
	channelID string,
	peers []Endpoint2,
) (*peer.ProposalResponse, error) {
	var peers2 []Endpoint

	for index := range peers {
		ep, err := Endpoint2ToEndpoint(peers[index])
		if err != nil {
			return nil, err
		}
		peers2 = append(peers2, ep)
	}

	return Query(
		chaincode,
		mspOpt,
		args,
		privateData,
		channelID,
		//"",
		peers2,
	)
}

// Query
// chaincode Path,Name,IsInit,Version,Type are necessary
// peerGrpcOpt Timeout is necessary
// mspOpt necessary
// args necessary
// privateData not necessary
// channelID necessary
// peerAddress necessary
func Query(
	chaincode ChainOpt,
	mspOpt MSPOpt,
	args [][]byte,
	privateData map[string][]byte,
	channelID string,
	//txID string,
	peers []Endpoint,
) (*peer.ProposalResponse, error) {
	invocation := getChaincodeInvocationSpec(
		chaincode.Path,
		chaincode.Name,
		chaincode.IsInit,
		chaincode.Version,
		peer.ChaincodeSpec_GOLANG,
		args,
	)
	signer, err := GetSigner(mspOpt.Path, mspOpt.Id)
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

	signedProp, err := protoutil.GetSignedProposal(prop, signer)
	if err != nil {
		return nil, fmt.Errorf("protoutil.GetSignedProposal() -> %v", err)
	}

	var peerClients []*peerclient.PeerClient
	var endorserClients []peer.EndorserClient
	for index := range peers {
		peerClient, err := peerclient.NewPeerClient(
			peers[index].Address,
			peers[index].GrpcTLSOpt.ServerNameOverride,
			clientcommon.WithTLS(peers[index].GrpcTLSOpt.Ca),
			clientcommon.WithClientCert(peers[index].GrpcTLSOpt.ClientKey, peers[index].GrpcTLSOpt.ClientCrt),
			clientcommon.WithTimeout(peers[index].GrpcTLSOpt.Timeout),
		)
		if err != nil {
			return nil, fmt.Errorf("NewPeerClient() -> %v", err)
		}
		peerClients = append(peerClients, peerClient)

		endorserClient, err := peerClient.Endorser()
		if err != nil {
			return nil, fmt.Errorf("peerClient.Endorser() -> %v", err)
		}
		endorserClients = append(endorserClients, endorserClient)
	}
	defer func() {
		for index := range peerClients {
			_ = peerClients[index].Close()
		}
	}()

	responses, err := processProposals(endorserClients, signedProp)
	if err != nil {
		return nil, fmt.Errorf("processProposals() -> %v", err)
	}
	fmt.Printf("txid: %s\n", txid)
	return responses[0], nil
}
