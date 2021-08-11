package api

import (
	"embed"
	_ "embed"
	"html/template"

	"github.com/gobuffalo/cli/internal/genny/newapp/core"
	"github.com/paganotoni/fsbox"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

var (
	//go:embed templates
	templates embed.FS
)

// New generator for creating a Buffalo API application
func New(opts *Options) (*genny.Group, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	gg, err := core.New(opts.Options)
	if err != nil {
		return gg, err
	}

	g := genny.New()
	data := map[string]interface{}{
		"opts": opts,
	}

	helpers := template.FuncMap{}

	t := gogen.TemplateTransformer(data, helpers)
	g.Transformer(t)
	g.Box(fsbox.New(templates, "templates", fsbox.OptionFSIgnoreGoEnv))

	gg.Add(g)

	return gg, nil
}
