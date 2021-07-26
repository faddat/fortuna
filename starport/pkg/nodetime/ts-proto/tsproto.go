// Package tsproto provides access to protoc-gen-ts_proto protoc plugin.
package tsproto

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/tendermint/starport/starport/pkg/nodetime"
)

const pluginName = "protoc-gen-ts_proto"

var (
	once       sync.Once
	binaryPath string
)

// BinaryPath returns the path to the binary of the ts-proto plugin so it can be passed to
// protoc via --plugin option.
//
// protoc is very picky about binary names of its plugins. for ts-proto, binary name
// will be protoc-gen-ts_proto.
// see why: https://github.com/stephenh/ts-proto/blob/7f76c05/README.markdown#quickstart.
func BinaryPath() (path string, err error) {
	once.Do(func() {
		if err = nodetime.PlaceBinary(); err != nil {
			return
		}

		tmpdir := os.TempDir()
		binaryPath = filepath.Join(tmpdir, pluginName)

		// comforting protoc by giving protoc-gen-ts_proto name to the plugin's binary.
		script := fmt.Sprintf(`#!/bin/bash
%s ts-proto "$@"
`, nodetime.BinaryPath)

		err = os.WriteFile(binaryPath, []byte(script), 0755)
	})

	return binaryPath, err
}
