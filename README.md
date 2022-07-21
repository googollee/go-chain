[![Unit tests](https://github.com/googollee/go-chain/actions/workflows/unittest.yaml/badge.svg)](https://github.com/googollee/go-chain/actions/workflows/unittest.yaml)

# go-chain

go-chain is a dependency injection to create a function to call a group of function orderly.

# Use case

## Middleware for HTTP endpoints

```go
func RequestContext(r *http.Request) context.Context {
  return r.Context()
}

func Auth(ctx context.Context, r *http.Request) error {
  return nil
}

func GetBody[T any] (r *http.Request) (ret T, err error) {
  err = json.NewDecoder(r.Body).Decode(&ret)
  return
}

func Handler(ctx context.Context, arg int) (string, error) {
  return fmt.Sprintf("%d", arg), nil
}

func Return[T any] (w http.ResponseWriter, arg T, err error) {
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  json.NewEncoder(w).Encode(arg)
}

func Server() {
  http.HandleFunc("/", chain.C[func(w http.ResponseWriter, r *http.Request)](
    RequestContext, Auth, GetBody[int], Handler, Return[string]
  ))
}
```

The function generating by `chain.C()` is similar to:

```go
func (w http.ResponseWriter, r *http.Request) {
  ctx := RequestContext(r)

  err := Auth(ctx, r)

  var i int
  if err == nil {
    i, err = GetBody[int](r)
  }

  var s string
  if err == nil {
    s, err = Handler(ctx, i)
  }

  Return[string](w, s, err)
}
```