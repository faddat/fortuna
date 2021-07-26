package conf

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/imdario/mergo"
)

var (
	// ErrCouldntLocateConfig returned when config.yml cannot be found in the source code.
	ErrCouldntLocateConfig = errors.New("could not locate a config.yml in your chain. please follow the link for how-to: https://github.com/tendermint/starport/blob/develop/docs/1%20Introduction/4%20Configuration.md")

	// FileNames holds a list of appropriate names for the config file.
	FileNames = []string{"config.yml", "config.yaml"}

	// DefaultConf holds default configuration.
	DefaultConf = Config{
		Host: Host{
			RPC:      ":26657",
			P2P:      ":26656",
			Prof:     ":6060",
			GRPC:     ":9090",
			API:      ":1317",
			Frontend: ":8080",
			DevUI:    ":12345",
		},
		Build: Build{
			Proto: Proto{
				Path: "proto",
				ThirdPartyPaths: []string{
					"third_party/proto",
					"proto_vendor",
				},
			},
		},
		Faucet: Faucet{
			Host: ":4500",
		},
	}
)

// Config is the user given configuration to do additional setup
// during serve.
type Config struct {
	Accounts  []Account              `yaml:"accounts"`
	Validator Validator              `yaml:"validator"`
	Faucet    Faucet                 `yaml:"faucet"`
	Client    Client                 `yaml:"client"`
	Build     Build                  `yaml:"build"`
	Init      Init                   `yaml:"init"`
	Genesis   map[string]interface{} `yaml:"genesis"`
	Host      Host                   `yaml:"host"`
}

// AccountByName finds account by name.
func (c Config) AccountByName(name string) (acc Account, found bool) {
	for _, acc := range c.Accounts {
		if acc.Name == name {
			return acc, true
		}
	}
	return Account{}, false
}

// Account holds the options related to setting up Cosmos wallets.
type Account struct {
	Name     string   `yaml:"name"`
	Coins    []string `yaml:"coins,omitempty"`
	Mnemonic string   `yaml:"mnemonic,omitempty"`
	Address  string   `yaml:"address,omitempty"`

	// The RPCAddress off the chain that account is issued at.
	RPCAddress string `yaml:"rpc_address,omitempty"`
}

// Validator holds info related to validator settings.
type Validator struct {
	Name   string `yaml:"name"`
	Staked string `yaml:"staked"`
}

// Build holds build configs.
type Build struct {
	Binary string `yaml:"binary"`
	Proto  Proto  `yaml:"proto"`
}

// Proto holds proto build configs.
type Proto struct {
	// Path is the relative path of where app's proto files are located at.
	Path string `yaml:"path"`

	// ThirdPartyPath is the relative path of where the third party proto files are
	// located that used by the app.
	ThirdPartyPaths []string `yaml:"third_party_paths"`
}

// Client configures code generation for clients.
type Client struct {
	// Vuex configures code generation for Vuex.
	Vuex Vuex `yaml:"vuex"`
}

// Vuex configures code generation for Vuex.
type Vuex struct {
	// Path configures out location for generated Vuex code.
	Path string `yaml:"path"`
}

// Faucet configuration.
type Faucet struct {
	// Name is faucet account's name.
	Name *string `yaml:"name"`

	// Coins holds type of coin denoms and amounts to distribute.
	Coins []string `yaml:"coins"`

	// CoinsMax holds of chain denoms and their max amounts that can be transferred
	// to single user.
	CoinsMax []string `yaml:"coins_max"`

	// Host is the host of the faucet server
	Host string `yaml:"host"`

	// Port number for faucet server to listen at.
	Port int `yaml:"port"`
}

// Init overwrites sdk configurations with given values.
type Init struct {
	// App overwrites appd's config/app.toml configs.
	App map[string]interface{} `yaml:"app"`

	// Config overwrites appd's config/config.toml configs.
	Config map[string]interface{} `yaml:"config"`

	// Home overwrites default home directory used for the app
	Home string `yaml:"home"`

	// CLIHome overwrites default CLI home directory used for launchpad app
	CLIHome string `yaml:"cli-home"`

	// KeyringBackend is the default keyring backend to use for blockchain initialization
	KeyringBackend string `yaml:"keyring-backend"`
}

// Host keeps configuration related to started servers.
type Host struct {
	RPC      string `yaml:"rpc"`
	P2P      string `yaml:"p2p"`
	Prof     string `yaml:"prof"`
	GRPC     string `yaml:"grpc"`
	API      string `yaml:"api"`
	Frontend string `yaml:"frontend"`
	DevUI    string `yaml:"dev-ui"`
}

// Parse parses config.yml into UserConfig.
func Parse(r io.Reader) (Config, error) {
	var conf Config
	if err := yaml.NewDecoder(r).Decode(&conf); err != nil {
		return conf, err
	}
	if err := mergo.Merge(&conf, DefaultConf); err != nil {
		return Config{}, err
	}
	return conf, validate(conf)
}

// ParseFile parses config.yml from the path.
func ParseFile(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, nil
	}
	defer file.Close()
	return Parse(file)
}

// validate validates user config.
func validate(conf Config) error {
	if len(conf.Accounts) == 0 {
		return &ValidationError{"at least 1 account is needed"}
	}
	if conf.Validator.Name == "" {
		return &ValidationError{"validator is required"}
	}
	return nil
}

// ValidationError is returned when a configuration is invalid.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("config is not valid: %s", e.Message)
}

// LocateDefault locates the default path for the config file, if no file found returns ErrCouldntLocateConfig.
func LocateDefault(root string) (path string, err error) {
	for _, name := range FileNames {
		path = filepath.Join(root, name)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		} else if !os.IsNotExist(err) {
			return "", err
		}
	}
	return "", ErrCouldntLocateConfig
}

// FaucetHost returns the faucet host to use
func FaucetHost(conf Config) string {
	// We keep supporting Port option for backward compatibility
	// TODO: drop this option in the future
	host := conf.Faucet.Host
	if conf.Faucet.Port != 0 {
		host = fmt.Sprintf(":%d", conf.Faucet.Port)
	}

	return host
}
