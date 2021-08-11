package build

import (
	"embed"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/paganotoni/fsbox"
	"github.com/stretchr/testify/require"
)

// TODO: once `buffalo new` is converted to use genny
// create an integration test that first generates a new application
// and then tries to build using genny/build.
var (
	//go:embed testdata/coke
	cokefs embed.FS

	coke = fsbox.New(cokefs, "testdata/coke")
)

var cokeRunner = func() *genny.Runner {
	run := gentest.NewRunner()
	run.Disk.AddBox(coke)

	content, _ := coke.FindString("dependencies.json")
	run.Disk.Add(genny.NewFileS("package.json", content))
	run.Root = "testdata/coke"

	return run
}

var eq = func(r *require.Assertions, s string, c *exec.Cmd) {
	if runtime.GOOS == "windows" {
		s = strings.Replace(s, "bin/build", `bin\build.exe`, 1)
		s = strings.Replace(s, "bin/foo", `bin\foo.exe`, 1)
	}
	r.Equal(s, strings.Join(c.Args, " "))
}

func Test_New(t *testing.T) {
	r := require.New(t)

	run := cokeRunner()

	opts := &Options{
		WithAssets:    true,
		WithBuildDeps: true,
		Environment:   "bar",
		App:           meta.New("."),
	}
	opts.App.Bin = "bin/foo"
	r.NoError(run.WithNew(New(opts)))
	run.Root = opts.App.Root

	r.NoError(run.Run())

	res := run.Results()

	// we should never leave any files modified or dropped
	r.Len(res.Files, 0)

	cmds := []string{"go get -d", "go build -tags bar -o bin/foo", "go mod tidy"}
	r.Len(res.Commands, len(cmds))
	for i, c := range res.Commands {
		eq(r, cmds[i], c)
	}
}

func Test_NewWithoutBuildDeps(t *testing.T) {
	envy.Temp(func() {
		envy.Set(envy.GO111MODULE, "off")
		r := require.New(t)

		run := cokeRunner()

		opts := &Options{
			WithAssets:    false,
			WithBuildDeps: false,
			Environment:   "bar",
			App:           meta.New("."),
		}
		opts.App.Bin = "bin/foo"
		r.NoError(run.WithNew(New(opts)))
		run.Root = opts.App.Root

		r.NoError(run.Run())

		res := run.Results()

		cmds := []string{"go get -d", "go build -tags bar -o bin/foo"}
		r.Len(res.Commands, len(cmds))
		for i, c := range res.Commands {
			eq(r, cmds[i], c)
		}
	})
}
