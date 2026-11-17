package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/daarxwalker/hui"
)

const usage = "usage: go tool hui assets [-dir DIR] [-check]"

func main() {
	if len(os.Args) < 2 || os.Args[1] != "assets" {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}
	flags := flag.NewFlagSet("assets", flag.ContinueOnError)
	dir := flags.String("dir", "hui", "directory to write to, relative to the go:generate directive")
	check := flags.Bool("check", false, "only report the files that differ, write nothing")
	if err := flags.Parse(os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}
	files := hui.Assets()
	if *check {
		changes, err := hui.CheckAssets(*dir, files)
		if err != nil {
			fmt.Fprintf(os.Stderr, "hui assets: %v\n", err)
			os.Exit(1)
		}
		var (
			foreign  int
			outdated int
		)
		for _, change := range changes {
			fmt.Fprintf(os.Stderr, "hui assets: %s/%s\n", *dir, change)
			if change.Kind == hui.AssetForeign {
				foreign++
				continue
			}
			outdated++
		}
		if foreign > 0 {
			fmt.Fprintf(os.Stderr, "hui assets: move or delete the foreign files, go generate refuses to write into %s while they are there\n", *dir)
		}
		if outdated > 0 {
			fmt.Fprintf(os.Stderr, "hui assets: %s is out of date, run go generate\n", *dir)
		}
		if len(changes) > 0 {
			os.Exit(1)
		}
		return
	}
	if err := hui.WriteAssets(*dir, files); err != nil {
		fmt.Fprintf(os.Stderr, "hui assets: %v\n", err)
		os.Exit(1)
	}
}
