package channel

import (
	"fmt"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/core/scc/cscc"
)

//JoinBySnapshotStatus get join by snapshot status
func JoinBySnapshotStatus(mspOpt chaincode.MSPOpt, peers chaincode.Endpoint) (*peer.JoinBySnapshotStatus, error) {
	spec := &peer.ChaincodeSpec{
		Type:        peer.ChaincodeSpec_GOLANG,
		ChaincodeId: &peer.ChaincodeID{Name: "cscc"},
		Input:       &peer.ChaincodeInput{Args: [][]byte{[]byte(cscc.JoinBySnapshotStatus)}},
	}

	resp, err := exec(mspOpt, peers, spec)
	if err != nil {
		return nil, fmt.Errorf("exec to peer failed: %v", err)
	}

	status := &peer.JoinBySnapshotStatus{}

	err = proto.Unmarshal(resp.Payload, status)
	if err != nil {
		return nil, fmt.Errorf("unmarshal proto failed: %v", err)
	}

	return status, nil
}
