package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// publishCommand publishes a catalog.
type publishCommand struct {
}

func init() {
	RegisterCommand("publish", func(flags *flag.FlagSet) Command {
		cmd := new(publishCommand)
		return cmd
	})
}

func (c *publishCommand) Describe() string {
	return "Publish a catalog."
}

func (c *publishCommand) Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s publish <pin>\n", os.Args[0])
}

func (c *publishCommand) Examples() []string {
	return []string{
		"ABCDE12345",
	}
}

func (c *publishCommand) Run(args []string) error {
	if len(args) != 1 {
		return errors.New("no pin specified")
	}

	pin := args[0]

	service, err := GetCatalogsService()
	if err != nil {
		return err
	}

	// Start publish
	_, err = service.Publish().PIN(pin).Do(context.Background())
	if err != nil {
		return err
	}

	// Get status every 5 seconds
	for {
		time.Sleep(5 * time.Second)

		status, err := service.PublishStatus().PIN(pin).Do(context.Background())
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stdout, "Step %6d of %6d   %03d%%\r",
			status.CurrentStep, status.TotalSteps, status.Percent)

		if status.Done {
			break
		}
	}

	fmt.Fprintf(os.Stdout, "%s\rDone\n", strings.Repeat(" ", 78))

	return nil
}
