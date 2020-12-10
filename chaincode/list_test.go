package chaincode

import (
	"fmt"
	"testing"
	"time"
)

func TestGetInstalled(t *testing.T) {
	resp, err := GetInstalled2(
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
		[]string{"127.0.0.1:7051"},
	)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)
	fmt.Println(resp.Response.Status, string(resp.Response.Payload))
}

func TestGetInstantiated(t *testing.T) {
	resp, err := GetInstantiated2(
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
		[]string{"127.0.0.1:7051"},
		"mychannel",
	)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)
	fmt.Println(resp.Response.Status, resp.GetResponse(), string(resp.Response.Payload), resp.Response.GetMessage())
}
