package info

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/gobuffalo/clara/v2/genny/rx"
	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_pkgChecks(t *testing.T) {
	r := require.New(t)
	d := t.TempDir()
	os.Chdir(d)

	bb := &bytes.Buffer{}
	run := gentest.NewRunner()
	opts := &Options{
		App: meta.New("."),
		Out: rx.NewWriter(bb),
	}

	err := ioutil.WriteFile("go.mod", []byte("module foo"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	run.WithRun(pkgChecks(opts, d))

	r.NoError(run.Run())

	res := bb.String()
	r.Contains(res, "Buffalo: go.mod")
}
