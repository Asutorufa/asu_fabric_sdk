package chaincode

import "github.com/hyperledger/fabric/bccsp/factory"

//GetFactory .
func GetFactory(s string) {
	switch s {
	case "sw":
		(&factory.SWFactory{}).Get(factory.GetDefaultOpts())
	case "pkcs11":
		(&factory.SWFactory{}).Get(factory.GetDefaultOpts())
	}
}
