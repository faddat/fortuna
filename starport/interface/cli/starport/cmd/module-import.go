package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

// NewModuleImport creates a new command to import an sdk module.
func NewModuleImport() *cobra.Command {
	c := &cobra.Command{
		Use:       "import [feature]",
		Short:     "Imports a new module to app.",
		Long:      "Use starport add wasm to add support for webassembly smart contracts to your blockchain.",
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: []string{"wasm"},
		RunE:      importModuleHandler,
	}
	return c
}

func importModuleHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	name := args[0]
	sc := scaffolder.New(appPath)
	if err := sc.ImportModule(name); err != nil {
		return err
	}

	s.Stop()

	fmt.Printf("\n🎉 Imported module `%s`.\n\n", name)
	return nil
}
