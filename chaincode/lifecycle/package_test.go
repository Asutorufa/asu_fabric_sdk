package lifecycle

import (
	"testing"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/hyperledger/fabric-protos-go/peer"
)

func TestPackage(t *testing.T) {
	t.Log(
		Package2(chaincode.ChainOpt{
			Path:  "/mnt/shareSSD/code/Fabric/first/blockchainsTest/",
			Type:  peer.ChaincodeSpec_GOLANG,
			Label: "Test Package",
		}, "/mnt/shareSSD/code/Fabric/first/test.tar.gz"),
	)
}
