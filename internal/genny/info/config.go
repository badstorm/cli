package info

import (
	"io/fs"
	"io/ioutil"
	"path"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/packd"
)

// ListWalker allows for a box that supports listing and walking
type ListWalker interface {
	packd.Lister
	packd.Walkable
}

func configs(opts *Options, root string) genny.RunFn {
	return func(r *genny.Runner) error {
		return filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}

			dat, err := ioutil.ReadFile(p)
			if err != nil {
				return err
			}

			opts.Out.Header("Buffalo: " + path.Join(p))
			opts.Out.WriteString(string(dat) + "\n")

			return nil
		})
	}
}
