package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// catalogsCommand lists your catalogs.
type catalogsCommand struct {
	take, skip int64
	sort       string
}

func init() {
	RegisterCommand("catalogs", func(flags *flag.FlagSet) Command {
		cmd := new(catalogsCommand)
		flags.Int64Var(&cmd.take, "take", 0, "Number of catalogs to take")
		flags.Int64Var(&cmd.skip, "skip", 0, "Number of catalogs to skip")
		flags.StringVar(&cmd.sort, "sort", "", "Sort order, e.g. name or id or -created")
		return cmd
	})
}

func (c *catalogsCommand) Describe() string {
	return "List catalogs."
}

func (c *catalogsCommand) Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s catalogs\n", os.Args[0])
}

func (c *catalogsCommand) Examples() []string {
	return []string{
		"-take=5",
		"-take=5 -skip=5",
		"-sort=-created,id",
	}
}

func (c *catalogsCommand) Run(args []string) error {
	service, err := GetCatalogsService()
	if err != nil {
		return err
	}

	svc := service.Search()
	if c.skip > 0 {
		svc = svc.Skip(c.skip)
	}
	if c.take > 0 {
		svc = svc.Take(c.take)
	}
	svc = svc.Sort(c.sort)

	res, err := svc.Do()
	if err != nil {
		return err
	}

	fmt.Printf("%d catalogs found.\n", res.TotalItems)
	fmt.Printf("%3s  %-50s %-10s\n", "ID", "Name", "Created")
	fmt.Printf("%s\n", strings.Repeat("=", 78))
	for _, cat := range res.Items {
		fmt.Printf("%3d. %-50s %-10s\n", cat.ID, cat.Name, cat.Created.Format("2006-01-02"))
	}

	return nil
}
