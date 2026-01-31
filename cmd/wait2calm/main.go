package main

import (
	_ "embed"
	"log/slog"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/nagylzs/wait2calm/internal/config"
	"github.com/nagylzs/wait2calm/internal/signal"
	"github.com/nagylzs/wait2calm/internal/version"
)

func main() {
	var opts = config.Opts{
		Debug:   false,
		Verbose: false,
	}
	posArgs, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	if opts.ShowVersion {
		version.PrintVersion()
		os.Exit(0)
	}

	// Set loglevel
	var programLevel = new(slog.LevelVar)
	if opts.Debug {
		programLevel.Set(slog.LevelDebug)
	} else if opts.Verbose {
		programLevel.Set(slog.LevelInfo)
	} else {
		programLevel.Set(slog.LevelWarn)
	}

	lw := os.Stderr
	h := slog.New(
		tint.NewHandler(lw, &tint.Options{
			NoColor: !isatty.IsTerminal(lw.Fd()),
			Level:   programLevel,
		}),
	)
	slog.SetDefault(h)

	signal.SetupSignalHandler()

	go func() {
		err = runMain(opts, posArgs)
		if err != nil {
			signal.Stop(1)
			return
		}
		signal.Stop(0)
	}()

	for !signal.IsStopping() {
		time.Sleep(time.Second)
	}

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func runMain(opts config.Opts, posArgs []string) error {
	slog.Info("TODO")
	return nil
}
