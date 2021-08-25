package build

import (
	"embed"
	"time"

	"github.com/gobuffalo/cli/internal/runtime"
	"github.com/paganotoni/fsbox"

	"github.com/gobuffalo/events"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/plushgen"
	"github.com/gobuffalo/plush/v4"
)

var (
	//go:embed templates
	templates embed.FS
)

// New generator for building a Buffalo application
// This powers the `buffalo build` command and can be
// used to programatically build/compile Buffalo
// applications.
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}
	g.ErrorFn = func(err error) {
		events.EmitError(EvtBuildStopErr, err, events.Payload{"opts": opts})
	}

	g.RunFn(func(r *genny.Runner) error {
		events.EmitPayload(EvtBuildStart, events.Payload{"opts": opts})
		return nil
	})

	g.Transformer(genny.Dot())

	// validate templates
	g.RunFn(ValidateTemplates(templateWalker(opts.App), opts.TemplateValidators))

	// rename main() to originalMain()
	g.RunFn(transformMain(opts))

	// add any necessary templates for the build
	box := fsbox.New(templates, "templates", fsbox.OptionFSIgnoreGoEnv)
	if err := g.Box(box); err != nil {
		return g, err
	}

	// configure plush
	ctx := plush.NewContext()
	ctx.Set("opts", opts)
	ctx.Set("buildTime", opts.BuildTime.Format(time.RFC3339))
	ctx.Set("buildVersion", opts.BuildVersion)
	ctx.Set("buffaloVersion", runtime.Version)
	g.Transformer(plushgen.Transformer(ctx))

	// create the ./a pkg
	ag, err := apkg(opts)
	if err != nil {
		return g, err
	}
	g.Merge(ag)

	if opts.WithAssets {
		// mount the assets generator
		ag, err := assets(opts)
		if err != nil {
			return g, err
		}
		g.Merge(ag)
	}

	// create the final go build command
	c, err := buildCmd(opts)
	if err != nil {
		return g, err
	}

	g.Command(c)
	g.RunFn(func(r *genny.Runner) error {
		events.EmitPayload(EvtBuildStop, events.Payload{"opts": opts})
		return nil
	})

	g.RunFn(Cleanup(opts))
	return g, nil
}
