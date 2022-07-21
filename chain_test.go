package chain

import "testing"

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
