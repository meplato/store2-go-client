package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/meplato/store2-go-client/products"
)

// uploadCommand uploads to a specific catalog.
type uploadCommand struct {
	verbose bool
	infile  string
	outfile string
}

func init() {
	RegisterCommand("upload", func(flags *flag.FlagSet) Command {
		cmd := new(uploadCommand)
		flags.BoolVar(&cmd.verbose, "v", false, "Print progress")
		flags.StringVar(&cmd.infile, "i", "", "Input file")
		return cmd
	})
}

func (c *uploadCommand) Describe() string {
	return "Upload a catalog."
}

func (c *uploadCommand) Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s upload <pin> < filename.csv\n", os.Args[0])
	fmt.Fprint(os.Stderr, `
The uploaded file must be in CSV format with a semicolon as a separator
and (optionally) enclosed by double-quotes. All rows in the CSV file must
have the same number of columns.

The first line is the header line and must include one or more of the
following columns: MODE, SPN, NAME, PRICE, ORDER_UNIT, MPN, MANUFACTURER,
ECLASS_VERSION, ECLASS_CODE, and TAX_CODE.
The header row must have the two columns MODE and SPN.

The MODE column of each row must have one of the following values:
C - The product should be created. The row must have the columns
    NAME, PRICE, and ORDER_UNIT.
D - The product should be deleted.
U - The product should be updated. All columns with a non-blank
      value will be updated.

Example file:

MODE;SPN;NAME;PRICE;ORDER_UNIT
C;1000;"Product 1000";19.50;PCE
C;2000;"Product 2000";0.50;PCE
U;2000;;0.49;EA
D;1000;;;

Upload will read the file line by line. It will first try to insert the
product with supplier part number (SPN) 1000. Then it will insert the
product with SPN 2000. The 4th row will update the price and the order unit
of product 2000 to 0.49 and EA respectively. Finally, the product 1000 is
deleted from the catalog.

Final notes:

The upload command is a very simple example to illustrate interacting with
a catalog in Store. You must probably use a more complete implementation
that e.g. does not interact with CSV files but with your product database.
However, the Store API remains unchanged regardless of the backend.

`)
}

func (c *uploadCommand) Examples() []string {
	return []string{
		"ABCDE12345 -v < catalogfile.csv",
		"ABCDE12345 -i catalogdata.csv",
	}
}

func (c *uploadCommand) Run(args []string) error {
	if len(args) != 1 {
		return errors.New("no pin specified")
	}

	pin := args[0]

	service, err := GetProductsService()
	if err != nil {
		return err
	}

	// Prepare input
	var in io.Reader
	if c.infile != "" {
		f, err := os.Open(c.infile)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	} else {
		in = os.Stdin
	}
	csvr := csv.NewReader(in)
	csvr.Comma = ';'

	// Parse header from input and initialize cell handlers
	header, err := csvr.Read()
	if err != nil {
		return err
	}
	if len(header) == 0 {
		return errors.New("no header row")
	}
	handlersByIndex := make(map[int]rowHandler)
	for i, cell := range header {
		h, found := rowHandlers[cell]
		if !found {
			return fmt.Errorf("found invalid column name %q", cell)
		}
		handlersByIndex[i] = h
	}

	// Read input file line-by-line
	var line int = 1
	for {
		record, err := csvr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		line++

		var r row
		r.Line = line

		for i, cell := range record {
			h, found := handlersByIndex[i]
			if !found {
				return fmt.Errorf("no handler for index %d", i)
			}
			if err := h(&r, cell); err != nil {
				return fmt.Errorf("line %d: %v", line, err)
			}
		}

		if c.verbose {
			fmt.Fprintf(os.Stdout, "line %6d\r", line)
		}

		// Validate the row
		if err := r.Validate(); err != nil {
			return fmt.Errorf("line %d: %v", err)
		}

		// Call Create, Update, or Delete API
		switch r.Mode {
		case "C":
			// Create a new product (or overwrite an existing)
			p := &products.CreateProduct{
				Spn:       r.SPN,
				Name:      *r.Name,
				Price:     *r.Price,
				OrderUnit: *r.OrderUnit,
			}
			if r.MPN != nil {
				p.Mpn = *r.MPN
			}
			if r.Manufacturer != nil {
				p.Manufacturer = *r.Manufacturer
			}
			if r.EclassVersion != nil && r.EclassCode != nil {
				p.Eclasses = append(p.Eclasses, &products.Eclass{
					Version: *r.EclassVersion,
					Code:    *r.EclassCode,
				})
			}
			if r.TaxCode != nil {
				p.TaxCode = *r.TaxCode
			}
			_, err := service.Create().PIN(pin).Area("work").Product(p).Do()
			if err != nil {
				return fmt.Errorf("line %d: create failed: %v", err)
			}
		case "U":
			// Update a product
			p := &products.UpdateProduct{
				Name:         r.Name,
				Price:        r.Price,
				OrderUnit:    r.OrderUnit,
				Mpn:          r.MPN,
				Manufacturer: r.Manufacturer,
				TaxCode:      r.TaxCode,
			}
			if r.EclassVersion != nil && r.EclassCode != nil {
				p.Eclasses = append(p.Eclasses, &products.Eclass{
					Version: *r.EclassVersion,
					Code:    *r.EclassCode,
				})
			}
			_, err := service.Update().PIN(pin).Area("work").Spn(r.SPN).Product(p).Do()
			if err != nil {
				return fmt.Errorf("line %d: update failed: %v", err)
			}
		case "D":
			// Delete a product
			err := service.Delete().PIN(pin).Area("work").Spn(r.SPN).Do()
			if err != nil {
				return fmt.Errorf("line %d: delete failed: %v", err)
			}
		}
	}

	if c.verbose {
		fmt.Fprintf(os.Stdout, "Read %d lines\n", line)
	}

	return nil
}

// row is an intermediary structure to read data into.
type row struct {
	Line          int
	Mode          string
	SPN           string
	Name          *string
	Price         *float64
	OrderUnit     *string
	MPN           *string
	Manufacturer  *string
	EclassVersion *string
	EclassCode    *string
	TaxCode       *string
}

// Validate checks for errors in a row. It also ensures that the given
// fields are valid with regard to the mode.
func (r *row) Validate() error {
	if !(r.Mode == "C" || r.Mode == "U" || r.Mode == "D") {
		return fmt.Errorf("unknown mode %q", r.Mode)
	}
	if r.SPN == "" {
		return errors.New("no SPN specified")
	}
	if r.Mode == "C" {
		if r.Name == nil || *r.Name == "" {
			return errors.New("no name specified")
		}
		if r.Price == nil || *r.Price < 0.0 {
			return errors.New("no price specified")
		}
		if r.OrderUnit == nil || *r.OrderUnit == "" {
			return errors.New("no order unit specified")
		}
	}
	return nil
}

// rowHandler handles the update of a specific cell and writes the parsed
// value into the field of a row.
type rowHandler func(r *row, cell string) error

// rowHandlers by column name.
var rowHandlers = map[string]rowHandler{
	"MODE":           handleMode,
	"SPN":            handleSPN,
	"NAME":           handleName,
	"PRICE":          handlePrice,
	"ORDER_UNIT":     handleOrderUnit,
	"MPN":            handleMPN,
	"MANUFACTURER":   handleManufacturer,
	"ECLASS_VERSION": handleEclassVersion,
	"ECLASS_CODE":    handleEclassCode,
	"TAX_CODE":       handleTaxCode,
}

func handleMode(r *row, cell string) error {
	r.Mode = strings.ToUpper(cell)
	return nil
}

func handleSPN(r *row, cell string) error {
	r.SPN = cell
	return nil
}

func handleName(r *row, cell string) error {
	if cell != "" {
		r.Name = &cell
	}
	return nil
}

func handlePrice(r *row, cell string) error {
	if cell != "" {
		if price, err := strconv.ParseFloat(cell, 64); err != nil {
			return fmt.Errorf("price %q is not a number", cell)
		} else {
			r.Price = &price
		}
	}
	return nil
}

func handleOrderUnit(r *row, cell string) error {
	if cell != "" {
		r.OrderUnit = &cell
	}
	return nil
}

func handleMPN(r *row, cell string) error {
	if cell != "" {
		r.MPN = &cell
	}
	return nil
}

func handleManufacturer(r *row, cell string) error {
	if cell != "" {
		r.Manufacturer = &cell
	}
	return nil
}

func handleEclassVersion(r *row, cell string) error {
	if cell != "" {
		r.EclassVersion = &cell
	}
	return nil
}

func handleEclassCode(r *row, cell string) error {
	if cell != "" {
		r.EclassCode = &cell
	}
	return nil
}

func handleTaxCode(r *row, cell string) error {
	if cell != "" {
		r.TaxCode = &cell
	}
	return nil
}
