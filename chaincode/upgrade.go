package chaincode

import (
	"github.com/Asutorufa/fabricsdk/client"
	"github.com/hyperledger/fabric/msp"
)

//Upgrade update a chaincode
func Upgrade(channelID string, cTor string,
	chainOpt ChainOpt, signer msp.SigningIdentity,
	peerClient client.PeerClient, ordererClients []client.OrdererClient) error {
	return Instantiate(channelID, cTor, chainOpt, signer, peerClient, ordererClients)
}
