package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
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
	fmt.Printf("%3s  %-50s %-10s %-10s\n", "ID", "Name", "Created", "PIN")
	fmt.Printf("%s\n", strings.Repeat("=", 78))
	for _, cat := range res.Items {
		fmt.Printf("%3d. %-50s %-10s %-10s\n", cat.ID, substring(cat.Name, 50), cat.Created.Format("2006-01-02"), cat.PIN)
	}

	return nil
}

func substring(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	var b bytes.Buffer
	for i, r := range s {
		if i >= n {
			break
		}
		b.WriteRune(r)
	}
	return b.String()
}
