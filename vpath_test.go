package vpath

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestVpath(t *testing.T) {

	tests := []struct {
		src   string
		vpath string
		want  interface{}
	}{
		{`{"id":1,"data":"d"}`, `id,data`, []interface{}{1.0, "d"}},
		{`[{"id":1},{"id":"2"}]`, `*/id`, []interface{}{1.0, "2"}},
		{`[{"id":1},{"id":"2"},{"id":"3"},{"id":"4"}]`, `0,3/id`, []interface{}{1.0, "4"}},
		{`[{"id":1},{"id":"2"},{"id":"3"},{"id":"4"}]`, `3/id`, []interface{}{"4"}},
		{`[{"id":1},{"id":"2"},{"id":"3"},{"id":"4"}]`, `1:4/id`, []interface{}{"2", "3", "4"}},
		{`[{"list":[{"id":1},{"id":"2"}]},{"list":[{"id":3},{"id":"4"}]}]`, `*/list/*/id`, []interface{}{1.0, "2", 3.0, "4"}},
		{`[{"list":[{"id":1},{"id":"2"}]},{"list":[{"id":3},{"id":"4"}]}]`, `*/list/1/id`, []interface{}{"2", "4"}},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			var val interface{}
			err := json.Unmarshal([]byte(tt.src), &val)
			if err != nil {
				t.Errorf("Error json '%v'", tt.src)
			}
			got := Lookup(val, tt.vpath)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVpath() = %v, want %v", got, tt.want)
			}
		})
	}
}
