package endorserclient

import (
	"fabricSDK/chaincode/peerclient"

	"github.com/hyperledger/fabric-protos-go/peer"
)

func NewEndorserClient(client *peerclient.PeerClient) (peer.EndorserClient, error) {
	return client.Endorser()
}
