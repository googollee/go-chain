package chain

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
)

func TestC(t *testing.T) {
	type MyString string
	type MyInt int

	handleString := func(str MyString) string {
		return string(str)
	}

	handleInt := func(i MyInt) int {
		return int(i)
	}

	f := C[func(MyString, MyInt) (string, int)](handleString, handleInt)

	str := MyString("str")
	i := MyInt(1)

	got1, got2 := f(str, i)

	if got1 != "str" || got2 != 1 {
		t.Errorf("want: str, 1, got: %v, %v", got1, got2)
	}
}

func TestCError(t *testing.T) {
	wantErr := errors.New("error")
	errFunc := func(i int) error {
		if i < 10 {
			return wantErr
		}
		return nil
	}

	followCalled := false
	followFunc := func() {
		followCalled = true
	}

	f := C[func(int) error](errFunc, followFunc)

	tests := []struct {
		input      int
		wantErr    error
		wantCalled bool
	}{
		{0, wantErr, false},
		{100, nil, true},
	}

	for _, test := range tests {
		name := fmt.Sprintf("C(errFunc, followFunc)(%d)", test.input)
		t.Run(name, func(t *testing.T) {
			followCalled = false
			gotErr := f(test.input)

			if want, got := test.wantErr, gotErr; want != got {
				t.Errorf("return error, want: %v, got: %v", want, got)
			}

			if want, got := test.wantCalled, followCalled; want != got {
				t.Errorf("follow func calls, want: %v, got: %v", want, got)
			}
		})
	}
}

func TestCDefer(t *testing.T) {
	var gotErr error
	var gotCode int

	init := func() (int, error) {
		return 0, nil
	}

	errHandler := func(code int, err error) {
		if err != nil {
			gotErr = err
			return
		}

		gotCode = code
	}

	atoi := func(s string) (int, error) {
		ret, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return int(ret), nil
	}

	f := C[func(string)](init, Defer(errHandler), atoi)

	tests := []struct {
		input    string
		wantCode int
		wantErr  bool
	}{
		{"10", 10, false},
	}

	for _, test := range tests {
		name := fmt.Sprintf("C(Defer(errHandler), atoi)(%q)", test.input)
		t.Run(name, func(t *testing.T) {
			gotCode = 0
			gotErr = nil

			f(test.input)

			if test.wantErr != (gotErr != nil) {
				t.Errorf("want error: %v, got error: %v", test.wantErr, gotErr)
			}

			if want, got := test.wantCode, gotCode; want != got {
				t.Errorf("code, want: %v, got: %v", want, got)
			}
		})
	}
}
