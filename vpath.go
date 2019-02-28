package vpath

import (
	"reflect"
	"strconv"
	"strings"
)

// Vpath value path
type vpath struct {
	fieldMap map[reflect.Type]map[string]int
}

// NewVpath Create a new Vpath
func newVpath() *vpath {
	return &vpath{
		fieldMap: map[reflect.Type]map[string]int{},
	}
}

func (v *vpath) field(typ reflect.Type) map[string]int {
	m, ok := v.fieldMap[typ]
	if ok {
		return m
	}
	num := typ.NumField()

	m = map[string]int{}
	for i := 0; i != num; i++ {
		field := typ.Field(i)
		name := field.Name
		m[name] = i
	}
	v.fieldMap[typ] = m
	return m
}

// Lookup to the value
func (v *vpath) Lookup(val interface{}, path string) []interface{} {
	ret := v.LookupBy([]reflect.Value{reflect.ValueOf(val)}, strings.Split(path, "/"))
	r := make([]interface{}, 0, len(ret))
	for _, v := range ret {
		r = append(r, v.Interface())
	}
	return r
}

// LookupBy to the value
func (v *vpath) LookupBy(vals []reflect.Value, path []string) (ret []reflect.Value) {
	path = clean(path)
	if len(path) == 0 {
		return vals
	}
	vs := strings.Split(path[0], ",")
	for _, val := range vals {
		ret = append(ret, v.lookup(val, vs)...)
	}
	return v.LookupBy(ret, path[1:])
}

func (v *vpath) lookup(val reflect.Value, vs []string) (ret []reflect.Value) {
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		l := val.Len()
		for _, v := range vs {
			if v == "*" {
				for i := 0; i != l; i++ {
					d := val.Index(i)
					ret = append(ret, d)
				}
				continue
			}
			sl := strings.Split(v, ":")
			switch len(sl) {
			case 1:
				index, err := strconv.ParseInt(v, 0, 0)
				if err != nil || int(index) >= l {
					continue
				}
				d := val.Index(int(index))
				ret = append(ret, d)
			case 2:
				offset := 0
				if sl[0] != "" {
					off, err := strconv.ParseInt(sl[0], 0, 0)
					if err != nil {
						continue
					}
					offset = int(off)
				}

				limit := l
				if sl[1] != "" {
					lim, err := strconv.ParseInt(sl[1], 0, 0)
					if err != nil {
						continue
					}
					limit = int(lim)
				}

				if limit >= l {
					limit = l
				}

				val = val.Slice(offset, limit)
				l = val.Len()
				for i := 0; i != l; i++ {

					d := val.Index(i)
					ret = append(ret, d)
				}
			}

		}

	case reflect.Map:
		typ := val.Type()
		switch typ.Key().Kind() {
		case reflect.String:
			for _, v := range vs {
				d := val.MapIndex(reflect.ValueOf(v))
				if d.Kind() != reflect.Invalid {
					ret = append(ret, d)
				}
			}
		default:
			//TODO
		}
	case reflect.Struct:
		typ := val.Type()
		m := v.field(typ)
		for _, v := range vs {
			i, ok := m[v]
			if ok {
				d := val.Field(i)
				ret = append(ret, d)
			}
		}
	case reflect.Ptr, reflect.Interface:
		if !val.IsNil() {
			return v.lookup(val.Elem(), vs)
		}
	}
	return ret
}

func clean(path []string) []string {
	ret := make([]string, 0, len(path))
	for i := 0; i != len(path); i++ {
		v := path[i]
		switch v {
		case ".":
		case "..":
			ret = ret[:len(ret)-1]
		default:
			ret = append(ret, v)
		}
	}
	return ret
}
