package channel

import (
	"fmt"
	"log"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/Asutorufa/fabricsdk/client"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/protoutil"
)

// Create create channel
func Create(channelID string, txFile []byte, mspOpt chaincode.MSPOpt, orderers []chaincode.Endpoint) (*common.Block, error) {
	env, err := getTxEnvelop(txFile)
	if err != nil {
		return nil, fmt.Errorf("get tx envelop failed: %v", err)
	}

	signer, err := chaincode.GetSigner(mspOpt.Path, mspOpt.ID)
	if err != nil {
		return nil, fmt.Errorf("get signer failed: %v", err)
	}

	signedEnv, err := sanityCheckAndSignConfigTx(channelID, env, signer)
	if err != nil {
		return nil, fmt.Errorf("signed envelop failed: %v", err)
	}

	for oi := range orderers {
		oc, err := client.NewOrdererClientSelf(
			orderers[oi].Address,
			orderers[oi].ServerNameOverride,
			client.WithClientCert(orderers[oi].ClientKey, orderers[oi].ClientCrt),
			client.WithTLS(orderers[oi].Ca),
		)
		if err != nil {
			log.Printf("create orderer [%s] client failed: %v\n", orderers[oi].Address, err)
			continue
		}
		defer oc.Close()

		bc, err := oc.Broadcast()
		if err != nil {
			log.Printf("get broadcast failed: %v\n", err)
			continue
		}

		err = bc.Send(signedEnv)
		if err != nil {
			log.Printf("send signed envelop failed: %v", err)
			continue
		}
		block, err := Fetch(mspOpt, orderers[oi], channelID, 0)
		if err != nil {
			log.Printf("fetch genesis block failed: %v", err)
			continue
		}
		return block, nil
	}
	return nil, fmt.Errorf("send signed envelop to all orderers failed")
}

func getTxEnvelop(file ...[]byte) (*common.Envelope, error) {
	l := len(file)

	if l > 1 {
		return nil, fmt.Errorf("only 0 or 1 file")
	}

	if l == 0 {
		return nil, fmt.Errorf("now not support create default channel envelop")
	}

	return protoutil.UnmarshalEnvelope(file[0])
}
