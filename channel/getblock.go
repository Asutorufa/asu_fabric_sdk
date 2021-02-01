package channel

import (
	"fmt"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/protoutil"
)

func getBlockEnvelop(
	channelID string,
	position *orderer.SeekPosition,
	tlsCertHash []byte,
	signer msp.SigningIdentity,
	bestEffort bool,
) (*common.Envelope, error) {
	seekInfo := &orderer.SeekInfo{
		Start:    position,
		Stop:     position,
		Behavior: orderer.SeekInfo_BLOCK_UNTIL_READY,
	}

	if bestEffort {
		seekInfo.ErrorResponse = orderer.SeekInfo_BEST_EFFORT
	}

	enveLOP, err := protoutil.CreateSignedEnvelopeWithTLSBinding(common.HeaderType_DELIVER_SEEK_INFO, channelID, signer, seekInfo, 0, 0, tlsCertHash)
	if err != nil {
		return nil, fmt.Errorf("create signed envelop with tls binding error -> %v", err)
	}

	return enveLOP, nil
}
