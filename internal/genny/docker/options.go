package docker

import (
	"github.com/gobuffalo/meta"
)

type Options struct {
	App     meta.App `json:"app"`
	Version string   `json:"version"`
}
