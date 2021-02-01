package lifecycle

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Asutorufa/fabricsdk/chaincode"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/core/chaincode/platforms/golang"
	"github.com/hyperledger/fabric/core/chaincode/platforms/java"
	"github.com/hyperledger/fabric/core/chaincode/platforms/node"
)

// Package package a chaincode
// chainOpt need: path, type, label
func Package(
	chainOpt chaincode.ChainOpt,
) ([]byte, error) {

	payload := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(payload)
	tw := tar.NewWriter(gw)

	normalizePath, err := NormalizePath(chainOpt.Type, chainOpt.Path)
	if err != nil {
		return nil, err
	}

	metadataBytes := []byte(fmt.Sprintf(`{"path":"%s","type":"%s","label":"%s"}`, normalizePath, chainOpt.Type.String(), chainOpt.Label))
	err = tw.WriteHeader(
		&tar.Header{
			Name: "metadata.json",
			Size: int64(len(metadataBytes)),
			Mode: 0100644,
		},
	)
	if err != nil {
		return nil, err
	}

	_, err = tw.Write(metadataBytes)
	if err != nil {
		return nil, err
	}

	codeBytes, err := GetDeploymentPayload(chainOpt.Type, chainOpt.Path)
	if err != nil {
		return nil, err
	}

	err = tw.WriteHeader(
		&tar.Header{
			Name: "code.tar.gz",
			Size: int64(len(codeBytes)),
			Mode: 0100644,
		},
	)
	if err != nil {
		return nil, err
	}

	_, err = tw.Write(codeBytes)
	if err != nil {
		return nil, err
	}

	err = tw.Close()
	if err == nil {
		err = gw.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("create tar error -> %v", err)
	}

	return payload.Bytes(), err
}

// Package2 to Package
func Package2(
	chainOpt chaincode.ChainOpt,
	outPutFile string,
) error {
	data, err := Package(chainOpt)
	if err != nil {
		return fmt.Errorf("get bytes error -> %v", err)
	}

	return ioutil.WriteFile(outPutFile, data, os.ModePerm)
}

// NormalizePath get path for different language chaincode
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

// GetDeploymentPayload get chaincode data from path for different language chaincode
func GetDeploymentPayload(tYPE peer.ChaincodeSpec_Type, path string) ([]byte, error) {
	switch tYPE {
	case peer.ChaincodeSpec_GOLANG:
		platform := &golang.Platform{}
		return platform.GetDeploymentPayload(path)
	case peer.ChaincodeSpec_NODE:
		platform := &node.Platform{}
		return platform.GetDeploymentPayload(path)
	case peer.ChaincodeSpec_JAVA:
		platform := &java.Platform{}
		return platform.GetDeploymentPayload(path)
	case peer.ChaincodeSpec_CAR:
	}

	return nil, fmt.Errorf("unsupport package platform -> %v", tYPE.String())
}
