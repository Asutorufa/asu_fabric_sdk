package orderclient

import (
	"context"
	"crypto/tls"
	"fabricSDK/chaincode/grpcclient"
	"io/ioutil"
	"log"
	"time"

	ordererProtos "github.com/hyperledger/fabric-protos-go/orderer"
)

//
//var (
//	address  string // orderer_address
//	override string // orderer_tls_serverhostoverride
//
//	connTimeout time.Duration // orderer_client_connTimeout
//
//	useTLS            bool          // orderer_tls_enabled
//	requirtClientCert bool          // orderer_tls_clientAuthRequired
//	timeShift         time.Duration // orderer_tls_handshakeTimeShift
//
//	caPEM []byte // orderer_tls_rootcert_file
//
//	keyPEM []byte // orderer_tls_clientKey_file
//
//	certPEM []byte // orderer_tls_clientCert_file
//
//	grpcClient *grpcclient.GRPCClient
//)
//
//func init() {
//	x := grpcclient.ClientConfig{}
//
//	x.SecOpts.UseTLS = useTLS
//	x.SecOpts.RequireClientCert = requirtClientCert
//	x.SecOpts.TimeShift = timeShift
//
//	if x.SecOpts.UseTLS {
//		x.SecOpts.ServerRootCAs = [][]byte{caPEM}
//	}
//
//	x.SecOpts.Key = keyPEM
//
//	if x.SecOpts.RequireClientCert {
//		x.SecOpts.Certificate = certPEM
//	}
//
//	var err error
//	grpcClient, err = grpcclient.NewGRPCClient(x)
//	if err != nil {
//		panic(err)
//	}
//}

type OrdererClient struct {
	GRPCClient *grpcclient.GRPCClient
	address    string
	sn         string
}

func WithTimeout(duration time.Duration) func(client *grpcclient.ClientConfig) {
	if duration == 0 {
		return func(client *grpcclient.ClientConfig) {
			client.Timeout = 3 * time.Second
		}
	}
	return func(client *grpcclient.ClientConfig) {
		client.Timeout = duration
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

func WithTLS2(caPEM []byte) func(client *grpcclient.ClientConfig) {
	return func(client *grpcclient.ClientConfig) {
		client.SecOpts.UseTLS = true
		client.SecOpts.ServerRootCAs = [][]byte{caPEM}
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

func WithClientCert2(keyPEM, certPEM []byte) func(client *grpcclient.ClientConfig) {
	return func(client *grpcclient.ClientConfig) {
		client.SecOpts.RequireClientCert = true
		client.SecOpts.Key = keyPEM
		client.SecOpts.Certificate = certPEM
	}
}

func NewOrdererClient(address, override string, Opt ...func(config *grpcclient.ClientConfig)) (o *OrdererClient, err error) {
	config := &grpcclient.ClientConfig{}

	for index := range Opt {
		Opt[index](config)
	}

	o = new(OrdererClient)
	o.address = address
	o.sn = override
	o.GRPCClient, err = grpcclient.NewGRPCClient(config)
	if err != nil {
		return nil, err
	}

	return
}

func (o *OrdererClient) Broadcast() (ordererProtos.AtomicBroadcast_BroadcastClient, error) {
	conn, err := o.GRPCClient.NewConnection(o.address, grpcclient.ServerNameOverride(o.sn))
	if err != nil {
		return nil, err
	}
	return ordererProtos.NewAtomicBroadcastClient(conn).Broadcast(context.TODO())
}

func (o *OrdererClient) Deliver() (ordererProtos.AtomicBroadcast_DeliverClient, error) {
	conn, err := o.GRPCClient.NewConnection(o.address, grpcclient.ServerNameOverride(o.sn))
	if err != nil {
		return nil, err
	}
	return ordererProtos.NewAtomicBroadcastClient(conn).Deliver(context.TODO())
}

func (o *OrdererClient) Certificate() tls.Certificate {
	return o.GRPCClient.Certificate()
}
