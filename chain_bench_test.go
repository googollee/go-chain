package chain

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type httpError struct {
	code    int
	message string
}

func (e httpError) Error() string {
	return e.message
}

func (e httpError) Code() int {
	return e.code
}

type User struct {
	ID   int
	Name string
}

type UpdateUserArg struct {
	Name string
}

type Users struct{}

func (u *Users) Auth(ctx context.Context, r *http.Request) (*User, error) {
	return &User{ID: 1}, nil
}

func (u *Users) GetUpdateUserArg(ctx context.Context, r *http.Request) (arg UpdateUserArg, err error) {
	err = json.NewDecoder(r.Body).Decode(&arg)
	if err != nil {
		err = httpError{
			code:    http.StatusBadRequest,
			message: err.Error(),
		}
	}

	return
}

func (u *Users) UpdateUser(ctx context.Context, user *User, arg UpdateUserArg) (*User, error) {
	user.Name = arg.Name
	return user, nil
}

func RequestContext(r *http.Request) context.Context {
	ret := r.Context()
	if ret == nil {
		ret = context.Background()
	}
	return ret
}

type HTTPError interface {
	error
	Code() int
}

func Response[T any](w http.ResponseWriter, resp T, err error) {
	if err != nil {
		if httpErr, ok := err.(HTTPError); ok {
			w.WriteHeader(httpErr.Code())
			w.Write([]byte(httpErr.Error()))
			return
		}
	}

	json.NewEncoder(w).Encode(&resp)
}

func BenchmarkC(b *testing.B) {
	b.StopTimer()

	users := &Users{}
	f := C[http.HandlerFunc](RequestContext, func() (*User, error) { return nil, nil }, Defer(Response[*User]), users.Auth, users.GetUpdateUserArg, users.UpdateUser)
	arg := UpdateUserArg{
		Name: "name",
	}
	argData, _ := json.Marshal(&arg)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/users/1", bytes.NewReader(argData))
		resp := httptest.NewRecorder()

		f(resp, req)
	}
}

func BenchmarkFunc(b *testing.B) {
	b.StopTimer()

	users := &Users{}
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := RequestContext(r)

		var user *User
		var err error
		defer Response(w, user, err)

		user, err = users.Auth(ctx, r)
		if err != nil {
			return
		}

		arg, err := users.GetUpdateUserArg(ctx, r)
		if err != nil {
			return
		}

		user, err = users.UpdateUser(ctx, user, arg)
		if err != nil {
			return
		}
	}
	arg := UpdateUserArg{
		Name: "name",
	}
	argData, _ := json.Marshal(&arg)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/users/1", bytes.NewReader(argData))
		resp := httptest.NewRecorder()

		f(resp, req)
	}
}
