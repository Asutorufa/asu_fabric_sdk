package lifecycle

import (
	"testing"
	"time"

	"github.com/Asutorufa/fabricsdk/chaincode"
)

func TestGetInstalledPackage(t *testing.T) {
	resp, err := GetInstalledPackage2(
		chaincode.ChainOpt{
			PackageID: "basic:794b463b3862b555435ae30621c1dc148780186b0755e3d797a3926a44dfd9b3",
		},
		chaincode.MSPOpt{
			Path: "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp",
			ID:   "Org1MSP",
		},
		[]chaincode.EndpointWithPath{
			{
				Address: "127.0.0.1:7051",
				GrpcTLSOptWithPath: chaincode.GrpcTLSOptWithPath{
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
	t.Log(resp.Response.Status, resp.Response.Message, len(resp.Response.Payload))
}
