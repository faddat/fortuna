package chain

import (
	"context"

	starportconf "github.com/tendermint/starport/starport/chainconf"
	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
)

// TODO omit -cli log messages for Stargate.

type Plugin interface {
	// Name of a Cosmos version.
	Name() string

	// Setup performs the initial setup for plugin.
	Setup(context.Context) error

	// ConfigCommands returns step.Exec configuration for config commands.
	Configure(context.Context, chaincmdrunner.Runner, string) error

	// GentxCommand returns step.Exec configuration for gentx command.
	Gentx(context.Context, chaincmdrunner.Runner, Validator) (path string, err error)

	// PostInit hook.
	PostInit(string, starportconf.Config) error

	// StartCommands returns step.Exec configuration to start servers.
	Start(context.Context, chaincmdrunner.Runner, starportconf.Config) error

	// Home returns the blockchain node's home dir.
	Home() string

	// CLIHome returns the cli blockchain node's home dir.
	CLIHome() string

	// Version of the plugin.
	Version() cosmosver.MajorVersion

	// SupportsIBC reports if app support IBC.
	SupportsIBC() bool
}

func (c *Chain) pickPlugin() Plugin {
	switch c.Version.Major() {
	case cosmosver.Launchpad:
		return newLaunchpadPlugin(c.app)
	case cosmosver.Stargate:
		return newStargatePlugin(c.app)
	}
	panic("unknown cosmos version")
}
