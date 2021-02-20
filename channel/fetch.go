package channel

import (
	"fmt"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/Asutorufa/fabricsdk/client"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/hyperledger/fabric/common/util"
)

// Fetch fetch specific block from orderer
func Fetch(mspOpt chaincode.MSPOpt, orderers chaincode.Endpoint, channelID string, blockNum uint64) (*common.Block, error) {
	ordererClient, err := client.NewOrdererClientSelf(
		orderers.Address,
		orderers.ServerNameOverride,
		client.WithClientCert(orderers.ClientKey, orderers.ClientCrt),
		client.WithTLS(orderers.Ca),
		client.WithTimeout(orderers.Timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("get orderer [%s] client error -> %v", orderers.Address, err)
	}

	deliver, err := ordererClient.Deliver()
	if err != nil {
		return nil, fmt.Errorf("get orderer client deliver error -> %v", err)
	}

	signer, err := chaincode.GetSigner(mspOpt.Path, mspOpt.Id)
	if err != nil {
		return nil, fmt.Errorf("get signer error -> %v", err)
	}

	var tlsCertHash []byte
	if len(ordererClient.Certificate().Certificate) > 0 {
		tlsCertHash = util.ComputeSHA256(ordererClient.Certificate().Certificate[0])
	}

	env, err := getBlockEnvelop(channelID, &orderer.SeekPosition{
		Type: &orderer.SeekPosition_Specified{
			Specified: &orderer.SeekSpecified{
				Number: blockNum,
			},
		},
	}, tlsCertHash, signer, true)
	if err != nil {
		return nil, fmt.Errorf("get block envelop error -> %v", err)
	}

	err = deliver.Send(env)
	if err != nil {
		return nil, fmt.Errorf("deliver send error -> %v", err)
	}

	resp, err := deliver.Recv()
	if err != nil {
		return nil, fmt.Errorf("recv from deliver error -> %v", err)
	}

	switch resp.Type.(type) {
	case *orderer.DeliverResponse_Status:
		return nil, fmt.Errorf("Expect block, but get status: %v", resp)
	case *orderer.DeliverResponse_Block:
		resp, err = deliver.Recv()
		if err != nil {
			return nil, fmt.Errorf("recv from deliver error -> %v", err)
		}
		if resp.GetStatus() != common.Status_SUCCESS {
			return nil, fmt.Errorf("response status code [%d] is not successful", resp.GetStatus())
		}
		return resp.GetBlock(), nil
	default:
		return nil, fmt.Errorf("unknown type: %T", resp)
	}
}
