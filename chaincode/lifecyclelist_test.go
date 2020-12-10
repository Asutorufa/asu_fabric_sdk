package chaincode

import (
	"testing"
	"time"
)

func TestQueryInstalled(t *testing.T) {
	resp, err := QueryInstalled2(
		GrpcTLSOpt2{
			ClientCrtPath:      "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.crt",
			ClientKeyPath:      "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.key",
			CaPath:             "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/ca.crt",
			ServerNameOverride: "peer0.org1.example.com",
			Timeout:            6 * time.Second,
		},
		MSPOpt{
			Path: "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp",
			Id:   "Org1MSP",
		},
		"127.0.0.1:7051",
	)

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(resp.Response, resp.Response.Status, resp.Response.Message, string(resp.Response.Payload))
}

func TestQueryApproved(t *testing.T) {
	resp, err := QueryApproved2(
		ChainOpt{
			Name: "basic",
			//Sequence: 1,
		},
		GrpcTLSOpt2{
			ClientCrtPath:      "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.crt",
			ClientKeyPath:      "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.key",
			CaPath:             "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/ca.crt",
			ServerNameOverride: "peer0.org1.example.com",
			Timeout:            6 * time.Second,
		},
		MSPOpt{
			Path: "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp",
			Id:   "Org1MSP",
		},
		"mychannel",
		"127.0.0.1:7051",
	)

	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp)
	t.Log(resp.Response, resp.Response.Status, resp.Response.Message, string(resp.Response.Payload))
}

func TestQueryCommitted(t *testing.T) {
	resp, err := QueryCommitted2(
		ChainOpt{
			//Name: "basic",
		},
		GrpcTLSOpt2{
			ClientCrtPath:      "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.crt",
			ClientKeyPath:      "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.key",
			CaPath:             "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/ca.crt",
			ServerNameOverride: "peer0.org1.example.com",
			Timeout:            6 * time.Second,
		},
		MSPOpt{
			Path: "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp",
			Id:   "Org1MSP",
		},
		"mychannel",
		"127.0.0.1:7051",
	)

	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp)
	t.Log(resp.Response, resp.Response.Status, resp.Response.Message, string(resp.Response.Payload))
}
