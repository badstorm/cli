package info

import (
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
)

// New returns a generator that performs buffalo
// related rx checks
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	g.RunFn(appDetails(opts))
	g.RunFn(configs(opts, filepath.Join(opts.App.Root, "config")))
	// g.RunFn(pkgChecks(opts, opts.App.Root))

	return g, nil
}
