package orderclient

import (
	"context"
	"crypto/tls"
	"fabricSDK/chaincode/client/grpcclient"

	"google.golang.org/grpc"

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
	grpcConn   *grpc.ClientConn
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
	o.grpcConn, err = o.GRPCClient.NewConnection(o.address, grpcclient.ServerNameOverride(o.sn))
	return
}

func (o *OrdererClient) Broadcast() (ordererProtos.AtomicBroadcast_BroadcastClient, error) {
	return ordererProtos.NewAtomicBroadcastClient(o.grpcConn).Broadcast(context.TODO())
}

func (o *OrdererClient) Deliver() (ordererProtos.AtomicBroadcast_DeliverClient, error) {
	return ordererProtos.NewAtomicBroadcastClient(o.grpcConn).Deliver(context.TODO())
}

func (o *OrdererClient) Certificate() tls.Certificate {
	return o.GRPCClient.Certificate()
}

func (o *OrdererClient) Close() error {
	return o.grpcConn.Close()
}
