package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"cloud.google.com/go/spanner/spansql"
)

const (
	defaultFileBasename     = "spanner_er"
	exitCodeOK          int = 0
	exitCodeError           = 10 + iota
	exitCodeArgsError
)

type cli struct{}

func (cli *cli) run(args []string) int {
	var (
		file   string
		output string
		t      string
	)
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	flags.StringVar(&file, "s", "", "spanner schema file")
	flags.StringVar(&output, "o", "", "output file name.default is spanner_er.<type>(pass to dot option -o)")
	flags.StringVar(&t, "T", "png", "output file type. default is png(pass to dot option -T)")
	if err := flags.Parse(args); err != nil {
		flags.Usage()
		return exitCodeArgsError
	}
	if file == "" {
		flags.Usage()
		return exitCodeArgsError
	}
	if output == "" {
		output = fmt.Sprintf("%s.%s", defaultFileBasename, t)
	}
	body, err := cli.read(file)
	if err != nil {
		log.Print(err)
		return exitCodeError
	}
	tables, err := parse(body)
	if err != nil {
		log.Print(err)
		return exitCodeError
	}
	graph, err := NewGraph()
	if err != nil {
		log.Print(err)
		return exitCodeError
	}

	if err := graph.ApplyTables(tables); err != nil {
		log.Print(err)
		return exitCodeError
	}
	s := graph.String()
	r := strings.NewReader(s)
	c := exec.Command("dot", fmt.Sprintf("-T%s", t), "-o", output)
	c.Stdin = r
	c.Start()
	c.Wait()

	return exitCodeOK
}

func (cli *cli) read(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	body := string(data)
	return body, nil

}

func parse(sqls string) ([]spansql.CreateTable, error) {
	// spansql not allow backquote
	sqls = strings.Replace(sqls, "`", "", -1)
	d, err := spansql.ParseDDL(sqls)
	if err != nil {
		return nil, err
	}
	tables := []spansql.CreateTable{}
	for _, e := range d.List {
		switch v := e.(type) {
		case spansql.CreateTable:
			tables = append(tables, v)
		}
	}
	return tables, nil
}
