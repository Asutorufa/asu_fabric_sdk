package channel

import (
	"fmt"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/core/scc/qscc"
)

//GetChannelInfo get channel infos
func GetChannelInfo(channelID string, mspOpt chaincode.MSPOpt, peers chaincode.Endpoint) (*common.BlockchainInfo, error) {
	resp, err := exec(mspOpt, peers, &peer.ChaincodeSpec{
		Type:        peer.ChaincodeSpec_GOLANG,
		ChaincodeId: &peer.ChaincodeID{Name: "qscc"},
		Input: &peer.ChaincodeInput{
			Args: [][]byte{[]byte(qscc.GetChainInfo), []byte(channelID)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("exec to peer failed: %v", err)
	}

	i := &common.BlockchainInfo{}

	err = proto.Unmarshal(resp.Payload, i)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %v", err)
	}

	return i, nil
}
