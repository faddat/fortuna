package moduleimport

import (
	"embed"
	"fmt"
	"strings"

	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/templates/module"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
)

var (
	//go:embed stargate/* stargate/**/*
	fsStargate embed.FS

	stargateTemplate = xgenny.NewEmbedWalker(fsStargate, "stargate/")
)

// New ...
func NewImportStargate(opts *ImportOptions) (*genny.Generator, error) {
	g := genny.New()
	g.RunFn(appModifyStargate(opts))
	g.RunFn(rootModifyStargate(opts))
	if err := g.Box(stargateTemplate); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("AppName", opts.AppName)
	ctx.Set("title", strings.Title)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{binaryNamePrefix}}", opts.BinaryNamePrefix))
	return g, nil
}

// app.go modification on Stargate when importing wasm
func appModifyStargate(opts *ImportOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := module.PathAppGo
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateImport := `%[1]v
		"strings"
		"github.com/CosmWasm/wasmd/x/wasm"
		wasmclient "github.com/CosmWasm/wasmd/x/wasm/client"`
		replacementImport := fmt.Sprintf(templateImport, module.PlaceholderSgAppModuleImport)
		content := strings.Replace(f.String(), module.PlaceholderSgAppModuleImport, replacementImport, 1)

		templateEnabledProposals := `var (
			// If EnabledSpecificProposals is "", and this is "true", then enable all x/wasm proposals.
			// If EnabledSpecificProposals is "", and this is not "true", then disable all x/wasm proposals.
			ProposalsEnabled = "false"
			// If set to non-empty string it must be comma-separated list of values that are all a subset
			// of "EnableAllProposals" (takes precedence over ProposalsEnabled)
			// https://github.com/CosmWasm/wasmd/blob/02a54d33ff2c064f3539ae12d75d027d9c665f05/x/wasm/internal/types/proposal.go#L28-L34
			EnableSpecificProposals = ""
		)
		
		// GetEnabledProposals parses the ProposalsEnabled / EnableSpecificProposals values to
		// produce a list of enabled proposals to pass into wasmd app.
		func GetEnabledProposals() []wasm.ProposalType {
			if EnableSpecificProposals == "" {
				if ProposalsEnabled == "true" {
					return wasm.EnableAllProposals
				}
				return wasm.DisableAllProposals
			}
			chunks := strings.Split(EnableSpecificProposals, ",")
			proposals, err := wasm.ConvertToProposals(chunks)
			if err != nil {
				panic(err)
			}
			return proposals
		}`
		content = strings.Replace(content, module.PlaceholderSgWasmAppEnabledProposals, templateEnabledProposals, 1)

		templateGovProposalHandlers := `%[1]v
		govProposalHandlers = wasmclient.ProposalHandlers`
		replacementProposalHandlers := fmt.Sprintf(templateGovProposalHandlers, module.PlaceholderSgAppGovProposalHandlers)
		content = strings.Replace(content, module.PlaceholderSgAppGovProposalHandlers, replacementProposalHandlers, 1)

		templateModuleBasic := `%[1]v
		wasm.AppModuleBasic{},`
		replacementModuleBasic := fmt.Sprintf(templateModuleBasic, module.PlaceholderSgAppModuleBasic)
		content = strings.Replace(content, module.PlaceholderSgAppModuleBasic, replacementModuleBasic, 1)

		templateKeeperDeclaration := `%[1]v
		wasmKeeper wasm.Keeper`
		replacementKeeperDeclaration := fmt.Sprintf(templateKeeperDeclaration, module.PlaceholderSgAppKeeperDeclaration)
		content = strings.Replace(content, module.PlaceholderSgAppKeeperDeclaration, replacementKeeperDeclaration, 1)

		templateEnabledProposalsArgument := `%[1]v
		enabledProposals []wasm.ProposalType,`
		replacementEnabledProposalsArgument := fmt.Sprintf(templateEnabledProposalsArgument, module.PlaceholderSgAppNewArgument)
		content = strings.Replace(content, module.PlaceholderSgAppNewArgument, replacementEnabledProposalsArgument, 1)

		templateStoreKey := `%[1]v
		wasm.StoreKey,`
		replacementStoreKey := fmt.Sprintf(templateStoreKey, module.PlaceholderSgAppStoreKey)
		content = strings.Replace(content, module.PlaceholderSgAppStoreKey, replacementStoreKey, 1)

		templateKeeperDefinition := `%[1]v
		var wasmRouter = bApp.Router()
		wasmDir := filepath.Join(homePath, "wasm")
	
		wasmConfig, err := wasm.ReadWasmConfig(appOpts)
		if err != nil {
			panic("error while reading wasm config: " + err.Error())
		}

		// The last arguments can contain custom message handlers, and custom query handlers,
		// if we want to allow any custom callbacks
		supportedFeatures := "staking"
		app.wasmKeeper = wasm.NewKeeper(
			appCodec,
			keys[wasm.StoreKey],
			app.GetSubspace(wasm.ModuleName),
			app.AccountKeeper,
			app.BankKeeper,
			app.StakingKeeper,
			app.DistrKeeper,
			wasmRouter,
			wasmDir,
			wasmConfig,
			supportedFeatures,
			nil,
			nil,
		)
	
		// The gov proposal types can be individually enabled
		if len(enabledProposals) != 0 {
			govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.wasmKeeper, enabledProposals))
		}`
		replacementKeeperDefinition := fmt.Sprintf(templateKeeperDefinition, module.PlaceholderSgAppKeeperDefinition)
		content = strings.Replace(content, module.PlaceholderSgAppKeeperDefinition, replacementKeeperDefinition, 1)

		templateAppModule := `%[1]v
		wasm.NewAppModule(&app.wasmKeeper, app.StakingKeeper),`
		replacementAppModule := fmt.Sprintf(templateAppModule, module.PlaceholderSgAppAppModule)
		content = strings.Replace(content, module.PlaceholderSgAppAppModule, replacementAppModule, 1)

		templateInitGenesis := `%[1]v
		wasm.ModuleName,`
		replacementInitGenesis := fmt.Sprintf(templateInitGenesis, module.PlaceholderSgAppInitGenesis)
		content = strings.Replace(content, module.PlaceholderSgAppInitGenesis, replacementInitGenesis, 1)

		templateParamSubspace := `%[1]v
		paramsKeeper.Subspace(wasm.ModuleName)`
		replacementParamSubspace := fmt.Sprintf(templateParamSubspace, module.PlaceholderSgAppParamSubspace)
		content = strings.Replace(content, module.PlaceholderSgAppParamSubspace, replacementParamSubspace, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// app.go modification on Stargate when importing wasm
func rootModifyStargate(opts *ImportOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := "cmd/" + opts.BinaryNamePrefix + "d/cmd/root.go"
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateImport := `%[1]v
		"github.com/CosmWasm/wasmd/x/wasm"`
		replacementImport := fmt.Sprintf(templateImport, module.PlaceholderSgRootImport)
		content := strings.Replace(f.String(), module.PlaceholderSgRootImport, replacementImport, 1)

		templateCommand := `%[1]v
		AddGenesisWasmMsgCmd(app.DefaultNodeHome),`
		replacementCommand := fmt.Sprintf(templateCommand, module.PlaceholderSgRootCommands)
		content = strings.Replace(content, module.PlaceholderSgRootCommands, replacementCommand, 1)

		templateInitFlags := `%[1]v
		wasm.AddModuleInitFlags(startCmd)`
		replacementInitFlags := fmt.Sprintf(templateInitFlags, module.PlaceholderSgRootInitFlags)
		content = strings.Replace(content, module.PlaceholderSgRootInitFlags, replacementInitFlags, 1)

		templateenabledProposals := `%[1]v
		app.GetEnabledProposals(),`
		replacementAppArgument := fmt.Sprintf(templateenabledProposals, module.PlaceholderSgRootAppArgument)
		content = strings.Replace(content, module.PlaceholderSgRootAppArgument, replacementAppArgument, 1)

		replacementExportArgument := fmt.Sprintf(templateenabledProposals, module.PlaceholderSgRootExportArgument)
		content = strings.Replace(content, module.PlaceholderSgRootExportArgument, replacementExportArgument, 1)

		replacementNoHeightExportArgument := fmt.Sprintf(templateenabledProposals, module.PlaceholderSgRootNoHeightExportArgument)
		content = strings.Replace(content, module.PlaceholderSgRootNoHeightExportArgument, replacementNoHeightExportArgument, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
