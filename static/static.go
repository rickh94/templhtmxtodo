package static

import (
	"embed"
	"templtodo3/config"

	"github.com/benbjohnson/hashfs"
)

//go:embed css/* js/* img/*
var Static embed.FS

var HashStatic = hashfs.NewFS(Static)

func StaticUrl(name string) string {
	return config.AppConfig.StaticHostname + "/static/" + HashStatic.HashName(name)
}
