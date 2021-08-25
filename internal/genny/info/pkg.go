package info

import (
	"io/ioutil"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
)

func pkgChecks(opts *Options, folder string) genny.RunFn {
	return func(r *genny.Runner) error {
		dat, err := ioutil.ReadFile(filepath.Join(folder, "go.mod"))
		if err != nil {
			opts.Out.WriteString("go.mod not found, skipping package details.")

			return nil
		}

		opts.Out.Header("\nBuffalo: go.mod")
		opts.Out.Write(dat)

		return nil
	}
}
