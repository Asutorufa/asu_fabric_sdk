package chaincode

import (
	"fmt"
	"testing"
	"time"
)

func get(t *testing.T, a string) {
	resp, err := Query2(
		ChainOpt{Path: "sacc", Name: "sacc", IsInit: true, Version: "1.0.4"},
		MSPOpt{
			Path: "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/peers/peer-0-baas98/msp",
			ID:   "baas98",
		},
		[][]byte{[]byte("get"), []byte(a)},
		map[string][]byte{},
		"channel1",
		[]EndpointWithPath{
			{
				Address: "192.168.9.196:30060",
				GrpcTLSOptWithPath: GrpcTLSOptWithPath{
					CaPath:             "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/ca.crt",
					ServerNameOverride: "peer-0-baas98",
					Timeout:            6 * time.Second},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)
	fmt.Println(resp.Response.Status, string(resp.Response.Payload))
}

func TestQuery(t *testing.T) {
	get(t, "a")
	get(t, "b")
}

func get2(t *testing.T, b [][]byte) {
	resp, err := Query2(
		ChainOpt{Path: "assetTransfer", Name: "basic", IsInit: false, Version: "1.0"},
		MSPOpt{
			Path: "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp",
			ID:   "Org1MSP",
		},
		b,
		map[string][]byte{},
		"mychannel",
		[]EndpointWithPath{
			{
				Address: "127.0.0.1:7051",
				GrpcTLSOptWithPath: GrpcTLSOptWithPath{
					CaPath:             "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/ca.crt",
					ServerNameOverride: "peer0.org1.example.com",
					Timeout:            6 * time.Second,
				},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}
	// fmt.Println(resp)
	fmt.Println("status code:", resp.Response.Status, "Payload: -> ", string(resp.Response.Payload))
}

func TestQuery2(t *testing.T) {
	get2(t, [][]byte{[]byte("get"), []byte("語彙")})
	fmt.Println()
	get2(t, [][]byte{[]byte("get2"), []byte("君は")})
	fmt.Println()
	get2(t, [][]byte{[]byte("getCreator")})
}
