package channel

import (
	"fmt"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/core/scc/cscc"
)

//GetChannels get all channels from a peer
func GetChannels(mspOpt chaincode.MSPOpt, peers chaincode.Endpoint) ([]*peer.ChannelInfo, error) {
	resp, err := exec(mspOpt, peers, &peer.ChaincodeSpec{
		Type:        peer.ChaincodeSpec_GOLANG,
		ChaincodeId: &peer.ChaincodeID{Name: "cscc"},
		Input: &peer.ChaincodeInput{
			Args: [][]byte{[]byte(cscc.GetChannels)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("exec to peer failed: %v", err)
	}

	q := &peer.ChannelQueryResponse{}

	err = proto.Unmarshal(resp.Payload, q)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %v", err)
	}

	return q.Channels, nil
}
