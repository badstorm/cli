package coke

import "embed"

var (
	//go:embed templates
	Templates embed.FS

	//go:embed public
	Assets embed.FS

	//go:embed config
	Config embed.FS
)
