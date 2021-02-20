package lifecycle

import (
	"testing"
	"time"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/peer/lifecycle"
)

func TestApproveForMyOrg(t *testing.T) {
	resp, err := ApproveForMyOrg2(
		chaincode.ChainOpt{
			Name:              "basic",
			Version:           "1.0",
			Sequence:          1,
			Label:             "basic",
			EndorsementPlugin: "escc",
			ValidationPlugin:  "vscc",
			PackageID:         "basic:794b463b3862b555435ae30621c1dc148780186b0755e3d797a3926a44dfd9b3",
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
		[]chaincode.Endpoint2{
			{
				Address: "127.0.0.1:7050",
				GrpcTLSOpt2: chaincode.GrpcTLSOpt2{
					CaPath:             "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/ordererOrganizations/example.com/users/Admin@example.com/tls/ca.crt",
					ServerNameOverride: "orderer.example.com",
					Timeout:            6 * time.Second,
				},
			},
		},
	)

	if err != nil {
		t.Error(err)
	}

	t.Log(string(resp.Payload))

	x := &lifecycle.ApproveChaincodeDefinitionForMyOrgResult{} // nothing
	err = proto.Unmarshal(resp.Payload, x)
	if err != nil {
		t.Error(err)
	}
}

func TestApproveForMyOrg2(t *testing.T) {
	resp, err := ApproveForMyOrg2(
		chaincode.ChainOpt{
			Name:              "basic",
			Version:           "1.0",
			Sequence:          1,
			Label:             "basic",
			EndorsementPlugin: "escc",
			ValidationPlugin:  "vscc",
			PackageID:         "basic:794b463b3862b555435ae30621c1dc148780186b0755e3d797a3926a44dfd9b3",
		},
		chaincode.MSPOpt{
			Path: "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp",
			Id:   "Org2MSP",
		},
		"mychannel",
		[]chaincode.Endpoint2{
			{
				Address: "127.0.0.1:9051",
				GrpcTLSOpt2: chaincode.GrpcTLSOpt2{
					CaPath:             "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/tls/ca.crt",
					ServerNameOverride: "peer0.org2.example.com",
					Timeout:            6 * time.Second,
				},
			},
		},
		[]chaincode.Endpoint2{
			{
				Address: "127.0.0.1:7050",
				GrpcTLSOpt2: chaincode.GrpcTLSOpt2{
					CaPath:             "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/ordererOrganizations/example.com/users/Admin@example.com/tls/ca.crt",
					ServerNameOverride: "orderer.example.com",
					Timeout:            6 * time.Second,
				},
			},
		},
	)

	if err != nil {
		t.Error(err)
	}

	t.Log(resp.Response.Status)

	x := &lifecycle.ApproveChaincodeDefinitionForMyOrgResult{} // nothing
	err = proto.Unmarshal(resp.Payload, x)
	if err != nil {
		t.Error(err)
	}
}
