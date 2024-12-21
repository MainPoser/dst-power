package config

import "path"

var (
	WorkDir   string
	StaticDir string
	UiDir     string
)

func MediaDir() string {
	return path.Join(WorkDir, "media")
}
