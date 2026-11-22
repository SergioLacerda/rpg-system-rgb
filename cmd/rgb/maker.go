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
	OutDir    string
}

func runMaker(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing maker subcommand (want process|stats|export)")
	}
	switch args[0] {
	case "process":
		return runMakerProcess(args[1:])
	case "stats":
		return runMakerStats(args[1:])
	case "export":
		return runMakerExport(args[1:])
	default:
		return fmt.Errorf("unknown maker subcommand %q (want process|stats|export)", args[0])
	}
}

func runMakerProcess(args []string) error {
	opts, err := parseMakerProcessFlags(args)
	if err != nil {
		return err
	}
	out, err := app.ProcessMakerBatch(opts.BatchFile, opts.IndexFile, opts.OutDir)
	if err != nil {
		return err
	}
	fmt.Print(out)
	return nil
}

func runMakerStats(args []string) error {
	dir, err := parseMakerStatsFlags(args)
	if err != nil {
		return err
	}
	out, err := app.MakerWorkspaceStats(dir)
	if err != nil {
		return err
	}
	fmt.Print(out)
	return nil
}

func runMakerExport(args []string) error {
	dir, outDir, err := parseMakerExportFlags(args)
	if err != nil {
		return err
	}
	out, err := app.ExportMakerChatPack(dir, outDir)
	if err != nil {
		return err
	}
	fmt.Print(out)
	return nil
}

func parseMakerExportFlags(args []string) (string, string, error) {
	var dir, outDir string
	flags := flag.NewFlagSet("maker export", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&dir, "dir", "", "path to Maker workspace directory")
	flags.StringVar(&outDir, "out", "", "path to output chat-pack directory")
	if err := flags.Parse(args); err != nil {
		return "", "", err
	}
	if flags.NArg() != 0 {
		return "", "", fmt.Errorf("unexpected maker export arguments: %v", flags.Args())
	}
	if dir == "" || outDir == "" {
		return "", "", fmt.Errorf("maker export requires --dir and --out")
	}
	return dir, outDir, nil
}

func parseMakerStatsFlags(args []string) (string, error) {
	var dir string
	flags := flag.NewFlagSet("maker stats", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&dir, "dir", "", "path to Maker workspace directory")
	if err := flags.Parse(args); err != nil {
		return "", err
	}
	if flags.NArg() != 0 {
		return "", fmt.Errorf("unexpected maker stats arguments: %v", flags.Args())
	}
	if dir == "" {
		return "", fmt.Errorf("maker stats requires --dir")
	}
	return dir, nil
}

func parseMakerProcessFlags(args []string) (makerProcessOptions, error) {
	var opts makerProcessOptions
	flags := flag.NewFlagSet("maker process", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&opts.BatchFile, "batch", "", "path to Maker batch JSON file")
	flags.StringVar(&opts.IndexFile, "index", "", "optional path to semantic index JSON file")
	flags.StringVar(&opts.OutDir, "out-dir", "", "optional directory to write rendered package files")
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
