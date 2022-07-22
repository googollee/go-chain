package chain

import "reflect"

type DeferFunc struct {
	value reflect.Value
}

var deferType = reflect.TypeOf((*DeferFunc)(nil)).Elem()

func Defer(fn any) DeferFunc {
	return DeferFunc{
		value: reflect.ValueOf(fn),
	}
}
