// this is a small program that convert between LaTeX accents and Unicode characters
// it is a simple program that uses the `strings` package to convert between the two
// usage :
// > laxents -to-unicode <string>
// > laxents -to-latex <string>
// > laxents -to-unicode -i <file> -o <file>
// > laxents -to-latex -i <file> -o <file>
// if no input file or string is provided, it will read from stdin
// if no output file is provided, it will write to stdout
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
)

// check is a helper function to check for errors
func check(e error, msg string) {
	if e != nil {
		fmt.Fprintln(os.Stderr, msg)
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}

// the parameters for the program
type Args struct {
	ToUnicode bool   `arg:"-u,--to-unicode" help:"convert from LaTeX to Unicode"`
	ToLatex   bool   `arg:"-l,--to-latex" help:"convert from Unicode to LaTeX"`
	Input     string `arg:"-i,--input" help:"input file"`
	Output    string `arg:"-o,--output" help:"output file"`
	Text      string `arg:"positional" help:"string to convert"`
}

func (Args) Description() string {
	return "convert between LaTeX accents and Unicode characters"
}

func (Args) Epilog() string {
	// get the executable name without the path
	exe := filepath.Base(os.Args[0])
	// remove the .exe extension on Windows
	if filepath.Ext(exe) == ".exe" {
		exe = exe[:len(exe)-4]
	}
	return fmt.Sprintf(`Examples:
	%s -to-latex "déçû"
	%s -to-unicode "d\\'e\\c{c}\\^{u}"
	%s -to-unicode -i input.tex -o output.tex
	%s -to-latex -i input.tex -o output.tex
	cat input.tex | %s -to-unicode
	`, exe, exe, exe, exe, exe)
}

// main function
func main() {
	var args Args
	p, err := arg.NewParser(arg.Config{}, &args)
	check(err, "cannot create parser")
	p.Parse(os.Args[1:])

	// get the replacer
	var rep *strings.Replacer
	if args.ToUnicode && args.ToLatex {
		check(fmt.Errorf("cannot specify both -to-unicode and -to-latex"), "")
	}
	if !args.ToUnicode && !args.ToLatex {
		check(fmt.Errorf("must specify either -to-unicode or -to-latex"), "")
	}
	if args.ToUnicode {
		rep = toUnicode()
	} else {
		rep = toLatex()
	}

	// get the input
	var in io.Reader
	if args.Input != "" && args.Text != "" {
		check(fmt.Errorf("cannot specify both a file and a string"), "")
	}
	if args.Input != "" {
		f, err := os.Open(args.Input)
		check(err, "cannot open input file")
		defer f.Close()
		in = f
	} else if args.Text != "" {
		in = strings.NewReader(args.Text)
	} else {
		// check if there is data on stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			in = os.Stdin
		}
	}
	if in == nil {
		// display the usage message if there is no input
		p.WriteUsage(os.Stderr)
		os.Exit(1)
	}
	in = utfReader(in)

	// get the output
	var out io.Writer
	if args.Output == "" {
		out = os.Stdout
	} else {
		f, err := os.Create(args.Output)
		check(err, "cannot create output file")
		defer f.Close()
		out = f
	}

	// convert the input
	err = CopyWithReplace(in, out, rep)
	check(err, "problem during conversion")
}
