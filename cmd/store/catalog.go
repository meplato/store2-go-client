package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// catalogCommand gets details about one catalog.
type catalogCommand struct {
	pin string
}

func init() {
	RegisterCommand("catalog", func(flags *flag.FlagSet) Command {
		cmd := new(catalogCommand)
		return cmd
	})
}

func (c *catalogCommand) Describe() string {
	return "Print catalog information."
}

func (c *catalogCommand) Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s catalog <pin>\n", os.Args[0])
}

func (c *catalogCommand) Examples() []string {
	return []string{
		"ABCDE12345",
		"ABCDE12345 BEEF1C0DE1",
	}
}

func (c *catalogCommand) Run(args []string) error {
	if len(args) == 0 {
		return errors.New("no pin specified")
	}

	service, err := GetCatalogsService()
	if err != nil {
		return err
	}

	for i, pin := range args {
		c, err := service.Get().PIN(pin).Do()
		if err != nil {
			return err
		}

		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("%20s: %s\n", "PIN", c.PIN)
		fmt.Printf("%20s: %s\n", "Name", c.Name)
		fmt.Printf("%20s: %v\n", "Created", c.Created)
		if c.NumProductsWork != nil {
			fmt.Printf("%20s: %d\n", "# products work", *c.NumProductsWork)
		} else {
			fmt.Printf("%20s: %d\n", "# products work", 0)
		}
		if c.NumProductsLive != nil {
			fmt.Printf("%20s: %d\n", "# products live", *c.NumProductsLive)
		} else {
			fmt.Printf("%20s: %d\n", "# products live", 0)
		}
	}

	return nil
}
