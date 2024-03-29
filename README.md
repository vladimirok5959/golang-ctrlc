# golang-ctrlc

Wrapper with context for safe application closing by signals like ctrl+c or kill

## How to use

```sh
go get github.com/vladimirok5959/golang-ctrlc
```

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/vladimirok5959/golang-ctrlc/ctrlc"
)

func main() {
	MyAppFunc := func(ctx context.Context, shutdown context.CancelFunc) *[]ctrlc.Iface {
		// Some custom logic
		// With goroutine inside
		test := Run()

		// err1 := fmt.Errorf("Startup error 1")
		// err2 := fmt.Errorf("Startup error 2")
		// return ctrlc.MakeError(shutdown, ctrlc.AppError(err1), ctrlc.AppError(err2))

		// Http web server
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("New web request (%s)!\n", r.URL.Path)

			// Do something hard inside (12 seconds)
			for i := 0; i < 12000; i++ {
				select {
				case <-ctx.Done():
					// Interrupt request by server
					fmt.Printf("[BY SERVER] OK, I will cancel (%s)!\n", r.URL.Path)
					return
				case <-r.Context().Done():
					// Interrupt request by client
					fmt.Printf("[BY CLIENT] OK, I will cancel (%s)!\n", r.URL.Path)
					return
				default:
					// Main some logic
					// Some very long logic, just for example
					time.Sleep(1 * time.Millisecond)
				}
			}

			fmt.Printf("After 12 seconds!\n")
			w.Header().Set("Content-Type", "text/html")
			if _, err := w.Write([]byte(`<div>After 12 seconds!</div>`)); err != nil {
				fmt.Printf("%s\n", err.Error())
				return
			}
			if _, err := w.Write([]byte(`<div>` + r.URL.Path + `</div>`)); err != nil {
				fmt.Printf("%s\n", err.Error())
				return
			}
		})
		srv := &http.Server{Addr: "127.0.0.1:8080", Handler: mux}
		go func() {
			fmt.Printf("Starting web server: http://127.0.0.1:8080/\n")
			if err := srv.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					fmt.Printf("Web server startup error: %s\n", err.Error())
					// Application can't working without http web server
					// Call cancel context func on error
					shutdown()
					return
				}
			}
		}()

		return &[]ctrlc.Iface{test, srv}
	}

	// Run application
	ctrlc.App(MyAppFunc)
}
```

## Explanation

```go
type Iface interface {
	Shutdown(ctx context.Context) error
}

type CallbackFunc func(ctx context.Context, shutdown context.CancelFunc) *[]Iface

func App(f CallbackFunc)
func AppWithTimeOut(t time.Duration, f CallbackFunc)
```

**t** in `App` function 8 seconds is a default timeout for `Shutdown` function, if Shutdown function will be not closed, context will be canceled aftet this value and will print error to console. **f** - is function for main application code with cancel context inside, `ctx` will be triggered when application got one of close/terminate signals and we can terminate application by calling `shutdown` function

Callback function `f` in `App` must return pointer to array of interfaces which must contain/implement `Shutdown(ctx context.Context) error` function for correct shutdown
