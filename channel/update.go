package channel

import (
	"fmt"
	"log"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/Asutorufa/fabricsdk/chaincode/client/clientcommon"
	"github.com/Asutorufa/fabricsdk/chaincode/client/orderclient"
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/common/configtx"
	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/protoutil"
)

// Update update channel config
func Update(channelID string, mspOpt chaincode.MSPOpt, orderers []chaincode.Endpoint, updateConfig []byte) error {
	ctxEnv, err := protoutil.UnmarshalEnvelope(updateConfig)
	if err != nil {
		return fmt.Errorf("unmarshal envelope error -> %v", err)
	}

	signer, err := chaincode.GetSigner(mspOpt.Path, mspOpt.Id)
	if err != nil {
		return fmt.Errorf("get msp signer error -> %v", err)
	}
	chCrtEnv, err := sanityCheckAndSignConfigTx(channelID, ctxEnv, signer)
	if err != nil {
		return fmt.Errorf("check envelop with error -> %v", err)
	}

	for oi := range orderers {
		ordererClient, err := orderclient.NewOrdererClient(
			orderers[oi].Address,
			orderers[oi].ServerNameOverride,
			clientcommon.WithClientCert(orderers[oi].ClientKey, orderers[oi].ClientCrt),
			clientcommon.WithTLS(orderers[oi].Ca),
		)
		if err != nil {
			log.Printf("initialize new orderer [%s] client error -> %v\n", orderers[oi].Address, err)
			continue
		}
		defer ordererClient.Close()

		bc, err := ordererClient.Broadcast()
		if err != nil {
			log.Printf("get orderer broad cast error -> %v\n", err)
			continue
		}
		err = bc.Send(chCrtEnv)
		if err != nil {
			log.Printf("send envelop error -> %v\n", err)
			continue
		}
		return nil
	}
	return fmt.Errorf("send envelop error")
}

// copy from github.com/hyperledger/fabric/internal/peer/channel/update.go
func sanityCheckAndSignConfigTx(channelID string, envConfigUpdate *cb.Envelope, signer msp.SigningIdentity) (*cb.Envelope, error) {
	payload, err := protoutil.UnmarshalPayload(envConfigUpdate.Payload)
	if err != nil {
		return nil, fmt.Errorf("bad payload")
	}

	if payload.Header == nil || payload.Header.ChannelHeader == nil {
		return nil, fmt.Errorf("bad header")
	}

	ch, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall channel header")
	}

	if ch.Type != int32(cb.HeaderType_CONFIG_UPDATE) {
		return nil, fmt.Errorf("bad type")
	}

	if ch.ChannelId == "" {
		return nil, fmt.Errorf("empty channel id")
	}

	// Specifying the chainID on the CLI is usually redundant, as a hack, set it
	// here if it has not been set explicitly
	if channelID == "" {
		channelID = ch.ChannelId
	}

	if ch.ChannelId != channelID {
		return nil, fmt.Errorf("mismatched channel ID %s != %s", ch.ChannelId, channelID)
	}

	configUpdateEnv, err := configtx.UnmarshalConfigUpdateEnvelope(payload.Data)
	if err != nil {
		return nil, fmt.Errorf("Bad config update env")
	}

	sigHeader, err := protoutil.NewSignatureHeader(signer)
	if err != nil {
		return nil, err
	}

	configSig := &cb.ConfigSignature{
		SignatureHeader: protoutil.MarshalOrPanic(sigHeader),
	}

	configSig.Signature, err = signer.Sign(util.ConcatenateBytes(configSig.SignatureHeader, configUpdateEnv.ConfigUpdate))
	if err != nil {
		return nil, err
	}

	configUpdateEnv.Signatures = append(configUpdateEnv.Signatures, configSig)

	return protoutil.CreateSignedEnvelope(cb.HeaderType_CONFIG_UPDATE, channelID, signer, configUpdateEnv, 0, 0)
}
