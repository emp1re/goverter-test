package cli

import (
	"flag"
	"fmt"

	"github.com/emp1re/goverter-test/config"

	"github.com/emp1re/goverter-test"
)

type Strings []string

func (s Strings) String() string {
	return fmt.Sprint([]string(s))
}

func (s *Strings) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func Parse(args []string) (*goverter.GenerateConfig, error) {
	switch len(args) {
	case 0:
		return nil, usage("unknown")
	case 1:
		return nil, usageErr("missing command: ", args[0])
	default:
		if args[1] != "gen" {
			return nil, usageErr("unknown command: "+args[1], args[0])
		}
	}

	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fs.Usage = func() {}

	var global Strings
	fs.Var(&global, "global", "")
	fs.Var(&global, "g", "")

	buildTags := fs.String("build-tags", "goverter", "")
	outputConstraint := fs.String("output-constraint", "!goverter", "")

	if err := fs.Parse(args[2:]); err != nil {
		return nil, usageErr(err.Error(), args[0])
	}
	patterns := fs.Args()

	if len(patterns) == 0 {
		return nil, usageErr("missing PATTERN", args[0])
	}

	c := goverter.GenerateConfig{
		PackagePatterns:       patterns,
		BuildTags:             *buildTags,
		OutputBuildConstraint: *outputConstraint,
		Global: config.RawLines{
			Lines:    global,
			Location: "command line (-g, -global)",
		},
	}

	return &c, nil
}

func usageErr(err, cmd string) error {
	return fmt.Errorf("Error: %s\n%s", err, usage(cmd))
}

func usage(cmd string) error {
	return fmt.Errorf(`Usage:
  %s gen [OPTIONS] PACKAGE...

PACKAGE(s):
  Define the import paths goverter will use to search for converter interfaces.
  You can define multiple packages and use the special ... golang pattern to
  select multiple packages. See $ go help packages

OPTIONS:
  -build-tags [tags]: (default: goverter)
      a comma-separated list of additional build tags to consider satisfied
      during the loading of conversion interfaces. See 'go help buildconstraint'.
      Can be disabled by supplying an empty string.

  -g [value], -global [value]:
      apply settings to all defined converters. For a list of available
      settings see: https://goverter.jmattheis.de/reference/settings

  -h, --help:
      display this help page

  -output-constraint [constraint]: (default: !goverter)
      A build constraint added to all files generated by goverter.
      Can be disabled by supplying an empty string.

Examples:
  %s gen ./example/simple ./example/complex
  %s gen ./example/...
  %s gen github.com/emp1re/goverter-test/example/simple
  %s gen -g 'ignoreMissing no' -g 'skipCopySameType' ./simple

Documentation:
  Full documentation is available here: https://goverter.jmattheis.de`, cmd, cmd, cmd, cmd, cmd)
}
