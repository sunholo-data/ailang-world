package store

import (
	"context"
	"reflect"
	"testing"
)

func TestNoDurableCompatibilityWrappers(t *testing.T) {
	typ := reflect.TypeOf(&Store{})
	for _, name := range []string{"Commit", "AppendIntent", "GetReceipt"} {
		method, ok := typ.MethodByName(name)
		if !ok {
			t.Errorf("missing %s", name)
			continue
		}
		if method.Type.NumIn() < 2 || method.Type.In(1) != reflect.TypeOf((*context.Context)(nil)).Elem() {
			t.Errorf("%s is not context-first", name)
		}
		if _, ok := typ.MethodByName(name + "Context"); ok {
			t.Errorf("compatibility wrapper %sContext remains", name)
		}
	}
}
