package chaincode

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/peer"
	policydsl "github.com/hyperledger/fabric/common/policydsl"
	"github.com/hyperledger/fabric/common/util"
)

// Chaincode-related variables.
//var (
//	chaincodeLang         string = "golang"
//	chaincodeCtorJSON     string = "{}" //
//	chaincodePath         string = ""
//	chaincodeName         string = ""
//	chaincodeUsr          string = "" // Not used
//	chaincodeQueryRaw     bool
//	chaincodeQueryHex     bool
//	channelID             string = "" //
//	chaincodeVersion      string = ""
//	policy                string = ""
//	escc                  string = ""
//	vscc                  string = ""
//	policyMarshalled      []byte
//	transient             string = ""
//	isInit                bool   = false //
//	collectionsConfigFile string = ""
//	collectionConfigBytes []byte
//	peerAddresses         []string      = []string{""}     //
//	tlsRootCertFiles      []string      = []string{""}     //
//	connectionProfile     string        = ""               //
//	waitForEvent          bool          = false            //
//	waitForEventTimeout   time.Duration = 30 * time.Second //
//)

//var (
//	install     = "install"
//	instantiate = "instantiate"
//	upgrade     = "upgrade"
//	pack        = "package"
//)

//func check(option string) (err error) {
//	if chaincodeName == "" {
//		return errors.New("chaincode name is empty")
//	}
//
//	if option == install || option == instantiate || option == upgrade || option == pack {
//		if chaincodeVersion == "" {
//			return errors.New("chaincode version is empty")
//		}
//	}
//
//	if escc == "" {
//		escc = "escc"
//	}
//
//	if vscc == "" {
//		vscc = "vscc"
//	}
//
//	if policy != "" {
//		p, err := policydsl.FromString(policy)
//		if err != nil {
//			return err
//		}
//		policyMarshalled = protoutil.MarshalOrPanic(p)
//	}
//
//	if connectionProfile != "" {
//		_, collectionConfigBytes, err = GetCollectionConfigFromFile(connectionProfile)
//		if err != nil {
//			return err
//		}
//	}
//
//	// Check that non-empty chaincode parameters contain only Args as a key.
//	// Type checking is done later when the JSON is actually unmarshaled
//	// into a pb.ChaincodeInput. To better understand what's going
//	// on here with JSON parsing see http://blog.golang.org/json-and-go -
//	// Generic JSON with interface{}
//	if chaincodeCtorJSON != "{}" {
//		var f interface{}
//		err := json.Unmarshal([]byte(chaincodeCtorJSON), &f)
//		if err != nil {
//			//return errors.Wrap(err, "chaincode argument error")
//			return err
//		}
//		m := f.(map[string]interface{})
//		sm := make(map[string]interface{})
//		for k := range m {
//			sm[strings.ToLower(k)] = m[k]
//		}
//		_, argsPresent := sm["args"]
//		_, funcPresent := sm["function"]
//		if !argsPresent || (len(m) == 2 && !funcPresent) || len(m) > 2 {
//			return errors.New("non-empty JSON chaincode parameters must contain the following keys: 'Args' or 'Function' and 'Args'")
//		}
//		//} else {
//		//	if cmd == nil || (cmd != chaincodeInstallCmd && cmd != chaincodePackageCmd) {
//		//		return errors.New("empty JSON chaincode parameters must contain the following keys: 'Args' or 'Function' and 'Args'")
//		//	}
//	}
//	return nil
//}

//func checkChaincodeCmdParams(cmd *cobra.Command) error {
//	// we need chaincode name for everything, including deploy
//	if chaincodeName == common.UndefinedParamValue {
//		return errors.Errorf("must supply value for %s name parameter", chainFuncName)
//	}
//
//	if cmd.Name() == instantiateCmdName || cmd.Name() == installCmdName ||
//		cmd.Name() == upgradeCmdName || cmd.Name() == packageCmdName {
//		if chaincodeVersion == common.UndefinedParamValue {
//			return errors.Errorf("chaincode version is not provided for %s", cmd.Name())
//		}
//
//		if escc != common.UndefinedParamValue {
//			logger.Infof("Using escc %s", escc)
//		} else {
//			logger.Info("Using default escc")
//			escc = "escc"
//		}
//
//		if vscc != common.UndefinedParamValue {
//			logger.Infof("Using vscc %s", vscc)
//		} else {
//			logger.Info("Using default vscc")
//			vscc = "vscc"
//		}
//
//		if policy != common.UndefinedParamValue {
//			p, err := policydsl.FromString(policy)
//			if err != nil {
//				return errors.Errorf("invalid policy %s", policy)
//			}
//			policyMarshalled = protoutil.MarshalOrPanic(p)
//		}
//
//		if collectionsConfigFile != common.UndefinedParamValue {
//			var err error
//			_, collectionConfigBytes, err = GetCollectionConfigFromFile(collectionsConfigFile)
//			if err != nil {
//				return errors.WithMessagef(err, "invalid collection configuration in file %s", collectionsConfigFile)
//			}
//		}
//	}
//
//	// Check that non-empty chaincode parameters contain only Args as a key.
//	// Type checking is done later when the JSON is actually unmarshaled
//	// into a pb.ChaincodeInput. To better understand what's going
//	// on here with JSON parsing see http://blog.golang.org/json-and-go -
//	// Generic JSON with interface{}
//	if chaincodeCtorJSON != "{}" {
//		var f interface{}
//		err := json.Unmarshal([]byte(chaincodeCtorJSON), &f)
//		if err != nil {
//			return errors.Wrap(err, "chaincode argument error")
//		}
//		m := f.(map[string]interface{})
//		sm := make(map[string]interface{})
//		for k := range m {
//			sm[strings.ToLower(k)] = m[k]
//		}
//		_, argsPresent := sm["args"]
//		_, funcPresent := sm["function"]
//		if !argsPresent || (len(m) == 2 && !funcPresent) || len(m) > 2 {
//			return errors.New("non-empty JSON chaincode parameters must contain the following keys: 'Args' or 'Function' and 'Args'")
//		}
//	} else {
//		if cmd == nil || (cmd != chaincodeInstallCmd && cmd != chaincodePackageCmd) {
//			return errors.New("empty JSON chaincode parameters must contain the following keys: 'Args' or 'Function' and 'Args'")
//		}
//	}
//
//	return nil
//}

// GetCollectionConfigFromFile retrieves the collection configuration
// from the supplied file; the supplied file must contain a
// json-formatted array of collectionConfigJson elements
func GetCollectionConfigFromFile(ccFile string) (*peer.CollectionConfigPackage, []byte, error) {
	fileBytes, err := ioutil.ReadFile(ccFile)
	if err != nil {
		//return nil, nil, errors.Wrapf(err, "could not read file '%s'", ccFile)
		return nil, nil, errors.New("could not read file")
	}

	return getCollectionConfigFromBytes(fileBytes)
}

type endorsementPolicy struct {
	ChannelConfigPolicy string `json:"channelConfigPolicy,omitempty"`
	SignaturePolicy     string `json:"signaturePolicy,omitempty"`
}

type collectionConfigJson struct {
	Name              string             `json:"name"`
	Policy            string             `json:"policy"`
	RequiredPeerCount *int32             `json:"requiredPeerCount"`
	MaxPeerCount      *int32             `json:"maxPeerCount"`
	BlockToLive       uint64             `json:"blockToLive"`
	MemberOnlyRead    bool               `json:"memberOnlyRead"`
	MemberOnlyWrite   bool               `json:"memberOnlyWrite"`
	EndorsementPolicy *endorsementPolicy `json:"endorsementPolicy,omitempty"`
}

// getCollectionConfig retrieves the collection configuration
// from the supplied byte array; the byte array must contain a
// json-formatted array of collectionConfigJson elements
func getCollectionConfigFromBytes(cconfBytes []byte) (*peer.CollectionConfigPackage, []byte, error) {
	cconf := &[]collectionConfigJson{}
	err := json.Unmarshal(cconfBytes, cconf)
	if err != nil {
		//return nil, nil, errors.Wrap(err, "could not parse the collection configuration")
		return nil, nil, errors.New("could not parse the collection configuration")
	}

	ccarray := make([]*peer.CollectionConfig, 0, len(*cconf))
	for _, cconfitem := range *cconf {
		p, err := policydsl.FromString(cconfitem.Policy)
		if err != nil {
			//return nil, nil, errors.WithMessagef(err, "invalid policy %s", cconfitem.Policy)
			return nil, nil, errors.New("invalid policy")
		}

		cpc := &peer.CollectionPolicyConfig{
			Payload: &peer.CollectionPolicyConfig_SignaturePolicy{
				SignaturePolicy: p,
			},
		}

		var ep *peer.ApplicationPolicy
		if cconfitem.EndorsementPolicy != nil {
			signaturePolicy := cconfitem.EndorsementPolicy.SignaturePolicy
			channelConfigPolicy := cconfitem.EndorsementPolicy.ChannelConfigPolicy
			ep, err = getApplicationPolicy(signaturePolicy, channelConfigPolicy)
			if err != nil {
				return nil, nil, errors.New("invalid endorsement policy [%#v]")
			}
		}

		// Set default requiredPeerCount and MaxPeerCount if not specified in json
		requiredPeerCount := int32(0)
		maxPeerCount := int32(1)
		if cconfitem.RequiredPeerCount != nil {
			requiredPeerCount = *cconfitem.RequiredPeerCount
		}
		if cconfitem.MaxPeerCount != nil {
			maxPeerCount = *cconfitem.MaxPeerCount
		}

		cc := &peer.CollectionConfig{
			Payload: &peer.CollectionConfig_StaticCollectionConfig{
				StaticCollectionConfig: &peer.StaticCollectionConfig{
					Name:              cconfitem.Name,
					MemberOrgsPolicy:  cpc,
					RequiredPeerCount: requiredPeerCount,
					MaximumPeerCount:  maxPeerCount,
					BlockToLive:       cconfitem.BlockToLive,
					MemberOnlyRead:    cconfitem.MemberOnlyRead,
					MemberOnlyWrite:   cconfitem.MemberOnlyWrite,
					EndorsementPolicy: ep,
				},
			},
		}

		ccarray = append(ccarray, cc)
	}

	ccp := &peer.CollectionConfigPackage{Config: ccarray}
	ccpBytes, err := proto.Marshal(ccp)
	return ccp, ccpBytes, err
}

func getApplicationPolicy(signaturePolicy, channelConfigPolicy string) (*peer.ApplicationPolicy, error) {
	if signaturePolicy == "" && channelConfigPolicy == "" {
		// no policy, no problem
		return nil, nil
	}

	if signaturePolicy != "" && channelConfigPolicy != "" {
		// mo policies, mo problems
		return nil, errors.New(`cannot specify both "--signature-policy" and "--channel-config-policy"`)
	}

	var applicationPolicy *peer.ApplicationPolicy
	if signaturePolicy != "" {
		signaturePolicyEnvelope, err := policydsl.FromString(signaturePolicy)
		if err != nil {
			//return nil, errors.Errorf("invalid signature policy: %s", signaturePolicy)
			return nil, errors.New("invalid signature")
		}

		applicationPolicy = &peer.ApplicationPolicy{
			Type: &peer.ApplicationPolicy_SignaturePolicy{
				SignaturePolicy: signaturePolicyEnvelope,
			},
		}
	}

	if channelConfigPolicy != "" {
		applicationPolicy = &peer.ApplicationPolicy{
			Type: &peer.ApplicationPolicy_ChannelConfigPolicyReference{
				ChannelConfigPolicyReference: channelConfigPolicy,
			},
		}
	}

	return applicationPolicy, nil
}

//var chaincodeCmd = &cobra.Command{
//	Use:   chainFuncName,
//	Short: fmt.Sprint(chainCmdDes),
//	Long:  fmt.Sprint(chainCmdDes),
//	PersistentPreRun: func(cmd *cobra.Command, args []string) {
//		common.InitCmd(cmd, args)
//		common.SetOrdererEnv(cmd, args)
//	},
//}

//var flags *pflag.FlagSet

//func init() {
//	resetFlags()
//}

// Explicitly define a method to facilitate tests
//func resetFlags() {
//flags = &pflag.FlagSet{}
//
//flags.StringVarP(&chaincodeLang, "lang", "l", "golang",
//	fmt.Sprintf("Language the %s is written in", chainFuncName))
//flags.StringVarP(&chaincodeCtorJSON, "ctor", "c", "{}",
//	fmt.Sprintf("Constructor message for the %s in JSON format", chainFuncName))
//flags.StringVarP(&chaincodePath, "path", "p", common.UndefinedParamValue,
//	fmt.Sprintf("Path to %s", chainFuncName))
//flags.StringVarP(&chaincodeName, "name", "n", common.UndefinedParamValue,
//	"Name of the chaincode")
//flags.StringVarP(&chaincodeVersion, "version", "v", common.UndefinedParamValue,
//	"Version of the chaincode specified in install/instantiate/upgrade commands")
//flags.StringVarP(&chaincodeUsr, "username", "u", common.UndefinedParamValue,
//	"Username for chaincode operations when security is enabled")
//flags.StringVarP(&channelID, "channelID", "C", "",
//	"The channel on which this command should be executed")
//flags.StringVarP(&policy, "policy", "P", common.UndefinedParamValue,
//	"The endorsement policy associated to this chaincode")
//flags.StringVarP(&escc, "escc", "E", common.UndefinedParamValue,
//	"The name of the endorsement system chaincode to be used for this chaincode")
//flags.StringVarP(&vscc, "vscc", "V", common.UndefinedParamValue,
//	"The name of the verification system chaincode to be used for this chaincode")
//flags.BoolVarP(&isInit, "isInit", "I", false,
//	"Is this invocation for init (useful for supporting legacy chaincodes in the new lifecycle)")
//flags.BoolVarP(&getInstalledChaincodes, "installed", "", false,
//	"Get the installed chaincodes on a peer")
//flags.BoolVarP(&getInstantiatedChaincodes, "instantiated", "", false,
//	"Get the instantiated chaincodes on a channel")
//flags.StringVar(&collectionsConfigFile, "collections-config", common.UndefinedParamValue,
//	"The fully qualified path to the collection JSON file including the file name")
//flags.StringArrayVarP(&peerAddresses, "peerAddresses", "", []string{common.UndefinedParamValue},
//	"The addresses of the peers to connect to")
//flags.StringArrayVarP(&tlsRootCertFiles, "tlsRootCertFiles", "", []string{common.UndefinedParamValue},
//	"If TLS is enabled, the paths to the TLS root cert files of the peers to connect to. The order and number of certs specified should match the --peerAddresses flag")
//flags.StringVarP(&connectionProfile, "connectionProfile", "", common.UndefinedParamValue,
//	"Connection profile that provides the necessary connection information for the network. Note: currently only supported for providing peer connection information")
//flags.BoolVar(&waitForEvent, "waitForEvent", false,
//	"Whether to wait for the event from each peer's deliver filtered service signifying that the 'invoke' transaction has been committed successfully")
//flags.DurationVar(&waitForEventTimeout, "waitForEventTimeout", 30*time.Second,
//	"Time to wait for the event from each peer's deliver filtered service signifying that the 'invoke' transaction has been committed successfully")
//flags.BoolVarP(&createSignedCCDepSpec, "cc-package", "s", false,
//	"create CC deployment spec for owner endorsements instead of raw CC deployment spec")
//flags.BoolVarP(&signCCDepSpec, "sign", "S", false,
//	"if creating CC deployment spec package for owner endorsements, also sign it with local MSP")
//flags.StringVarP(&instantiationPolicy, "instantiate-policy", "i", "",
//	"instantiation policy for the chaincode")
//

// chaincodeInput is wrapper around the proto defined ChaincodeInput message that
// is decorated with a custom JSON unmarshaller.
type chaincodeInput struct {
	peer.ChaincodeInput
}

// UnmarshalJSON converts the string-based REST/JSON input to
// the []byte-based current ChaincodeInput structure.
func (c *chaincodeInput) UnmarshalJSON(b []byte) error {
	sa := struct {
		Function string
		Args     []string
	}{}
	err := json.Unmarshal(b, &sa)
	if err != nil {
		return err
	}
	allArgs := sa.Args
	if sa.Function != "" {
		allArgs = append([]string{sa.Function}, sa.Args...)
	}
	c.Args = util.ToChaincodeArgs(allArgs...)
	return nil
}

// getChaincodeSpec get chaincode spec from the cli cmd parameters
//func getChaincodeSpec() (*peer.ChaincodeSpec, error) {
//	spec := &peer.ChaincodeSpec{}
//	if err := check(""); err != nil {
//		// unset usage silence because it's a command line usage error
//		//cmd.SilenceUsage = false
//		return spec, err
//	}
//
//	// Build the spec
//	input := chaincodeInput{}
//	if err := json.Unmarshal([]byte(chaincodeCtorJSON), &input); err != nil {
//		//return spec, errors.Wrap(err, "chaincode argument error")
//		return spec,err
//	}
//	input.IsInit = isInit
//
//	chaincodeLang = strings.ToUpper(chaincodeLang)
//	spec = &peer.ChaincodeSpec{
//		Type:        peer.ChaincodeSpec_Type(peer.ChaincodeSpec_Type_value[chaincodeLang]),
//		ChaincodeId: &peer.ChaincodeID{Path: chaincodePath, Name: chaincodeName, Version: chaincodeVersion},
//		Input:       &input.ChaincodeInput,
//	}
//	return spec, nil
//}
