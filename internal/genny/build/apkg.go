package build

import (
	"github.com/gobuffalo/genny/v2"
)

func apkg(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	// g.RunFn(copyInflections)
	// g.RunFn(copyDatabase)
	// g.RunFn(addDependencies)

	return g, nil
}
