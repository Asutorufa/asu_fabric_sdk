package channel

import (
	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/core/scc/cscc"
)

// JoinBySnapshot join channel by snapshot
func JoinBySnapshot(mspOpt chaincode.MSPOpt, peers chaincode.Endpoint, snapshotPath string) (*peer.ProposalResponse, error) {
	return join(mspOpt, peers, &peer.ChaincodeSpec{
		Type: peer.ChaincodeSpec_GOLANG,
		ChaincodeId: &peer.ChaincodeID{
			Name: "cscc",
		},
		Input: &peer.ChaincodeInput{
			Args: [][]byte{[]byte(cscc.JoinChainBySnapshot), []byte(snapshotPath)},
		},
	})
}
