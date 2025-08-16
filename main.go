package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"

	"github.com/ainvaltin/nu-plugin"
)

func main() {
	err := run()
	if errors.Is(err, nu.ErrGoodbye) {
		return
	}
	if err != nil {
		slog.Error("plugin run", "err", err)
		os.Exit(1)
	}
}

var commands []*nu.Command

func run() (err error) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	p, err := nu.New(commands, "0.0.1", nil)
	if err != nil {
		return
	}

	err = p.Run(ctx)
	return
}
