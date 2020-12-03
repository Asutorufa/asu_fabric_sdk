package chaincode

import (
	"fmt"
	"testing"
	"time"
)

func get(t *testing.T, a string) {
	resp, err := Query2(
		ChainOpt{Path: "sacc", Name: "sacc", IsInit: true, Version: "1.0.4"},
		GrpcTLSOpt2{
			ClientCrtPath:      "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/client.crt",
			ClientKeyPath:      "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/client.key",
			CaPath:             "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/ca.crt",
			ServerNameOverride: "peer-0-baas98",
			Timeout:            6 * time.Second,
		},
		[][]byte{[]byte("get"), []byte(a)},
		"channel1",
		[]string{"192.168.9.196:30060"},
	)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)
	fmt.Println(resp.Response.Status, string(resp.Response.Payload))
}

func TestQuery(t *testing.T) {
	get(t, "a")
	get(t, "b")
}

func set(t *testing.T, a, b string) {
	resp, err := Invoke2(
		ChainOpt{Path: "sacc", Name: "sacc", IsInit: true, Version: "1.0.4"},
		GrpcTLSOpt2{
			ClientCrtPath:      "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/client.crt",
			ClientKeyPath:      "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/client.key",
			CaPath:             "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/ca.crt",
			ServerNameOverride: "peer-0-baas98",
			Timeout:            6 * time.Second,
		},
		GrpcTLSOpt2{
			ClientCrtPath:      "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/client.crt",
			ClientKeyPath:      "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/client.key",
			CaPath:             "/mnt/shareSSD/code/YunPhant/wasabi_3/src/wasabi/backEnd/conf/nfs_data/baas98/msp/baas98/users/Admin@baas98/tls/ca.crt",
			ServerNameOverride: "orderer-0-baas98",
			Timeout:            6 * time.Second,
		},
		[][]byte{[]byte("set"), []byte(a), []byte(b)},
		"channel1",
		[]string{"192.168.9.196:30060"},
		"192.168.9.196:30062",
	)

	if err != nil {
		t.Error(err)
	}

	t.Log(resp)
}

func TestInvoke(t *testing.T) {
	set(t, "a", "xiaoxiao")
	set(t, "b", "xiaoxiao2")
}
