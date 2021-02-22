package chaincode

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/Asutorufa/fabricsdk/client"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/core/chaincode/platforms/golang"
	"github.com/hyperledger/fabric/core/chaincode/platforms/java"
	"github.com/hyperledger/fabric/core/chaincode/platforms/node"
	"github.com/hyperledger/fabric/core/common/ccpackage"
	"github.com/hyperledger/fabric/core/common/ccprovider"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/protoutil"
)

//Install install a chaincode before fabric 2.0
//cTor eg: '{"Args":["init","a","100","b","200"]}'
// isPackage whether chainOpt.Path is a .tar.gz package or not
func Install(chainOpt ChainOpt, mspOpt MSPOpt, peers Endpoint, cTor string, isPackage bool) (proposalResponse *peer.ProposalResponse, err error) {
	peerClient, err := client.NewPeerClient(peers.Address, peers.ServerNameOverride, client.WithTLS(peers.Ca), client.WithClientCert(peers.ClientKey, peers.ClientCrt))
	if err != nil {
		return nil, fmt.Errorf("create new peer[%s] client failed: %v", peers.Address, err)
	}
	signer, err := GetSigner(mspOpt.Path, mspOpt.ID)
	if err != nil {
		return nil, fmt.Errorf("get signer from msp [id:%s,path:%s] failed: %v", mspOpt.ID, mspOpt.Path, err)
	}

	return InternalInstall(chainOpt, signer, peerClient, cTor, isPackage)
}

//InternalInstall install a chaincode
func InternalInstall(chainOpt ChainOpt, signer msp.SigningIdentity, peerClient *client.PeerClient, cTor string, isPackage bool) (proposalResponse *peer.ProposalResponse, err error) {
	deploymentPayload, err := getDeploymentPayload(chainOpt, cTor, isPackage)
	if err != nil {
		return nil, fmt.Errorf("get deployment failed: %v", err)
	}

	creator, err := signer.Serialize()
	if err != nil {
		return nil, fmt.Errorf("signer serialize failed: %v", err)
	}

	proposal, _, err := protoutil.CreateInstallProposalFromCDS(deploymentPayload, creator)
	if err != nil {
		return nil, fmt.Errorf("create install proposal failed: %v", err)
	}

	signedProposal, err := protoutil.GetSignedProposal(proposal, signer)
	if err != nil {
		return nil, fmt.Errorf("signed proposal failed: %v", err)
	}

	endorser, err := peerClient.Endorser()
	if err != nil {
		return nil, fmt.Errorf("get endorser from peer client failed: %v", err)
	}

	resp, err := endorser.ProcessProposal(context.Background(), signedProposal)
	if err != nil {
		return nil, fmt.Errorf("endorser process proposal failed: %v", err)
	}

	if resp == nil {
		return nil, errors.New("error during install: received nil proposal response")
	}

	if resp.Response == nil {
		return nil, errors.New("error during install: received proposal response with nil response")
	}

	if resp.Response.Status != int32(common.Status_SUCCESS) {
		return nil, fmt.Errorf("install failed with status: %d - %s", resp.Response.Status, resp.Response.Message)
	}

	return resp, nil
}

func getDeploymentPayload(chainOpt ChainOpt, cTor string, isPackage bool) (
	deployment *peer.ChaincodeDeploymentSpec, err error) {
	if !isPackage {
		input := &peer.ChaincodeInput{}

		err = proto.Unmarshal([]byte(cTor), input)
		if err != nil {
			return nil, fmt.Errorf("unmarshal ctor failed: %v", err)
		}

		spec := &peer.ChaincodeSpec{
			Type: chainOpt.Type,
			ChaincodeId: &peer.ChaincodeID{
				Path:    chainOpt.Path,
				Name:    chainOpt.Name,
				Version: chainOpt.Version,
			},
			Input: input,
		}

		code, err := GetDeploymentPayload(chainOpt.Type, chainOpt.Path)
		if err != nil {
			return nil, fmt.Errorf("get deployment payload failed: %v", err)
		}
		return &peer.ChaincodeDeploymentSpec{ // deploymentSPec
			ChaincodeSpec: spec,
			CodePackage:   code,
		}, nil
	}

	pkgBytes, err := ioutil.ReadFile(chainOpt.Path)
	if err != nil {
		return nil, fmt.Errorf("read chaincode package file failed: %v", err)
	}

	ccpack, err := ccprovider.GetCCPackage(pkgBytes, factory.GetDefault())
	if err != nil {
		return nil, fmt.Errorf("get chaincode package failed: %v", err)
	}
	// either CDS or Envelope
	o := ccpack.GetPackageObject()

	// try CDS first
	cds, ok := o.(*peer.ChaincodeDeploymentSpec)
	if !ok || cds == nil {
		// try Envelope next
		env, ok := o.(*common.Envelope)
		if !ok || env == nil {
			return nil, errors.New("error extracting valid chaincode package")
		}

		// this will check for a valid package Envelope
		_, sCDS, err := ccpackage.ExtractSignedCCDepSpec(env)
		if err != nil {
			return nil, fmt.Errorf("extracting valid signed chaincode package failed: %v", err)
		}

		// ...and get the CDS at last
		cds, err = protoutil.UnmarshalChaincodeDeploymentSpec(sCDS.ChaincodeDeploymentSpec)
		if err != nil {
			return nil, fmt.Errorf("extracting chaincode deployment spec failed: %v", err)
		}

		err = validateDeploymentSpec(cds.ChaincodeSpec.Type, cds.CodePackage)
		if err != nil {
			return nil, fmt.Errorf("chaincode deployment spec validation failed: %v", err)
		}
	}

	// get the chaincode details from cds
	cName := cds.ChaincodeSpec.ChaincodeId.Name
	cVersion := cds.ChaincodeSpec.ChaincodeId.Version

	// if user provided chaincodeName, use it for validation
	if chainOpt.Name != "" && chainOpt.Name != cName {
		return nil, fmt.Errorf("chaincode name %s does not match name %s in package", chainOpt.Name, cName)
	}

	// if user provided chaincodeVersion, use it for validation
	if chainOpt.Version != "" && chainOpt.Version != cVersion {
		return nil, fmt.Errorf("chaincode version %s does not match version %s in packages", chainOpt.Version, cVersion)
	}

	return cds, nil
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

func validateDeploymentSpec(tYPE peer.ChaincodeSpec_Type, ccPkg []byte) error {

	switch tYPE {
	case peer.ChaincodeSpec_GOLANG:
		platform := &golang.Platform{}
		return platform.ValidateCodePackage(ccPkg)
	case peer.ChaincodeSpec_NODE:
		platform := &node.Platform{}
		return platform.ValidateCodePackage(ccPkg)
	case peer.ChaincodeSpec_JAVA:
		platform := &java.Platform{}
		return platform.ValidateCodePackage(ccPkg)
	case peer.ChaincodeSpec_CAR:
	}

	return fmt.Errorf("unsupport package platform -> %v", tYPE.String())
}
