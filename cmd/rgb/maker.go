package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/SergioLacerda/rpg-system-rgb/internal/app"
)

type makerProcessOptions struct {
	BatchFile string
	IndexFile string
}

func runMaker(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing maker subcommand (want process)")
	}
	switch args[0] {
	case "process":
		opts, err := parseMakerProcessFlags(args[1:])
		if err != nil {
			return err
		}
		out, err := app.ProcessMakerBatch(opts.BatchFile, opts.IndexFile)
		if err != nil {
			return err
		}
		fmt.Print(out)
		return nil
	default:
		return fmt.Errorf("unknown maker subcommand %q (want process)", args[0])
	}
}

func parseMakerProcessFlags(args []string) (makerProcessOptions, error) {
	var opts makerProcessOptions
	flags := flag.NewFlagSet("maker process", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&opts.BatchFile, "batch", "", "path to Maker batch JSON file")
	flags.StringVar(&opts.IndexFile, "index", "", "optional path to semantic index JSON file")
	if err := flags.Parse(args); err != nil {
		return makerProcessOptions{}, err
	}
	if flags.NArg() != 0 {
		return makerProcessOptions{}, fmt.Errorf("unexpected maker process arguments: %v", flags.Args())
	}
	if opts.BatchFile == "" {
		return makerProcessOptions{}, fmt.Errorf("maker process requires --batch")
	}
	return opts, nil
}
