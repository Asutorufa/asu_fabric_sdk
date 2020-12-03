package chaincode

//var (
//	OrderingEndpoint           string = ""
//	tlsEnabled                 bool = false
//	clientAuth                 bool = false
//	caFile                     string = ""
//	keyFile                    string = ""
//	certFile                   string = ""
//	ordererTLSHostnameOverride string = ""
//	connTimeout                time.Duration = 3 * time.Second
//	tlsHandshakeTimeShift      time.Duration = 0
//)




// AddOrdererFlags adds flags for orderer-related commands
//func AddOrdererFlags(cmd *cobra.Command) {
//	flags := cmd.PersistentFlags()
//
//	flags.StringVarP(&OrderingEndpoint, "orderer", "o", "",
//		"Ordering service endpoint")
//	flags.BoolVarP(&tlsEnabled, "tls", "", false,
//		"Use TLS when communicating with the orderer endpoint")
//	flags.BoolVarP(&clientAuth, "clientauth", "", false,
//		"Use mutual TLS when communicating with the orderer endpoint")
//	flags.StringVarP(&caFile, "cafile", "", "",
//		"Path to file containing PEM-encoded trusted certificate(s) for the ordering endpoint")
//	flags.StringVarP(&keyFile, "keyfile", "", "",
//		"Path to file containing PEM-encoded private key to use for mutual TLS communication with the orderer endpoint")
//	flags.StringVarP(&certFile, "certfile", "", "",
//		"Path to file containing PEM-encoded X509 public key to use for mutual TLS communication with the orderer endpoint")
//	flags.StringVarP(&ordererTLSHostnameOverride, "ordererTLSHostnameOverride", "", "",
//		"The hostname override to use when validating the TLS connection to the orderer")
//	flags.DurationVarP(&connTimeout, "connTimeout", "", 3*time.Second,
//		"Timeout for client to connect")
//	flags.DurationVarP(&tlsHandshakeTimeShift, "tlsHandshakeTimeShift", "", 0,
//		"The amount of time to shift backwards for certificate expiration checks during TLS handshakes with the orderer endpoint")
//}
