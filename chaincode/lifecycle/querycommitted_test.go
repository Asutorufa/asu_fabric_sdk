package lifecycle

import (
	"testing"
	"time"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/peer/lifecycle"
)

func TestQueryCommitted(t *testing.T) {
	name := "basic"
	resp, err := QueryCommitted2(
		chaincode.ChainOpt{
			Name: name,
		},
		chaincode.MSPOpt{
			Path: "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp",
			Id:   "Org1MSP",
		},
		"mychannel",
		[]chaincode.Endpoint2{
			{
				Address: "127.0.0.1:7051",
				GrpcTLSOpt2: chaincode.GrpcTLSOpt2{
					CaPath:             "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/ca.crt",
					ServerNameOverride: "peer0.org1.example.com",
					Timeout:            6 * time.Second,
				},
			},
		},
	)

	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp.Response.Status, resp.Response.Message)

	var s proto.Message
	if name == "" {
		s = &lifecycle.QueryChaincodeDefinitionsResult{}
	} else {
		s = &lifecycle.QueryChaincodeDefinitionResult{}
	}

	err = proto.UnmarshalMerge(resp.Response.Payload, s)
	if err != nil {
		t.Error(err)
	}
	t.Log(s)
}
