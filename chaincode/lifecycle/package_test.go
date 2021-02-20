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
			Label: "basic", // label 不能有特殊字符 不然会安装失败
		}, "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/basic.tar.gz"),
	)
}
