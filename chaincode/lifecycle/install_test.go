package lifecycle

import (
	"testing"
	"time"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/peer/lifecycle"
)

func TestInstall(t *testing.T) {
	resp, err := Install2(
		chaincode.ChainOpt{ // this chaincode package must use lifecycle package
			Path: "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/basic.tar.gz",
		},
		chaincode.MSPOpt{
			Path: "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp",
			ID:   "Org1MSP",
		},
		[]chaincode.EndpointWithPath{
			{
				Address: "127.0.0.1:7051",
				GrpcTLSOptWithPath: chaincode.GrpcTLSOptWithPath{
					ClientCrtPath:      "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.crt",
					ClientKeyPath:      "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.key",
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

	x := &lifecycle.InstallChaincodeResult{}
	err = proto.Unmarshal(resp.Payload, x)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(x)
}

func TestInstall2(t *testing.T) {
	resp, err := Install2(
		chaincode.ChainOpt{ // this chaincode package must use lifecycle package
			Path: "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/basic.tar.gz",
		},
		chaincode.MSPOpt{
			Path: "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp",
			ID:   "Org2MSP",
		},
		[]chaincode.EndpointWithPath{
			{
				Address: "127.0.0.1:9051",
				GrpcTLSOptWithPath: chaincode.GrpcTLSOptWithPath{
					CaPath:             "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/tls/ca.crt",
					ServerNameOverride: "peer0.org2.example.com",
					Timeout:            6 * time.Second,
				},
			},
		},
	)

	if err != nil {
		t.Error(err)
		return
	}

	x := &lifecycle.InstallChaincodeResult{}
	err = proto.Unmarshal(resp.Payload, x)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(x)
}
