package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/ast"
	"github.com/cloudspannerecosystem/memefish/token"
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
		help   bool
		file   string
		output string
		t      string
	)
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	flags.BoolVar(&help, "h", false, "print help")
	flags.StringVar(&file, "s", "", "spanner schema file")
	flags.StringVar(&output, "o", "", "output file name.default is spanner_er.<type>(pass to dot option -o)")
	flags.StringVar(&t, "T", "png", "output file type. default is png(pass to dot option -T)")
	if err := flags.Parse(args); err != nil {
		return exitCodeArgsError
	}
	if help {
		flags.Usage()
		return exitCodeOK
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
	log.Printf("Parsed %d tables", len(tables))
	for i, t := range tables {
		log.Printf("Table %d: %s", i, getTableName(t))
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
	log.Printf("Graph DOT content: %s", s)
	r := strings.NewReader(s)
	c := exec.Command("dot", fmt.Sprintf("-T%s", t), "-o", output)
	c.Stdin = r
	var stderr strings.Builder
	c.Stderr = &stderr
	err = c.Run()
	if err != nil {
		log.Printf("Error running dot command: %v", err)
		log.Printf("Stderr: %s", stderr.String())
		return exitCodeError
	}
	log.Printf("Output file created: %s", output)

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

func parse(sqls string) ([]*ast.CreateTable, error) {
	// Log the original SQL for debugging
	log.Printf("Original SQL: %s", sqls)

	// Split the SQL by semicolons to get individual statements
	statements := strings.Split(sqls, ";")

	var tables []*ast.CreateTable
	for _, stmt := range statements {
		// Skip empty statements
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		// Log each statement for debugging
		log.Printf("Parsing statement: %s", stmt)

		// Create a new Parser instance for each statement
		file := &token.File{
			Buffer: stmt,
		}
		p := &memefish.Parser{
			Lexer: &memefish.Lexer{File: file},
		}

		// Parse the statement
		parsedStmt, err := p.ParseStatement()
		if err != nil {
			log.Printf("Error parsing statement: %v", err)
			continue
		}

		// If it's a CREATE TABLE statement, add it to our list
		if createTable, ok := parsedStmt.(*ast.CreateTable); ok {
			log.Printf("Found CREATE TABLE: %s", getTableName(createTable))
			tables = append(tables, createTable)
		} else {
			log.Printf("Statement is not a CREATE TABLE: %T", parsedStmt)
		}
	}

	return tables, nil
}

// Helper function to get the name from a CreateTable
func getTableName(t *ast.CreateTable) string {
	if t.Name != nil && len(t.Name.Idents) > 0 {
		return t.Name.Idents[len(t.Name.Idents)-1].Name
	}
	return ""
}
