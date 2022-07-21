// Package chain provides a injecting way to call a group of functions orderly.
package chain

import (
	"fmt"
	"reflect"
)

// C creates a new function as the type of Func, which calls given functions fn orderly.
// Each fn will be called with inputs of Func or returns from previous fn, picking with the right type.
// Returns of Func are picking from inputs of Func and all returns of fn.
func C[Func any](fn ...any) Func {
	var f *Func
	tfunc := reflect.TypeOf(f).Elem()

	chain := newChain(tfunc, fn)
	chain.Check()

	ret := reflect.MakeFunc(tfunc, chain.Call)

	return ret.Interface().(Func)
}

type chain struct {
	inputs  []reflect.Type
	outputs []reflect.Type
	funcs   []reflect.Value
}

func newChain(funcType reflect.Type, funcs []any) *chain {
	if funcType.Kind() != reflect.Func {
		panic("chain.C should be instanced by a function.")
	}

	ret := &chain{}

	ret.inputs = make([]reflect.Type, funcType.NumIn())
	for i := range ret.inputs {
		ret.inputs[i] = funcType.In(i)
	}

	ret.outputs = make([]reflect.Type, funcType.NumOut())
	for i := range ret.outputs {
		ret.outputs[i] = funcType.Out(i)
	}

	ret.funcs = make([]reflect.Value, len(funcs))
	for i := range funcs {
		ret.funcs[i] = reflect.ValueOf(funcs[i])
	}

	return ret
}

func (c *chain) Check() {
	providers := make(map[reflect.Type]bool)
	for _, input := range c.inputs {
		providers[input] = true
	}

	for _, fn := range c.funcs {
		ft := fn.Type()
		for i := 0; i < ft.NumIn(); i++ {
			in := ft.In(i)
			if !providers[in] {
				panic(fmt.Sprintf("chain.C can't provide any instance of type %s for function %s as input", in, ft))
			}
		}

		for i := 0; i < ft.NumOut(); i++ {
			out := ft.Out(i)
			providers[out] = true
		}
	}

	for _, out := range c.outputs {
		if !providers[out] {
			panic(fmt.Sprintf("chain.C can't provide any instance of type %s as output", out))
		}
	}
}

func (c *chain) Call(args []reflect.Value) []reflect.Value {
	providers := make(map[reflect.Type]reflect.Value)

	for _, arg := range args {
		providers[arg.Type()] = arg
	}

	for _, fn := range c.funcs {
		ft := fn.Type()

		input := make([]reflect.Value, ft.NumIn())
		for i := range input {
			input[i] = providers[ft.In(i)]
		}

		output := fn.Call(input)
		for _, out := range output {
			providers[out.Type()] = out
		}
	}

	output := make([]reflect.Value, len(c.outputs))
	for i := range output {
		output[i] = providers[c.outputs[i]]
	}

	return output
}
