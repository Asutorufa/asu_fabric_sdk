package deliverclient

import (
	"fabricSDK/chaincode/peerclient"

	"github.com/hyperledger/fabric-protos-go/peer"
)

func NewDeliverClient(peer *peerclient.PeerClient) (peer.DeliverClient, error) {
	return peer.PeerDeliver()
}
