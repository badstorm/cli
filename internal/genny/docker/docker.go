package docker

import (
	"embed"
	"text/template"

	"github.com/gobuffalo/cli/internal/runtime"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/paganotoni/fsbox"
)

var (
	//go:embed templates
	templates embed.FS
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if len(opts.Version) == 0 {
		opts.Version = runtime.Version
	}

	data := map[string]interface{}{
		"opts": opts,
	}

	g.Box(fsbox.New(templates, "templates", fsbox.OptionFSIgnoreGoEnv))

	helpers := template.FuncMap{}
	t := gogen.TemplateTransformer(data, helpers)
	g.Transformer(t)
	g.Transformer(genny.Replace("dot-", "."))

	return g, nil
}
