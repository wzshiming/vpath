package vpath

import (
	"reflect"
)

var defaul = newVpath()

// Lookup to the value
func Lookup(val interface{}, path string) []interface{} {
	return defaul.Lookup(val, path)
}

// LookupBy to the value
func LookupBy(vals []reflect.Value, path []string) (ret []reflect.Value) {
	return defaul.LookupBy(vals, path)
}
