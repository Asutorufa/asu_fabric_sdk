package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"

	"github.com/Asutorufa/fabricsdk/client/grpcclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//Client grpc client
type Client struct {
	address      string
	sn           string
	grpcConn     *grpc.ClientConn
	certificates []tls.Certificate
}

//NewClient new grpc client
func NewClient(address, override string, Opt ...func(config *grpcclient.ClientConfig)) (*Client, error) {
	config := &grpcclient.ClientConfig{}

	for oi := range Opt {
		Opt[oi](config)
	}

	client := &Client{
		address: address,
		sn:      override,
	}

	var opt []grpc.DialOption

	if config.SecOpts.UseTLS || config.SecOpts.RequireClientCert {
		c := &tls.Config{
			ServerName: override,
		}
		if config.SecOpts.UseTLS {
			certPool := x509.NewCertPool()
			for i := range config.SecOpts.ServerRootCAs {
				certPool.AppendCertsFromPEM(config.SecOpts.ServerRootCAs[i])
			}
			c.RootCAs = certPool
		}

		if config.SecOpts.RequireClientCert {
			cert, err := tls.X509KeyPair(config.SecOpts.Certificate, config.SecOpts.Key)
			if err != nil {
				return nil, err
			}

			c.Certificates = append(c.Certificates, cert)
		}
		client.certificates = c.Certificates
		opt = append(opt, grpc.WithTransportCredentials(credentials.NewTLS(c)))
	} else {
		opt = append(opt, grpc.WithInsecure())
	}

	// TODO KeepALive

	opt = append(opt, grpc.WithBlock()) // 阻塞
	opt = append(opt, grpc.FailOnNonTempDialError(true))

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	var err error
	client.grpcConn, err = grpc.DialContext(ctx, address, opt...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

//Certificate get all certificates
func (c *Client) Certificate() tls.Certificate {
	// return o.GRPCClient.Certificate()
	cert := tls.Certificate{}
	if len(c.certificates) > 0 {
		cert = c.certificates[0]
	}
	return cert
}

//Close close grpc connection
func (c *Client) Close() error {
	return c.grpcConn.Close()
}
