package core

import (
	"embed"
	"html/template"
	"os/exec"

	"github.com/BurntSushi/toml"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/paganotoni/fsbox"
)

var (
	//go:embed templates
	templates embed.FS
)

func rootGenerator(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	// validate opts
	if err := opts.Validate(); err != nil {
		return g, err
	}

	g.Command(exec.Command("go", "mod", "init", opts.App.PackagePkg))

	g.Transformer(genny.Replace("dot-", "."))

	// add common templates
	err := g.Box(fsbox.New(templates, "templates", fsbox.OptionFSIgnoreGoEnv))
	if err != nil {
		return g, err
	}

	data := map[string]interface{}{
		"opts": opts,
	}

	helpers := template.FuncMap{}

	t := gogen.TemplateTransformer(data, helpers)
	g.Transformer(t)

	g.RunFn(func(r *genny.Runner) error {
		f := genny.NewFile("config/buffalo-app.toml", nil)
		if err := toml.NewEncoder(f).Encode(opts.App); err != nil {
			return err
		}
		return r.File(f)
	})

	return g, nil
}
