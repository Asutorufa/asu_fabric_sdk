package client

import (
	"context"

	"github.com/Asutorufa/fabricsdk/client/grpcclient"
	ordererProtos "github.com/hyperledger/fabric-protos-go/orderer"
)

//OrdererClient orderer client use grpc
type OrdererClient struct {
	GRPCClient *grpcclient.GRPCClient
	Client
}

//NewOrdererClient create new orderer client
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

//NewOrdererClientSelf create new orderer client
func NewOrdererClientSelf(address, override string, Opt ...func(config *grpcclient.ClientConfig)) (*OrdererClient, error) {
	c, err := NewClient(address, override, Opt...)
	if err != nil {
		return nil, err
	}

	return &OrdererClient{
		Client: *c,
	}, nil
}

//Broadcast orderer broadcast client
func (o *OrdererClient) Broadcast() (ordererProtos.AtomicBroadcast_BroadcastClient, error) {
	return ordererProtos.NewAtomicBroadcastClient(o.grpcConn).Broadcast(context.TODO())
}

//Deliver orderer deliver client
func (o *OrdererClient) Deliver() (ordererProtos.AtomicBroadcast_DeliverClient, error) {
	return ordererProtos.NewAtomicBroadcastClient(o.grpcConn).Deliver(context.TODO())
}
