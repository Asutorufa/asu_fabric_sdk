package chaincode

import (
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/core/chaincode/platforms/golang"
)

func Package() {

}

func NormalizePath(tYPE peer.ChaincodeSpec_Type, path string) (string, error) {
	switch tYPE {
	case peer.ChaincodeSpec_GOLANG:
		platform := &golang.Platform{}
		return platform.NormalizePath(path)
	case peer.ChaincodeSpec_NODE:
	case peer.ChaincodeSpec_CAR:
	case peer.ChaincodeSpec_JAVA:
	}
	return path, nil
}
