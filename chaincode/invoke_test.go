package chaincode

import (
	"testing"
	"time"
)

func set(t *testing.T, a, b string) {
	resp, err := Invoke2(
		ChainOpt{Path: "sacc", Name: "sacc", IsInit: true, Version: "1.0.4"},

		MSPOpt{
			Path: "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/peers/peer-0-baas98/msp",
			Id:   "baas98",
		},
		[][]byte{[]byte("set"), []byte(a), []byte(b)},
		map[string][]byte{},
		"channel1",
		[]Endpoint2{
			{
				Address: "192.168.9.196:30060",
				GrpcTLSOpt2: GrpcTLSOpt2{
					ClientCrtPath:      "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/client.crt",
					ClientKeyPath:      "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/client.key",
					CaPath:             "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/ca.crt",
					ServerNameOverride: "peer-0-baas98",
					Timeout:            6 * time.Second,
				},
			},
		},
		Endpoint2{
			Address: "192.168.9.196:30062",
			GrpcTLSOpt2: GrpcTLSOpt2{
				ClientCrtPath:      "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/client.crt",
				ClientKeyPath:      "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/client.key",
				CaPath:             "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/ca.crt",
				ServerNameOverride: "orderer-0-baas98",
				Timeout:            6 * time.Second,
			},
		},
	)

	if err != nil {
		t.Error(err)
	}

	t.Log(resp)
}

func TestInvoke(t *testing.T) {
	set(t, "a", "xiaoxiao")
	set(t, "b", "xiaoxiao2")
}

func set2(t *testing.T, a, b string) {
	resp, err := Invoke2(
		ChainOpt{Path: "sacc", Name: "sacc", IsInit: false, Version: "3.0"},

		MSPOpt{
			Path: "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp",
			Id:   "Org1MSP",
		},
		[][]byte{[]byte("set"), []byte(a), []byte(b)},
		map[string][]byte{},
		"channel1",
		[]Endpoint2{
			{
				Address: "127.0.0.1:7051",
				GrpcTLSOpt2: GrpcTLSOpt2{
					ClientCrtPath:      "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.crt",
					ClientKeyPath:      "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.key",
					CaPath:             "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/ca.crt",
					ServerNameOverride: "peer0.org1.example.com",
					Timeout:            6 * time.Second,
				},
			},
		},
		Endpoint2{
			Address: "127.0.0.1:7050",
			GrpcTLSOpt2: GrpcTLSOpt2{
				ClientCrtPath:      "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/ordererOrganizations/example.com/users/Admin@example.com/tls/client.crt",
				ClientKeyPath:      "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/ordererOrganizations/example.com/users/Admin@example.com/tls/client.key",
				CaPath:             "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/ordererOrganizations/example.com/users/Admin@example.com/tls/ca.crt",
				ServerNameOverride: "orderer.example.com",
				Timeout:            6 * time.Second,
			},
		},
	)

	if err != nil {
		t.Error(err)
	}

	t.Log(resp)
}

func TestSet2(t *testing.T) {
	set2(t, "b", "d")
}
