package peerclient

import (
	"context"
	"crypto/tls"
	"fabricSDK/chaincode/grpcclient"
	"fmt"
	"io/ioutil"
	"log"
	"time"

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
	address    string
	sn         string
}

func WithTimeout(t time.Duration) func(config *grpcclient.ClientConfig) {
	return func(c *grpcclient.ClientConfig) {
		c.Timeout = t
	}
}

func WithTLS(caPEMPath string) func(client *grpcclient.ClientConfig) {
	data, err := ioutil.ReadFile(caPEMPath)
	if err != nil {
		log.Printf("caPEM read error, set to false -> %v\n", err)
		return func(client *grpcclient.ClientConfig) {}
	}
	return WithTLS2(data)
}

func WithTLS2(caPEM []byte) func(*grpcclient.ClientConfig) {
	return func(c *grpcclient.ClientConfig) {
		c.SecOpts.UseTLS = true
		c.SecOpts.ServerRootCAs = [][]byte{caPEM}
	}
}

func WithClientCert(keyPEMPath, certPEMPath string) func(client *grpcclient.ClientConfig) {
	key, err := ioutil.ReadFile(keyPEMPath)
	if err != nil {
		log.Printf("client key read error, set to false -> %v\n", err)
		return func(client *grpcclient.ClientConfig) {}
	}
	cert, err := ioutil.ReadFile(certPEMPath)
	if err != nil {
		log.Printf("client key read error, set to false -> %v\n", err)
		return func(client *grpcclient.ClientConfig) {}
	}
	return WithClientCert2(key, cert)
}

func WithClientCert2(keyPEM, certPEM []byte) func(*grpcclient.ClientConfig) {
	return func(c *grpcclient.ClientConfig) {
		c.SecOpts.UseTLS = true
		c.SecOpts.RequireClientCert = true
		c.SecOpts.Key = keyPEM
		c.SecOpts.Certificate = certPEM
	}
}

func NewPeerClient(address, override string, Opt ...func(*grpcclient.ClientConfig)) (p *PeerClient, err error) {
	config := &grpcclient.ClientConfig{}

	for index := range Opt {
		Opt[index](config)
	}
	fmt.Println(string(config.SecOpts.Certificate))
	fmt.Println(config.SecOpts.ServerRootCAs)
	fmt.Println(string(config.SecOpts.Key))
	p = new(PeerClient)
	p.address = address
	p.sn = override
	p.GrpcClient, err = grpcclient.NewGRPCClient(config)
	if err != nil {
		return nil, err
	}
	return
}

// Endorser returns a client for the Endorser service
func (pc *PeerClient) Endorser() (peer.EndorserClient, error) {
	conn, err := pc.GrpcClient.NewConnection(pc.address, grpcclient.ServerNameOverride(pc.sn))
	if err != nil {
		return nil, fmt.Errorf("%v: endorser client failed to connect to %s", err, pc.address)
	}
	return peer.NewEndorserClient(conn), nil
}

// Deliver returns a client for the Deliver service
func (pc *PeerClient) Deliver() (peer.Deliver_DeliverClient, error) {
	conn, err := pc.GrpcClient.NewConnection(pc.address, grpcclient.ServerNameOverride(pc.sn))
	if err != nil {
		//return nil, errors.WithMessagef(err, "deliver client failed to connect to %s", pc.Address)
		return nil, err
	}
	return peer.NewDeliverClient(conn).Deliver(context.TODO())
}

// PeerDeliver returns a client for the Deliver service for peer-specific use
// cases (i.e. DeliverFiltered)
func (pc *PeerClient) PeerDeliver() (peer.DeliverClient, error) {
	conn, err := pc.GrpcClient.NewConnection(pc.address, grpcclient.ServerNameOverride(pc.sn))
	if err != nil {
		//return nil, errors.WithMessagef(err, "deliver client failed to connect to %s", pc.Address)
		return nil, err
	}
	return peer.NewDeliverClient(conn), nil
}

// Certificate returns the TLS client certificate (if available)
func (pc *PeerClient) Certificate() tls.Certificate {
	return pc.GrpcClient.Certificate()
}

// SnapshotClient returns a client for the snapshot service
func (pc *PeerClient) SnapshotClient() (peer.SnapshotClient, error) {
	conn, err := pc.GrpcClient.NewConnection(pc.address, grpcclient.ServerNameOverride(pc.sn))
	if err != nil {
		//return nil, errors.WithMessagef(err, "snapshot client failed to connect to %s", pc.Address)
		return nil, err
	}
	return peer.NewSnapshotClient(conn), nil
}
