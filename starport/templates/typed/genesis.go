package typed

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/genny"
)

func (t *typedStargate) genesisModify(opts *Options, g *genny.Generator) {
	g.RunFn(t.genesisProtoModify(opts))
	g.RunFn(t.genesisTypesModify(opts))
	g.RunFn(t.genesisModuleModify(opts))
}

func (t *typedStargate) genesisProtoModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/genesis.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateProtoImport := `%[1]v
import "%[2]v/%[3]v.proto";`
		replacementProtoImport := fmt.Sprintf(templateProtoImport, PlaceholderGenesisProtoImport, opts.ModuleName, opts.TypeName)
		content := strings.Replace(f.String(), PlaceholderGenesisProtoImport, replacementProtoImport, 1)

		// Determine the new field number
		fieldNumber := strings.Count(content, PlaceholderGenesisProtoStateField) + 1

		templateProtoState := `%[1]v
		repeated %[2]v %[3]vList = %[4]v; %[5]v`
		replacementProtoState := fmt.Sprintf(
			templateProtoState,
			PlaceholderGenesisProtoState,
			strings.Title(opts.TypeName),
			opts.TypeName,
			fieldNumber,
			PlaceholderGenesisProtoStateField,
		)
		content = strings.Replace(content, PlaceholderGenesisProtoState, replacementProtoState, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) genesisTypesModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/genesis.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateTypesImport := `"fmt"`
		content := strings.Replace(f.String(), PlaceholderGenesisTypesImport, templateTypesImport, 1)

		templateTypesDefault := `%[1]v
%[2]vList: []*%[2]v{},`
		replacementTypesDefault := fmt.Sprintf(templateTypesDefault, PlaceholderGenesisTypesDefault, strings.Title(opts.TypeName))
		content = strings.Replace(content, PlaceholderGenesisTypesDefault, replacementTypesDefault, 1)

		templateTypesValidate := `%[1]v
// Check for duplicated ID in %[2]v
%[2]vIdMap := make(map[uint64]bool)

for _, elem := range gs.%[3]vList {
	if _, ok := %[2]vIdMap[elem.Id]; ok {
		return fmt.Errorf("duplicated id for %[2]v")
	}
	%[2]vIdMap[elem.Id] = true
}`
		replacementTypesValidate := fmt.Sprintf(
			templateTypesValidate,
			PlaceholderGenesisTypesValidate,
			opts.TypeName,
			strings.Title(opts.TypeName),
		)
		content = strings.Replace(content, PlaceholderGenesisTypesValidate, replacementTypesValidate, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) genesisModuleModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/genesis.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateModuleInit := `%[1]v
// Set all the %[2]v
for _, elem := range genState.%[3]vList {
	k.Set%[3]v(ctx, *elem)
}

// Set %[2]v count
k.Set%[3]vCount(ctx, uint64(len(genState.%[3]vList)))
`
		replacementModuleInit := fmt.Sprintf(
			templateModuleInit,
			PlaceholderGenesisModuleInit,
			opts.TypeName,
			strings.Title(opts.TypeName),
		)
		content := strings.Replace(f.String(), PlaceholderGenesisModuleInit, replacementModuleInit, 1)

		templateModuleExport := `%[1]v
// Get all %[2]v
%[2]vList := k.GetAll%[3]v(ctx)
for _, elem := range %[2]vList {
	elem := elem
	genesis.%[3]vList = append(genesis.%[3]vList, &elem)
}
`
		replacementModuleExport := fmt.Sprintf(
			templateModuleExport,
			PlaceholderGenesisModuleExport,
			opts.TypeName,
			strings.Title(opts.TypeName),
		)
		content = strings.Replace(content, PlaceholderGenesisModuleExport, replacementModuleExport, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
