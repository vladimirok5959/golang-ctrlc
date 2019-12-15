package ctrlc

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const C_ICON_START = "🌟"
const C_ICON_WARN = "⚡️"
const C_ICON_HOT = "🔥"
const C_ICON_MAG = "✨"
const C_ICON_SC = "🌳"

type Iface interface {
	Shutdown(ctx context.Context) error
}

func App(t time.Duration, f func(ctx context.Context) *[]Iface) {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)

	fmt.Printf(
		C_ICON_START+" %s\n",
		cly(
			fmt.Sprintf(
				"Application started (%d sec)",
				t/time.Second,
			),
		),
	)

	sctx, shutdown := context.WithCancel(context.Background())
	ifaces := f(sctx)

	switch val := <-stop; val {
	case syscall.SIGINT:
		fmt.Printf(
			"\r"+C_ICON_WARN+" %s\n",
			cly(
				fmt.Sprintf(
					"Shutting down (interrupt) (%d sec)",
					t/time.Second,
				),
			),
		)
	case syscall.SIGTERM:
		fmt.Printf(
			C_ICON_WARN+" %s\n",
			cly(
				fmt.Sprintf(
					"Shutting down (terminate) (%d sec)",
					t/time.Second,
				),
			),
		)
	default:
		fmt.Printf(
			C_ICON_WARN+" %s\n",
			cly(
				fmt.Sprintf(
					"Shutting down (%d sec)",
					t/time.Second,
				),
			),
		)
	}

	shutdown()

	errors := false
	ctx, cancel := context.WithTimeout(context.Background(), t)
	for _, iface := range *ifaces {
		if err := iface.Shutdown(ctx); err != nil {
			errors = true
			fmt.Printf(
				C_ICON_HOT+" %s\n",
				clr(fmt.Sprintf(
					"Shutdown error (%T): %s",
					iface,
					err.Error(),
				)),
			)
		}
	}
	cancel()

	if errors {
		fmt.Printf(
			C_ICON_MAG+" %s\n",
			cly(
				fmt.Sprintf(
					"Application exited with errors (%d sec)",
					t/time.Second,
				),
			),
		)
		os.Exit(1)
	} else {
		fmt.Printf(
			C_ICON_SC+" %s\n",
			clg(
				fmt.Sprintf(
					"Application exited successfully (%d sec)",
					t/time.Second,
				),
			),
		)
	}
}