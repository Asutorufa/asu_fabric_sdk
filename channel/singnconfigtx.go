package channel

import (
	"fmt"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/hyperledger/fabric/protoutil"
)

//SignConfigTx sign config tx file to file
func SignConfigTx(channelID string, txFile []byte, mspOpt chaincode.MSPOpt) ([]byte, error) {
	env, err := protoutil.UnmarshalEnvelope(txFile)
	if err != nil {
		return nil, fmt.Errorf("unmarshalEnvelope Failed: %v", err)
	}

	signer, err := chaincode.GetSigner(mspOpt.Path, mspOpt.ID)
	if err != nil {
		return nil, fmt.Errorf("get signer failed: %v", err)
	}

	sTxEnv, err := sanityCheckAndSignConfigTx(channelID, env, signer)
	if err != nil {
		return nil, fmt.Errorf("sign config tx failed: %v", err)
	}

	return protoutil.Marshal(sTxEnv)
}
