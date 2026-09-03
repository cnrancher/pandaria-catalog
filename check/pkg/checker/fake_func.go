package checker

import "text/template"

// Add fake template functions here
var fakeFuncsMap = template.FuncMap{
	"default":         fakeFunction,
	"include":         fakeFunction,
	"trunc":           fakeFunction,
	"trimSuffix":      fakeFunction,
	"contains":        fakeFunction,
	"replace":         fakeFunction,
	"quote":           fakeFunction,
	"toYaml":          fakeFunction,
	"fromYaml":        fakeFunction,
	"empty":           fakeFunction,
	"has":             fakeFunction,
	"hasKey":          fakeFunction,
	"keys":            fakeFunction,
	"deepCopy":        fakeFunction,
	"set":             fakeFunction,
	"sha256sum":       fakeFunction,
	"list":            fakeFunction,
	"semverCompare":   fakeFunction,
	"nindent":         fakeFunction,
	"trim":            fakeFunction,
	"split":           fakeFunction,
	"regexReplaceAll": fakeFunction,
	"toString":        fakeFunction,
	"dict":            fakeFunction,
	"int":             fakeFunction,
	"append":          fakeFunction,
	"lookup":          fakeFunction,
	"b64dec":          fakeFunction,
	"splitList":       fakeFunction,
	"first":           fakeFunction,
	"b64enc":          fakeFunction,
}

func fakeFunction() string {
	return ""
}
