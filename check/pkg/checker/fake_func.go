package checker

import "text/template"

// Add fake template functions here
var fakeFuncsMap = template.FuncMap{
	"default":    fakeFunction,
	"include":    fakeFunction,
	"trunc":      fakeFunction,
	"trimSuffix": fakeFunction,
	"contains":   fakeFunction,
	"replace":    fakeFunction,
	"quote":      fakeFunction,
	"toYaml":     fakeFunction,
	"fromYaml":   fakeFunction,
	"empty":      fakeFunction,
	"has":        fakeFunction,
	"hasKey":     fakeFunction,
	"keys":       fakeFunction,
	"deepCopy":   fakeFunction,
	"set":        fakeFunction,
	"sha256sum":  fakeFunction,
	"list":       fakeFunction,
}

func fakeFunction() string {
	return ""
}
