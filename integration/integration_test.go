// Package integration_test integration test Starport and scaffolded apps.
package integration_test

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/tendermint/starport/starport/pkg/xexec"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if err := checkSystemRequirements(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func checkSystemRequirements() error {
	if !xexec.IsCommandAvailable("starport") {
		return errors.New("starport needs to be installed")
	}
	return nil
}
