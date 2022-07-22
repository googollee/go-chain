package chain

import "reflect"

// Defer injects a function defering to be called just before the call chain finishing.
func Defer(fn any) deferFunc {
	return deferFunc{
		value: reflect.ValueOf(fn),
	}
}

type deferFunc struct {
	value reflect.Value
}

var deferType = reflect.TypeOf((*deferFunc)(nil)).Elem()
