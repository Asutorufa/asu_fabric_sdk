package chaincode

import (
	"testing"
	"time"

	"github.com/hyperledger/fabric-protos-go/peer"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/peer/lifecycle"
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

	s := &lifecycle.QueryApprovedChaincodeDefinitionResult{}
	err = proto.UnmarshalMerge(resp.Response.Payload, s)
	if err != nil {
		t.Error(err)
	}
	t.Log(s)
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

	s := &lifecycle.QueryChaincodeDefinitionsResult{}
	err = proto.UnmarshalMerge(resp.Response.Payload, s)
	if err != nil {
		t.Error(err)
	}
	t.Log(s)
}

func TestCheckCommittedReadiness(t *testing.T) {
	resp, err := CheckCommittedReadiness2(
		ChainOpt{
			Path:              "assetTransfer",
			Name:              "basic",
			IsInit:            true,
			Version:           "1.0",
			PackageID:         "basic_1.0:4ec191e793b27e953ff2ede5a8bcc63152cecb1e4c3f301a26e22692c61967ad",
			Sequence:          1,
			EndorsementPlugin: "escc",
			ValidationPlugin:  "vscc",
			Policy:            "OR('Org1MSP.member')",
			Type:              peer.ChaincodeSpec_GOLANG,
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
	t.Log(resp.Response, resp.Response.Status, resp.Response.Message)

	s := lifecycle.CheckCommitReadinessResult{}

	err = proto.UnmarshalMerge(resp.Response.Payload, &s)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(s.Approvals)
}
