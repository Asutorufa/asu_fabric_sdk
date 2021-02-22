package main

import (
	"fmt"

	"github.com/Asutorufa/fabricsdk/chaincode"
)

func get(a string) {
	fmt.Println(chaincode.Query2(
		chaincode.ChainOpt{Path: "sacc", Name: "sacc", IsInit: true, Version: "1.0.4"},
		chaincode.MSPOpt{
			Path: "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/peers/peer-0-baas98/msp",
			ID:   "baas98",
		},
		[][]byte{[]byte("get"), []byte(a)},
		map[string][]byte{},
		"channel1",
		[]chaincode.EndpointWithPath{
			{
				Address: "192.168.9.196:30060",
			},
		},
	))
}
func main() {
	get("a")
	get("b")
}
