package main

import (
	"fabricSDK/chaincode"
	"fmt"
	"time"
)

func get(a string) {
	fmt.Println(chaincode.Query2(
		chaincode.ChainOpt{Path: "sacc", Name: "sacc", IsInit: true, Version: "1.0.4"},
		chaincode.GrpcTLSOpt2{
			ClientCrtPath:      "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/client.crt",
			ClientKeyPath:      "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/client.key",
			CaPath:             "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/ca.crt",
			ServerNameOverride: "peer-0-baas98",
			Timeout:            6 * time.Second,
		},
		[][]byte{[]byte("get"), []byte(a)},
		"channel1",
		[]string{"192.168.9.196:30060"},
	))
}
func main() {
	get("a")
	get("b")
}
