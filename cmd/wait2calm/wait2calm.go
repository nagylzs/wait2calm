package main

import (
	_ "embed"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/nagylzs/wait2calm/internal/config"
	"github.com/nagylzs/wait2calm/internal/signal"
	"github.com/nagylzs/wait2calm/internal/version"
	"github.com/shirou/gopsutil/v4/load"
)

const Granularity = time.Millisecond * 100

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

	var timedOut bool
	go func() {
		err, timedOut = runMain(opts, posArgs)
		if err != nil {
			signal.Stop(2)
			return
		}
		if timedOut && !opts.SuccessOnTimeout {
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
	}
	signal.Exit()
}

// waitSome waits for w amount of time, or less, if doNotWaitAfter has elapsed since started
// returns true if it actually waited w, returns false otherwise
func waitSome(w time.Duration, started time.Time, doNotWaitAfter time.Duration) bool {
	for w > 0 {
		elapsed := time.Since(started)
		if doNotWaitAfter > 0 && elapsed > doNotWaitAfter {
			slog.Warn("Do not wait anymore",
				"--do-not-wait-after", FormatFluent(doNotWaitAfter),
				"elapsed", FormatFluent(elapsed))
			return false
		}
		if w > Granularity {
			time.Sleep(Granularity)
			w -= Granularity
		} else {
			time.Sleep(w)
			break
		}
	}
	return true
}

func FormatFluent(d time.Duration) string {
	switch {
	case d >= time.Hour:
		return fmt.Sprintf("%.2fh", float64(d)/float64(time.Hour))
	case d >= time.Minute:
		return fmt.Sprintf("%.2fm", float64(d)/float64(time.Minute))
	case d >= time.Second:
		return fmt.Sprintf("%.2fs", float64(d)/float64(time.Second))
	case d >= time.Millisecond:
		return fmt.Sprintf("%.2fms", float64(d)/float64(time.Millisecond))
	case d >= time.Microsecond:
		return fmt.Sprintf("%.2fµs", float64(d)/float64(time.Microsecond))
	default:
		return fmt.Sprintf("%dns", d.Nanoseconds())
	}
}

// returns an error, and a flag that is set when it times out on --do-not-wait-after
func runMain(opts config.Opts, posArgs []string) (error, bool) {
	if opts.LoadType != 1 && opts.LoadType != 5 && opts.LoadType != 15 {
		return fmt.Errorf("invalid load type %d must be 1,5 or 15", opts.LoadType), false
	}
	if opts.ImmediateStartBelow == 0 && opts.DelayedStartBelow == 0 {
		slog.Warn("Warning, both ImmediateStartBelow and DelayedStartBelow are zero!")
	}
	if opts.ImmediateStartBelow < 0 || opts.DelayedStartBelow < 0 {
		return fmt.Errorf("both immediate and delayed 'start below' must not be negative"), false
	}
	if opts.ImmediateStartBelow > opts.DelayedStartBelow {
		slog.Warn("ImmediateStartBelow > DelayedStartBelow, are you sure?")
	}

	if opts.WaitBefore > opts.DoNotWaitAfter {
		slog.Warn("WaitBefore > DoNotWaitAfter? -> DoNotWaitAfter may be ignored")
	}

	started := time.Now()

	if opts.WaitBefore != time.Duration(0) {
		w := time.Duration(rand.Float64() * float64(opts.WaitBefore))
		slog.Info("Wait before first measurement",
			"max", FormatFluent(opts.WaitBefore),
			"actual", FormatFluent(w))
		if !waitSome(w, started, opts.DoNotWaitAfter) {
			return nil, true
		}
	}

	prevBelow := false
	for {
		info, err := load.Avg()
		if err != nil {
			return fmt.Errorf("load.Avg() failed: %v", err), false
		}
		var load float64
		switch opts.LoadType {
		case 1:
			load = info.Load1
		case 5:
			load = info.Load5
		case 15:
			load = info.Load15
		}
		if load <= opts.ImmediateStartBelow {
			slog.Info("Immediate start",
				"load", load,
				"immediate-start-below", opts.ImmediateStartBelow)
			break
		}
		if load < opts.DelayedStartBelow && prevBelow {
			slog.Info("Delayed start (second measurement)",
				"load", load,
				"delayed-start-below", opts.DelayedStartBelow)
			break
		}
		if load < opts.DelayedStartBelow {
			prevBelow = true
			slog.Info("Delayed start (first measurement)",
				"load", load,
				"delayed-start-below", opts.DelayedStartBelow)
		} else {
			prevBelow = false
			slog.Info("Do not start yet",
				"load", load,
				"delayed-start-below", opts.DelayedStartBelow)
		}
		slog.Info("Wait before next measurement",
			"measurement-interval", FormatFluent(opts.MeasurementInterval),
			"elapsed", FormatFluent(time.Since(started)))

		if !waitSome(opts.MeasurementInterval, started, opts.DoNotWaitAfter) {
			return nil, true
		}
	}

	if opts.WaitAfter != time.Duration(0) {
		w := time.Duration(rand.Float64() * float64(opts.WaitAfter))
		slog.Info("Wait after calm down",
			"max", FormatFluent(opts.WaitAfter),
			"actual", FormatFluent(w))
		if !waitSome(w, started, opts.DoNotWaitAfter) {
			return nil, true
		}
	}

	slog.Info("Calmed down!")
	return nil, false
}
