package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// downloadCommand downloads a specific catalog.
type downloadCommand struct {
	verbose bool
	pin     string
	area    string
	outfile string
}

func init() {
	RegisterCommand("download", func(flags *flag.FlagSet) Command {
		cmd := new(downloadCommand)
		flags.BoolVar(&cmd.verbose, "v", false, "Print progress")
		flags.StringVar(&cmd.pin, "pin", "", "PIN of catalog")
		flags.StringVar(&cmd.area, "area", "live", "Area to download (work/live)")
		flags.StringVar(&cmd.outfile, "o", "", "Output file")
		return cmd
	})
}

func (c *downloadCommand) Describe() string {
	return "Downloads a catalog."
}

func (c *downloadCommand) Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s download\n", os.Args[0])
}

func (c *downloadCommand) Examples() []string {
	return []string{
		"-pin=ABCDE12345 -v",
		"-pin=ABCDE12345 -o catalog.out",
	}
}

func (c *downloadCommand) Run(args []string) error {
	service, err := GetProductsService()
	if err != nil {
		return err
	}

	var out io.Writer
	if c.outfile != "" {
		f, err := os.OpenFile(c.outfile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	} else {
		out = os.Stdout
	}

	var n int
	var pageToken string
	for {
		res, err := service.Scroll().PIN(c.pin).Area(c.area).PageToken(pageToken).Do()
		if err != nil {
			return err
		}

		for _, item := range res.Items {
			n++
			fmt.Fprintf(out, "%v\n", item)
		}

		if res.PageToken == "" {
			break
		}
		pageToken = res.PageToken
	}

	if c.verbose {
		fmt.Fprintf(os.Stdout, "Downloaded %d products\n", n)
	}

	return nil
}
