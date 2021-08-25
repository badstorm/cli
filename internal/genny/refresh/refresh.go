package refresh

import (
	"embed"
	_ "embed"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/plushgen"
	"github.com/gobuffalo/plush/v4"
	"github.com/paganotoni/fsbox"
)

var (
	//go:embed templates
	templates embed.FS
)

// New generator to generate refresh templates
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()
	if err := opts.Validate(); err != nil {
		return g, err
	}
	g.Box(fsbox.New(templates, "templates", fsbox.OptionFSIgnoreGoEnv))

	ctx := plush.NewContext()
	ctx.Set("app", opts.App)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("dot-", "."))

	return g, nil
}
