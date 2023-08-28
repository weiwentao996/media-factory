package sources

import (
	"path"
	"runtime"
)

var Path string

func init() {
	Path = GetSourcesPath()
}

func GetSourcesPath() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
