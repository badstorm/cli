package build

import (
	"embed"
	"testing"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/packd"
	"github.com/paganotoni/fsbox"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/template_validator
var validatefs embed.FS

var goodTemplates = func() packd.Box {
	box := packd.NewMemoryBox()
	box.AddString("_ignored/c.html", "c")
	box.AddString("a.html", "a")
	box.AddString("b.md", "b")
	return box
}()

func Test_TemplateValidator_Good(t *testing.T) {
	r := require.New(t)

	tvs := []TemplateValidator{PlushValidator}

	run := gentest.NewRunner()
	run.WithRun(ValidateTemplates(goodTemplates, tvs))

	r.NoError(run.Run())
}

func Test_TemplateValidator_Bad(t *testing.T) {
	r := require.New(t)

	box := fsbox.New(validatefs, "testdata/template_validator/bad")
	tvs := []TemplateValidator{PlushValidator}

	run := gentest.NewRunner()
	run.WithRun(ValidateTemplates(box, tvs))

	err := run.Run()
	r.Error(err)
	r.Equal("template error in file a.html: line 1: no prefix parse function for > found\ntemplate error in file b.md: line 1: no prefix parse function for > found", err.Error())
}
