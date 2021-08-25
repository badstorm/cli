package info

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/clara/v2/genny/rx"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_configs(t *testing.T) {
	d := t.TempDir()
	os.Chdir(d)

	r := require.New(t)
	run := gentest.NewRunner()
	bb := &bytes.Buffer{}

	app := meta.New(".")
	opts := &Options{
		App: app,
		Out: rx.NewWriter(bb),
	}

	err := os.MkdirAll("config", 0777)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join("config", "buffalo-app.toml"), []byte("app"), 0777)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join("config", "buffalo-plugins.toml"), []byte("plugins"), 0777)
	if err != nil {
		t.Fatal(err)
	}

	run.WithRun(configs(opts, "config"))

	r.NoError(run.Run())

	x := bb.String()
	r.Contains(x, "Buffalo: config/buffalo-app.toml\napp")
	r.Contains(x, "Buffalo: config/buffalo-plugins.toml\nplugins")
}
