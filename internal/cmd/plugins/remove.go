package plugins

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/gobuffalo/cli/internal/plugins/plugdeps"
	"github.com/gobuffalo/cli/internal/takeon/github.com/markbates/errx"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/meta"
	"github.com/spf13/cobra"
)

var removeOptions = struct {
	dryRun bool
	vendor bool
}{}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "removes plugin from config/buffalo-plugins.toml",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("you must specify at least one package")
		}
		run := genny.WetRunner(context.Background())
		if removeOptions.dryRun {
			run = genny.DryRunner(context.Background())
		}

		app := meta.New(".")
		plugs, err := plugdeps.List(app)
		if err != nil && (errx.Unwrap(err) != plugdeps.ErrMissingConfig) {
			return err
		}

		for _, a := range args {
			a = strings.TrimSpace(a)
			bin := path.Base(a)
			plugs.Remove(plugdeps.Plugin{
				Binary: bin,
				GoGet:  a,
			})
		}

		run.WithRun(NewEncodePluginsRunner(app, plugs))
		return run.Run()
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&removeOptions.dryRun, "dry-run", "d", false, "dry run")
	removeCmd.Flags().BoolVar(&removeOptions.vendor, "vendor", false, "will install plugin binaries into ./plugins [WINDOWS not currently supported]")
}
