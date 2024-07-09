package netgen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/celestiaorg/celestia-app/v2/app"
	"github.com/celestiaorg/celestia-app/v2/app/encoding"
	"github.com/celestiaorg/celestia-app/v2/pkg/appconsts"
	"github.com/celestiaorg/celestia-app/v2/test/util/genesis"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/tendermint/tendermint/config"
	cmtjson "github.com/tendermint/tendermint/libs/json"
	cmtos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/p2p/pex"
	"github.com/tendermint/tendermint/privval"
)

// NodeInfo is a struct that contains the name, IP address, and network address
// of a node.
type NodeInfo struct {
	Name           string `json:"name"`
	IP             string `json:"ip"`
	NetworkAddress string `json:"network_address"`
}

func (n NodeInfo) PeerID() string {
	return fmt.Sprintf("%s@%s:26656", n.NetworkAddress, n.IP)
}

// Network maintains the initial state of the network. This includes the
// genesis, all relevant validators included in the genesis, and all accounts.
type Network struct {
	genesis *genesis.Genesis
	ecfg    encoding.Config

	validators map[string]NodeInfo
	accounts   []string
}

func NewNetwork(chainID string) (*Network, error) {
	codec := encoding.MakeConfig(app.ModuleEncodingRegisters...)
	g := genesis.NewDefaultGenesis().WithChainID(chainID).WithModifiers(genesis.ImmediateProposals(codec.Codec))

	return &Network{
		genesis:    g,
		validators: make(map[string]NodeInfo),
		ecfg:       codec,
	}, nil
}

// AddValidator adds a validator to the network. The validator is identified by
// its name which is assigned by pulumi as hardware is allocated. An addional
// account and keyring are saved to the payload directory that can be used by
// txsim.
func (n *Network) AddValidator(name, ip, payLoadRoot string) error {
	n.validators[name] = NodeInfo{
		Name: name,
		IP:   ip,
	}

	err := n.genesis.NewValidator(genesis.NewDefaultValidator(name))
	if err != nil {
		return err
	}

	// add a txsim key and keyring to each validator
	kr, err := keyring.New(app.Name, keyring.BackendTest,
		filepath.Join(payLoadRoot, name), nil, n.ecfg.Codec)
	if err != nil {
		return err
	}

	key, _, err := kr.NewMnemonic("txsim", keyring.English, "", "", hd.Secp256k1)
	if err != nil {
		return err
	}

	pk, err := key.GetPubKey()
	if err != nil {
		return err
	}

	err = n.genesis.AddAccount(genesis.Account{
		PubKey:  pk,
		Balance: 9999999999999999,
	})

	if err != nil {
		return err
	}

	return nil

}

func (n *Network) Peers() []string {
	var peers []string
	for _, v := range n.validators {
		peers = append(peers, v.PeerID())
	}
	return peers

}

func (n *Network) AddAccount(name string) {
	n.accounts = append(n.accounts, name)
}

func (n *Network) InitNodes(rootDir string) error {
	if len(n.accounts) != 0 {
		n.genesis.WithKeyringAccounts(genesis.NewKeyringAccounts(genesis.DefaultInitialBalance, n.accounts...)...)
	}

	// save the genesis file
	genesisPath := filepath.Join(rootDir, "genesis.json")

	genDoc, err := n.genesis.Export()
	if err != nil {
		return err
	}

	genBytes, err := cmtjson.MarshalIndent(genDoc, "", "  ")
	if err != nil {
		return err
	}

	// save the genesis file as configured
	err = cmtos.WriteFile(genesisPath, genBytes, 0o644)
	if err != nil {
		return err
	}

	vals := n.genesis.Validators()
	for i, v := range vals {
		vname := fmt.Sprintf("validator-%d", i)
		valPath := filepath.Join(rootDir, vname)
		nodeKeyFile := filepath.Join(valPath, "node_key.json")
		if err := cmtos.EnsureDir(filepath.Dir(nodeKeyFile), 0o777); err != nil {
			return err
		}

		// add the network key assigned by the genesis to that validator's payload
		nodeKey := &p2p.NodeKey{
			PrivKey: v.NetworkKey,
		}
		if err := nodeKey.SaveAs(nodeKeyFile); err != nil {
			return err
		}
		ninfo, has := n.validators[vname]
		if !has {
			return fmt.Errorf("No validator found %s", vname)
		}
		ninfo.NetworkAddress = string(nodeKey.ID())
		n.validators[vname] = ninfo

		// generate remaining private key file using the assigned consensus key
		pvStateFile := filepath.Join(valPath, "priv_validator_state.json")
		if err := cmtos.EnsureDir(filepath.Dir(pvStateFile), 0o777); err != nil {
			return err
		}
		pvKeyFile := filepath.Join(valPath, "priv_validator_key.json")
		if err := cmtos.EnsureDir(filepath.Dir(pvKeyFile), 0o777); err != nil {
			return err
		}
		filePV := privval.NewFilePV(v.ConsensusKey, pvKeyFile, pvStateFile)
		filePV.Save()

		cmtcfg, err := MakeConfig(vname)
		if err != nil {
			return err
		}
		config.WriteConfigFile(filepath.Join(rootDir, vname, "config.toml"), cmtcfg)

		appcfg := MakeAppConfig()
		serverconfig.WriteConfigFile(filepath.Join(rootDir, vname, "app.toml"), appcfg)
	}

	return nil
}

// SaveValidatorsToFile saves the validators map as a JSON to the given file.
func (n *Network) SaveValidatorsToFile(filename string) error {
	// Open the file for writing. Create it if it doesn't exist.
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode the validators map to JSON and write it to the file.
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Optional: format the JSON with indentation
	err = encoder.Encode(n.validators)
	if err != nil {
		return err
	}

	return nil
}

func (n *Network) SaveAddressBook(payloadRoot string, peers []string) error {
	addrBookFile := filepath.Join(payloadRoot, "addrbook.json")
	return WriteAddressBook(peers, addrBookFile)
}

func MakeConfig(name string, opts ...Option) (*config.Config, error) {
	cfg := config.DefaultConfig()
	cfg.Moniker = name
	cfg.RPC.ListenAddress = "tcp://0.0.0.0:26657"
	// cfg.P2P.ExternalAddress = fmt.Sprintf("tcp://%v", node.AddressP2P(false))
	// cfg.P2P.PersistentPeers = strings.Join(node.InitialPeers, ",")
	cfg.Instrumentation.Prometheus = false
	cfg.Mempool.Size = 5000
	cfg.Mempool.CacheSize = 10000
	cfg.Mempool.MaxTxBytes = 100_000_000
	cfg.Mempool.MaxTxsBytes = 1_000_000_000
	cfg.Mempool.Version = "v1"
	cfg.Mempool.TTLNumBlocks = 100
	cfg.Mempool.TTLDuration = 40 * time.Minute
	cfg.Mempool.MaxGossipDelay = 20 * time.Second
	cfg.TxIndex.Indexer = "kv"
	cfg.P2P.MaxNumInboundPeers = 15
	cfg.P2P.MaxNumOutboundPeers = 10
	cfg.P2P.PexReactor = true
	cfg.P2P.RecvRate = 5_120_000
	cfg.P2P.SendRate = 5_120_000
	cfg.RPC.MaxBodyBytes = 1_000_000_000
	cfg.RPC.MaxOpenConnections = 1000
	cfg.RPC.TimeoutBroadcastTxCommit = 60 * time.Second
	cfg.RPC.MaxSubscriptionClients = 1000
	cfg.Consensus.TimeoutPropose = time.Second * 10
	cfg.Consensus.TimeoutCommit = time.Second * 11
	cfg.Consensus.OnlyInternalWal = true
	cfg.Instrumentation.TraceBufferSize = 5000
	cfg.Instrumentation.TraceType = "local"

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg, nil
}

func MakeAppConfig() *serverconfig.Config {
	cfg := serverconfig.DefaultConfig()
	cfg.API.Enable = true
	cfg.GRPC.Enable = true
	cfg.GRPCWeb.Enable = false
	cfg.GRPC.MaxRecvMsgSize = 1_000_000_000
	cfg.GRPC.MaxSendMsgSize = 1_000_000_000

	// the default snapshot interval was determined by picking a large enough
	// value as to not dramatically increase resource requirements while also
	// being greater than zero so that there are more nodes that will serve
	// snapshots to nodes that state sync
	cfg.StateSync.SnapshotInterval = 0
	cfg.StateSync.SnapshotKeepRecent = 1
	cfg.MinGasPrices = fmt.Sprintf("%v%s", appconsts.DefaultMinGasPrice, app.BondDenom)
	return cfg
}

type Option func(*config.Config)

func WriteAddressBook(peers []string, file string) error {
	book := pex.NewAddrBook(file, false)
	for _, peer := range peers {
		addr, err := p2p.NewNetAddressString(peer)
		if err != nil {
			return fmt.Errorf("parsing peer address %s: %w", peer, err)
		}
		err = book.AddAddress(addr, addr)
		if err != nil {
			return fmt.Errorf("adding peer address %s: %w", peer, err)
		}
	}
	book.Save()
	return nil
}

// func addPeersToAddressBook(path string, peers []PeerPacket) error {
// 	err := os.MkdirAll(strings.Replace(path, "addrbook.json", "", -1), os.ModePerm)
// 	if err != nil {
// 		return err
// 	}

// 	addrBook := pex.NewAddrBook(path, false)
// 	err = addrBook.OnStart()
// 	if err != nil {
// 		return err
// 	}

// 	for _, peer := range peers {
// 		id, ip, peerPort, err := parsePeerID(peer.PeerID)
// 		if err != nil {
// 			return err
// 		}
// 		port, err := safeConvertIntToUint16(peerPort)
// 		if err != nil {
// 			return err
// 		}

// 		netAddr := p2p.NetAddress{
// 			ID:   p2p.ID(id),
// 			IP:   ip,
// 			Port: port,
// 		}

// 		err = addrBook.AddAddress(&netAddr, &netAddr)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	addrBook.Save()
// 	return nil
// }