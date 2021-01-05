package lifecycle

import (
	"testing"
	"time"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/peer/lifecycle"
)

func TestQueryInstalled(t *testing.T) {
	resp, err := QueryInstalled2(
		chaincode.MSPOpt{
			Path: "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp",
			Id:   "Org1MSP",
		},
		chaincode.Endpoint2{
			Address: "127.0.0.1:7051",
			GrpcTLSOpt2: chaincode.GrpcTLSOpt2{
				ClientCrtPath:      "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.crt",
				ClientKeyPath:      "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.key",
				CaPath:             "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/ca.crt",
				ServerNameOverride: "peer0.org1.example.com",
				Timeout:            6 * time.Second,
			},
		},
	)

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(resp.Response.Status, resp.Response.Message)

	s := &lifecycle.QueryInstalledChaincodesResult{}
	err = proto.UnmarshalMerge(resp.Response.Payload, s)
	if err != nil {
		t.Error(err)
	}
	t.Log(s)
}
