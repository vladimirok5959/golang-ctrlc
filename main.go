package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/vladimirok5959/golang-ctrlc/ctrlc"
)

func main() {
	ctrlc.App(8*time.Second, func(ctx context.Context) *[]ctrlc.Iface {
		// Some custom logic
		test := Run()

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
					time.Sleep(1 * time.Millisecond)
				}
			}

			fmt.Printf("After 12 seconds!\n")
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<div>After 12 seconds!</div>`))
			w.Write([]byte(`<div>` + r.URL.Path + `</div>`))
		})
		srv := &http.Server{Addr: "127.0.0.1:8080", Handler: mux}
		go func() {
			fmt.Printf("Starting web server: http://127.0.0.1:8080/\n")
			if err := srv.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					fmt.Printf("Web server startup error: %s\n", err.Error())
					return
				}
			}
		}()

		return &[]ctrlc.Iface{test, srv}
	})
}
