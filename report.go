package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	vegeta "github.com/tsenart/vegeta/lib"
)

func reportCmd() command {
	fs := flag.NewFlagSet("vegeta report", flag.ExitOnError)
	reporter := fs.String("reporter", "text", "Reporter [text, json, plot, hist[buckets]]")
	window := fs.Uint64("window", ^uint64(0), "Window size to aggregate")
	inputs := fs.String("inputs", "stdin", "Input files (comma separated)")
	output := fs.String("output", "stdout", "Output file")
	return command{fs, func(args []string) error {
		fs.Parse(args)
		return report(*reporter, *inputs, *output, *window)
	}}
}

// report validates the report arguments, sets up the required resources
// and writes the report every window duration
func report(reporter, inputs, output string, window uint64) error {
	if len(reporter) < 4 {
		return fmt.Errorf("bad reporter: %s", reporter)
	}

	files := strings.Split(inputs, ",")
	srcs := make([]io.Reader, len(files))
	for i, f := range files {
		in, err := file(f, false)
		if err != nil {
			return err
		}
		defer in.Close()
		srcs[i] = in
	}
	dec := vegeta.NewDecoder(srcs...)

	out, err := file(output, true)
	if err != nil {
		return err
	}
	defer out.Close()

	if window == 0 {
		return errors.New("bad window: got 0, want: [1..]")
	}

	var rep vegeta.Reporter
	switch reporter[:4] {
	case "text":
		rep = &vegeta.TextReporter{}
	case "json":
		rep = &vegeta.JSONReporter{}
	case "plot":
		if want := ^uint64(0); window != want {
			return fmt.Errorf("bad window for plot report: must be default %d", want)
		}
		rep = &vegeta.PlotReporter{}
	case "hist":
		if len(reporter) < 6 {
			return fmt.Errorf("bad buckets: '%s'", reporter[4:])
		}
		var hr vegeta.HistogramReporter
		if err := hr.Buckets.UnmarshalText([]byte(reporter[4:])); err != nil {
			return err
		}
		rep = &hr
	}

	return vegeta.Report(dec, rep, out, window)
}
