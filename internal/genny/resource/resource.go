package resource

import (
	"embed"
	"text/template"

	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/gobuffalo/packd"
	"github.com/paganotoni/fsbox"
)

var (
	//go:embed templates
	templates embed.FS
)

// New resource generator
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	if !opts.SkipTemplates {
		core := fsbox.New(templates, "templates/core", fsbox.OptionFSIgnoreGoEnv)

		if err := g.Box(core); err != nil {
			return g, err
		}
	}

	var abox packd.Box
	if opts.SkipModel {
		abox = fsbox.New(templates, "templates/standard", fsbox.OptionFSIgnoreGoEnv)
	} else {
		abox = fsbox.New(templates, "templates/use_model", fsbox.OptionFSIgnoreGoEnv)
	}

	if err := g.Box(abox); err != nil {
		return g, err
	}

	pres := presenter{
		App:   opts.App,
		Name:  name.New(opts.Name),
		Model: name.New(opts.Model),
		Attrs: opts.Attrs,
	}

	x := pres.Name.Resource().File().String()
	folder := pres.Name.Folder().Pluralize().String()
	g.Transformer(genny.Replace("resource-name", x))
	g.Transformer(genny.Replace("resource-use_model", x))
	g.Transformer(genny.Replace("folder-name", folder))
	g.Transformer(genny.Replace("[und]_", "_"))

	data := map[string]interface{}{
		"opts":    pres,
		"actions": actions(opts),
		"folder":  folder,
	}

	helpers := template.FuncMap{
		"camelize": func(s string) string {
			return flect.Camelize(s)
		},
	}

	g.Transformer(gogen.TemplateTransformer(data, helpers))
	g.RunFn(installPop(opts))
	g.RunFn(addResource(pres))

	return g, nil
}
