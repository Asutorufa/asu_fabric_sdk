package client

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/Asutorufa/fabricsdk/client/grpcclient"
)

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

func WithTLS2(caPEMPath string) func(client *grpcclient.ClientConfig) {
	data, err := ioutil.ReadFile(caPEMPath)
	if err != nil {
		log.Printf("caPEM read error, set to false -> %v\n", err)
		return func(client *grpcclient.ClientConfig) {}
	}
	return WithTLS(data)
}

func WithTLS(caPEM []byte) func(client *grpcclient.ClientConfig) {
	return func(client *grpcclient.ClientConfig) {
		client.SecOpts.UseTLS = true
		client.SecOpts.ServerRootCAs = [][]byte{caPEM}
	}
}

func WithClientCert2(keyPEMPath, certPEMPath string) func(client *grpcclient.ClientConfig) {
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
	return WithClientCert(key, cert)
}

func WithClientCert(keyPEM, certPEM []byte) func(client *grpcclient.ClientConfig) {
	return func(client *grpcclient.ClientConfig) {
		client.SecOpts.RequireClientCert = true
		client.SecOpts.Key = keyPEM
		client.SecOpts.Certificate = certPEM
	}
}
