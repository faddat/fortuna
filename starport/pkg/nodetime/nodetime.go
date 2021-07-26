// Package nodetime provides a single, and standalone NodeJS runtime executable that contains
// several NodeJS CLI programs bundled inside where those are reachable via subcommands.
// the CLI bundled programs are the ones that needed by Starport and more can added as needed.
package nodetime

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// BinaryPath is the path where nodetime binary is located in the fs.
const BinaryPath = "/tmp/nodetime"

// the list of CLIs included.
const (
	// CommandTSProto is https://github.com/stephenh/ts-proto.
	CommandTSProto = "ts-proto"

	// CommandTSC is https://github.com/microsoft/TypeScript.
	CommandTSC = "tsc"

	// CommandSTA is https://github.com/acacode/swagger-typescript-api.
	CommandSTA = "sta"
)

var (
	onceBinary      sync.Once
	oncePlaceBinary sync.Once
	binary          []byte
)

// Binary returns the binary bytes of the executable.
func Binary() []byte {
	onceBinary.Do(func() {
		// untar the binary.
		gzr, err := gzip.NewReader(bytes.NewReader(binaryCompressed))
		if err != nil {
			panic(err)
		}
		defer gzr.Close()

		tr := tar.NewReader(gzr)

		if _, err := tr.Next(); err != nil {
			panic(err)
		}

		if binary, err = io.ReadAll(tr); err != nil {
			panic(err)
		}
	})

	return binary
}

// PlaceBinary places the binary to BinaryPath.
func PlaceBinary() error {
	var err error

	oncePlaceBinary.Do(func() {
		// make sure that parent dir of the binary exists.
		if err = os.MkdirAll(filepath.Dir(BinaryPath), os.ModePerm); err != nil {
			return
		}

		// place the binary to BinaryPath.
		var f *os.File

		if f, err = os.OpenFile(BinaryPath, os.O_RDWR|os.O_CREATE, 0755); err != nil {
			return
		}
		defer f.Close()

		_, err = io.Copy(f, bytes.NewReader(Binary()))
	})

	return err
}
