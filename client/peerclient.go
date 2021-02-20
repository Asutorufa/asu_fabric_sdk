package client

import (
	"context"
	"fmt"

	"github.com/Asutorufa/fabricsdk/client/grpcclient"

	"github.com/hyperledger/fabric-protos-go/peer"
)

//var (
//	address  string        = ""
//	override string        = ""              // peer_tls_serverhostoverride
//	timeout  time.Duration = 3 * time.Second // peer_client_connTimeout
//
//	useTLS bool = false // peer_tls_enabled
//	caPEM  []byte
//
//	requireClientCert bool   = false // peer_tls_clientAuthRequired
//	kerPEM            []byte         // peer_tls_clientKey_file
//	certPEM           []byte         // peer_tls_client_file
//)

type PeerClient struct {
	GrpcClient *grpcclient.GRPCClient
	Client
}

func NewPeerClient(address, override string, Opt ...func(*grpcclient.ClientConfig)) (p *PeerClient, err error) {
	config := &grpcclient.ClientConfig{}

	for index := range Opt {
		Opt[index](config)
	}
	//fmt.Println(string(config.SecOpts.Certificate))
	//fmt.Println(config.SecOpts.ServerRootCAs)
	//fmt.Println(string(config.SecOpts.Key))
	p = new(PeerClient)
	p.address = address
	p.sn = override
	p.GrpcClient, err = grpcclient.NewGRPCClient(config)
	if err != nil {
		return nil, err
	}
	p.grpcConn, err = p.GrpcClient.NewConnection(p.address, grpcclient.ServerNameOverride(p.sn))
	return
}

func NewPeerClientSelf(address, override string, Opt ...func(config *grpcclient.ClientConfig)) (*PeerClient, error) {
	fmt.Printf("new peer [%s] client\n", address)
	c, err := NewClient(address, override, Opt...)
	if err != nil {
		return nil, err
	}
	return &PeerClient{
		Client: *c,
	}, nil
}

// Endorser returns a client for the Endorser service
func (pc *PeerClient) Endorser() (peer.EndorserClient, error) {
	return peer.NewEndorserClient(pc.grpcConn), nil
}

// Deliver returns a client for the Deliver service
func (pc *PeerClient) Deliver() (peer.Deliver_DeliverClient, error) {
	return peer.NewDeliverClient(pc.grpcConn).Deliver(context.TODO())
}

// PeerDeliver returns a client for the Deliver service for peer-specific use
// cases (i.e. DeliverFiltered)
func (pc *PeerClient) PeerDeliver() (peer.DeliverClient, error) {
	return peer.NewDeliverClient(pc.grpcConn), nil
}

// SnapshotClient returns a client for the snapshot service
func (pc *PeerClient) SnapshotClient() (peer.SnapshotClient, error) {
	return peer.NewSnapshotClient(pc.grpcConn), nil
}

// Certificate returns the TLS client certificate (if available)
// func (pc *PeerClient) Certificate() tls.Certificate {
// return pc.GrpcClient.Certificate()
// }

// func (pc *PeerClient) Close() (err error) {
// return pc.grpcConn.Close()
// }
