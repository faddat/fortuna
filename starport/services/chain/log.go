package chain

import (
	"io"
	"os"
	"strings"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/lineprefixer"
	"github.com/tendermint/starport/starport/pkg/prefixgen"
)

// prefixes holds prefix configuration for logs messages.
var prefixes = map[logType]struct {
	Name  string
	Color uint8
}{
	logStarport: {"starport", 202},
	logBuild:    {"build", 203},
	logAppd:     {"%s daemon", 204},
	logAppcli:   {"%s cli", 205},
}

// logType represents the different types of logs.
type logType int

const (
	logStarport logType = iota
	logBuild
	logAppd
	logAppcli
)

// std returns the cmdrunner steps to configure stdout and stderr to output logs by logType.
func (c *Chain) stdSteps(logType logType) []step.Option {
	std := c.stdLog(logType)
	return []step.Option{
		step.Stdout(std.out),
		step.Stderr(std.err),
	}
}

type std struct {
	out, err io.Writer
}

// std returns the stdout and stderr to output logs by logType.
func (c *Chain) stdLog(logType logType) std {
	prefixed := func(w io.Writer) *lineprefixer.Writer {
		var (
			prefix    = prefixes[logType]
			prefixStr string
			options   = prefixgen.Common(prefixgen.Color(prefix.Color))
			gen       = prefixgen.New(prefix.Name, options...)
		)
		if strings.Count(prefix.Name, "%s") > 0 {
			prefixStr = gen.Gen(c.app.Name)
		} else {
			prefixStr = gen.Gen()
		}
		return lineprefixer.NewWriter(w, func() string { return prefixStr })
	}
	var (
		stdout io.Writer = prefixed(c.stdout)
		stderr io.Writer = prefixed(c.stderr)
	)
	if logType == logStarport && c.logLevel == LogRegular {
		stdout = os.Stdout
		stderr = os.Stderr
	}
	return std{
		out: stdout,
		err: stderr,
	}
}

func (c *Chain) genPrefix(logType logType) string {
	prefix := prefixes[logType]

	return prefixgen.
		New(prefix.Name, prefixgen.Common(prefixgen.Color(prefix.Color))...).
		Gen(c.app.Name)
}
