package actions

import (
	"github.com/gobuffalo/buffalo/render"
	"github.com/paganotoni/fsbox"

	"coke"
)

var (
	r *render.Engine

	assetsBox    = fsbox.New(coke.Assets, "public")
	templatesBox = fsbox.New(coke.Templates, "templates")
)

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.plush.html",

		// Box containing all of the templates:
		TemplatesBox: templatesBox,
		AssetsBox:    assetsBox,

		// Add template helpers here:
		Helpers: render.Helpers{
			// uncomment for non-Bootstrap form helpers:
			// "form":     plush.FormHelper,
			// "form_for": plush.FormForHelper,
		},
	})
}
