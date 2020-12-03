package chaincode

// InitCmdFactory init the ChaincodeCmdFactory with default clients
//func InitCmdFactory(cmdName string, isEndorserRequired, isOrdererRequired bool, cryptoProvider bccsp.BCCSP) (*ChaincodeCmdFactory, error) {
//	var err error
//	var endorserClients []peer.EndorserClient
//	var deliverClients []peer.DeliverClient
//	if isEndorserRequired {
//		if err = validatePeerConnectionParameters(cmdName); err != nil {
//			//return nil, errors.WithMessage(err, "error validating peer connection parameters")
//
//			return nil,errors.New("error validating peer connection parameters")
//		}
//		for i, address := range peerAddresses {
//			var tlsRootCertFile string
//			if tlsRootCertFiles != nil {
//				tlsRootCertFile = tlsRootCertFiles[i]
//			}
//			endorserClient, err := common.GetEndorserClientFnc(address, tlsRootCertFile)
//			if err != nil {
//				return nil, errors.WithMessagef(err, "error getting endorser client for %s", cmdName)
//			}
//			endorserClients = append(endorserClients, endorserClient)
//			deliverClient, err := common.GetPeerDeliverClientFnc(address, tlsRootCertFile)
//			if err != nil {
//				return nil, errors.WithMessagef(err, "error getting deliver client for %s", cmdName)
//			}
//			deliverClients = append(deliverClients, deliverClient)
//		}
//		if len(endorserClients) == 0 {
//			return nil, errors.New("no endorser clients retrieved - this might indicate a bug")
//		}
//	}
//	certificate, err := common.GetCertificateFnc()
//	if err != nil {
//		return nil, errors.WithMessage(err, "error getting client certificate")
//	}
//
//	signer, err := common.GetDefaultSignerFnc()
//	if err != nil {
//		return nil, errors.WithMessage(err, "error getting default signer")
//	}
//
//	var broadcastClient common.BroadcastClient
//	if isOrdererRequired {
//		if len(common.OrderingEndpoint) == 0 {
//			if len(endorserClients) == 0 {
//				return nil, errors.New("orderer is required, but no ordering endpoint or endorser client supplied")
//			}
//			endorserClient := endorserClients[0]
//
//			orderingEndpoints, err := common.GetOrdererEndpointOfChainFnc(channelID, signer, endorserClient, cryptoProvider)
//			if err != nil {
//				return nil, errors.WithMessagef(err, "error getting channel (%s) orderer endpoint", channelID)
//			}
//			if len(orderingEndpoints) == 0 {
//				return nil, errors.Errorf("no orderer endpoints retrieved for channel %s, pass orderer endpoint with -o flag instead", channelID)
//			}
//			logger.Infof("Retrieved channel (%s) orderer endpoint: %s", channelID, orderingEndpoints[0])
//			// override viper env
//			viper.Set("orderer.address", orderingEndpoints[0])
//		}
//
//		broadcastClient, err = common.GetBroadcastClientFnc()
//		if err != nil {
//			return nil, errors.WithMessage(err, "error getting broadcast client")
//		}
//	}
//	return &ChaincodeCmdFactory{
//		EndorserClients: endorserClients,
//		DeliverClients:  deliverClients,
//		Signer:          signer,
//		BroadcastClient: broadcastClient,
//		Certificate:     certificate,
//	}, nil
//}

//func validatePeerConnectionParameters(cmdName string) error {
//	if connectionProfile != common.UndefinedParamValue {
//		networkConfig, err := common.GetConfig(connectionProfile)
//		if err != nil {
//			return err
//		}
//		if len(networkConfig.Channels[channelID].Peers) != 0 {
//			peerAddresses = []string{}
//			tlsRootCertFiles = []string{}
//			for peer, peerChannelConfig := range networkConfig.Channels[channelID].Peers {
//				if peerChannelConfig.EndorsingPeer {
//					peerConfig, ok := networkConfig.Peers[peer]
//					if !ok {
//						return errors.Errorf("peer '%s' is defined in the channel config but doesn't have associated peer config", peer)
//					}
//					peerAddresses = append(peerAddresses, peerConfig.URL)
//					tlsRootCertFiles = append(tlsRootCertFiles, peerConfig.TLSCACerts.Path)
//				}
//			}
//		}
//	}
//
//	// currently only support multiple peer addresses for invoke
//	multiplePeersAllowed := map[string]bool{
//		"invoke": true,
//	}
//	_, ok := multiplePeersAllowed[cmdName]
//	if !ok && len(peerAddresses) > 1 {
//		return errors.Errorf("'%s' command can only be executed against one peer. received %d", cmdName, len(peerAddresses))
//	}
//
//	if len(tlsRootCertFiles) > len(peerAddresses) {
//		logger.Warningf("received more TLS root cert files (%d) than peer addresses (%d)", len(tlsRootCertFiles), len(peerAddresses))
//	}
//
//	if viper.GetBool("peer.tls.enabled") {
//		if len(tlsRootCertFiles) != len(peerAddresses) {
//			return errors.Errorf("number of peer addresses (%d) does not match the number of TLS root cert files (%d)", len(peerAddresses), len(tlsRootCertFiles))
//		}
//	} else {
//		tlsRootCertFiles = nil
//	}
//
//	return nil
//}

// init viper
//func initViper(){
//	viper.SetEnvPrefix("core") // -> origin: common.CmdRoot
//	viper.AutomaticEnv()
//	replacer := strings.NewReplacer(".","_") // replace _ to .
//	viper.SetEnvKeyReplacer(replacer)

	//cryptoProvider := factory.GetDefault()
//}

//func main() {
//	// For environment variables.
//	viper.SetEnvPrefix(common.CmdRoot)
//	viper.AutomaticEnv()
//	replacer := strings.NewReplacer(".", "_")
//	viper.SetEnvKeyReplacer(replacer)
//
//	// Define command-line flags that are valid for all peer commands and
//	// subcommands.
//	mainFlags := mainCmd.PersistentFlags() // 初始化mainFlags的一些东西
//
//	mainFlags.String("logging-level", "", "Legacy logging level flag")
//	viper.BindPFlag("logging_level", mainFlags.Lookup("logging-level"))
//	mainFlags.MarkHidden("logging-level")
//
//	cryptoProvider := factory.GetDefault()
//
//	mainCmd.AddCommand(version.Cmd())
//	mainCmd.AddCommand(node.Cmd())
//	mainCmd.AddCommand(chaincode.Cmd(nil, cryptoProvider))
//	mainCmd.AddCommand(channel.Cmd(nil))
//	mainCmd.AddCommand(lifecycle.Cmd(cryptoProvider))
//	mainCmd.AddCommand(snapshot.Cmd(cryptoProvider))
//
//	// On failure Cobra prints the usage message and error string, so we only
//	// need to exit with a non-0 status
//	if mainCmd.Execute() != nil {
//		os.Exit(1)
//	}
//}
