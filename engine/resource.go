package engine

import (
	"regexp"
	"strings"
)

type Resource interface {
	NamedObject
	GetPath() string
}

type ReadableResource interface {
	ReadData() ([]byte, error)
	ReadDataAsync(callback ReadResourceCallback)
}

type ReadResourceCallback func(data []byte, err error)

type ResourceDefinition struct {
	path   string
	write  bool
}

var validResourcePathRegexp, _ = regexp.Compile("^(/[\\w.]+)+$")
func IsValidResourcePath(path string) bool {
	return validResourcePathRegexp.MatchString(path)
}

func GetResourceName(path string) string {
	var tokens = strings.Split(path, "/")
	return tokens[len(tokens)-1]
}