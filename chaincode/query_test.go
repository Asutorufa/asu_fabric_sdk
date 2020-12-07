package chaincode

import (
	"fmt"
	"testing"
	"time"
)

func get(t *testing.T, a string) {
	resp, err := Query2(
		ChainOpt{Path: "sacc", Name: "sacc", IsInit: true, Version: "1.0.4"},
		GrpcTLSOpt2{
			ClientCrtPath:      "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/client.crt",
			ClientKeyPath:      "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/client.key",
			CaPath:             "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/ca.crt",
			ServerNameOverride: "peer-0-baas98",
			Timeout:            6 * time.Second,
		},
		MSPOpt{
			Path: "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/peers/peer-0-baas98/msp",
			Id:   "baas98",
		},
		[][]byte{[]byte("get"), []byte(a)},
		"channel1",
		[]string{"192.168.9.196:30060"},
	)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)
	fmt.Println(resp.Response.Status, string(resp.Response.Payload))
}

func get2(t *testing.T, a string) {
	resp, err := Query2(
		ChainOpt{Path: "sacc", Name: "sacc", IsInit: false, Version: "3.0"},
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
		[][]byte{[]byte("get"), []byte(a)},
		"channel1",
		[]string{"127.0.0.1:7051"},
	)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)
	fmt.Println(resp.Response.Status, string(resp.Response.Payload))
}

func TestQuery2(t *testing.T) {
	get2(t, "a")
	get2(t, "b")
}

func TestQuery(t *testing.T) {
	get(t, "a")
	get(t, "b")
}
