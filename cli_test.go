package main

import (
	"os"
	"strings"
	"testing"
)

func Test_run(t *testing.T) {
	t.Parallel()
	cli := &cli{}
	t.Run("invalid arg case", func(t *testing.T) {
		testCases := []struct {
			name   string
			argStr string
		}{
			{
				name:   "no args",
				argStr: "",
			},
			{
				name:   "no option",
				argStr: "schema.sql",
			},
		}
		for _, tc := range testCases {
			args := strings.Split(tc.argStr, " ")
			if got, want := cli.run(args), exitCodeArgsError; got != want {
				t.Fatalf("%s exits %d, want %d\n", tc.name, got, want)
			}
		}
	})

	t.Run("error case", func(t *testing.T) {
		testCases := []struct {
			name   string
			argStr string
		}{
			{
				name:   "file does not exist",
				argStr: "-s not_exist.sql",
			},
		}
		for _, tc := range testCases {
			args := strings.Split(tc.argStr, " ")
			if got, want := cli.run(args), exitCodeError; got != want {
				t.Fatalf("%s exits %d, want %d\n", tc.name, got, want)
			}
		}
	})
	t.Run("valid case", func(t *testing.T) {
		output := "testdata/sample.png"
		if got, want := cli.run([]string{"-s", "testdata/sample.sql", "-o", output}), exitCodeOK; got != want {
			t.Fatalf("valid case exits %d, want %d\n", got, want)
		}
		_, err := os.Stat(output)
		if err != nil {
			t.Fatalf("cant create output file. %s", err)
		}
		if err := os.Remove(output); err != nil {
			t.Fatalf("remove test file failed. %s", err)
		}
	})
}
