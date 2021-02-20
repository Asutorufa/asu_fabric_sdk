package lifecycle

import (
	"testing"
	"time"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/peer/lifecycle"
)

func TestCheckCommittedReadiness(t *testing.T) {
	resp, err := CheckCommittedReadiness2(
		chaincode.ChainOpt{
			// Path:      "assetTransfer",
			Name: "basic",
			// IsInit:    true,
			// Version:   "1.0",
			PackageID: "basic:794b463b3862b555435ae30621c1dc148780186b0755e3d797a3926a44dfd9b3",
			Sequence:  1,
			// EndorsementPlugin: "escc",
			// ValidationPlugin:  "vscc",
			Policy: "OR('Org1MSP.member','Org2MSP.member')",
			// Type:             peer.ChaincodeSpec_GOLANG,
			// CollectionConfig: "",
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
	// t.Log(resp)
	// t.Log(resp.Response, resp.Response.Status, resp.Response.Message)

	s := lifecycle.CheckCommitReadinessResult{}

	err = proto.UnmarshalMerge(resp.Response.Payload, &s)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(s.Approvals)
}
