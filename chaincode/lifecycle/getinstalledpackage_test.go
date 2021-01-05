package lifecycle

import (
	"testing"
	"time"

	"github.com/Asutorufa/fabricsdk/chaincode"
)

func TestGetInstalledPackage(t *testing.T) {
	resp, err := GetInstalledPackage2(
		chaincode.ChainOpt{
			PackageID: "basic_1.0:4ec191e793b27e953ff2ede5a8bcc63152cecb1e4c3f301a26e22692c61967ad",
		},
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
	t.Log(resp.Response.Status, resp.Response.Message, len(resp.Response.Payload))
}
