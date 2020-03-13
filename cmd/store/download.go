package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

// downloadCommand downloads a specific catalog.
type downloadCommand struct {
	verbose bool
	area    string
	outfile string
}

func init() {
	RegisterCommand("download", func(flags *flag.FlagSet) Command {
		cmd := new(downloadCommand)
		flags.BoolVar(&cmd.verbose, "v", false, "Print progress")
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
		"ABCDE12345 -v",
		"ABCDE12345 -o catalog.out",
	}
}

func (c *downloadCommand) Run(args []string) error {
	if len(args) != 1 {
		return errors.New("no pin specified")
	}

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

	csvw := csv.NewWriter(out)
	csvw.Comma = ';'
	csvw.UseCRLF = true
	_ = csvw.Write([]string{"Supplier SKU", "Name", "Price", "Price Qty", "Currency", "Order unit", "Manufacturer", "Manufacturer SKU", "GTIN/EAN"})

	var n int
	var pageToken string
	for {
		res, err := service.Scroll().PIN(args[0]).Area(c.area).PageToken(pageToken).Do(context.Background())
		if err != nil {
			return err
		}

		for _, item := range res.Items {
			n++

			csvw.Write([]string{
				item.Spn,
				item.Name,
				fmt.Sprintf("%.2f", item.Price),
				fmt.Sprintf("%.2f", item.PriceQty),
				item.Currency,
				item.OrderUnit,
				item.Manufacturer,
				item.Mpn,
				item.Gtin,
			})
		}

		if res.PageToken == "" {
			break
		}
		pageToken = res.PageToken
	}

	csvw.Flush()

	if c.verbose {
		fmt.Fprintf(os.Stdout, "Downloaded %d products\n", n)
	}

	return nil
}
