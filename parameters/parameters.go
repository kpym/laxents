// parameters is a package that handles the command line parameters
package parameters

import (
	"errors"
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
		panic(fmt.Errorf("%s: %v", msg, e))
	}
}

// nopWriteCloser is a WriteCloser that does nothing on Close
type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }

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

func (Args) Epilogue() string {
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

// PrintHelp prints the help message
func PrintHelp() {
	p := arg.MustParse(&Args{})
	p.WriteHelp(os.Stdout)
}

// Parameters for the program
// returned by the Get function
type Parameters struct {
	ToUnicode bool
	In        io.ReadCloser
	Out       io.WriteCloser
}

var (
	// the missing input error
	ErrMissingParameters = errors.New("missing parameters")
)

// Get parses the command line arguments and returns the parameters
func Get(osArgs []string) (params *Parameters, err error) {
	// check if there are arguments
	if len(osArgs) == 0 {
		return nil, ErrMissingParameters
	}

	// capture the panic and return it as an error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	var args Args
	params = &Parameters{}

	p, err := arg.NewParser(arg.Config{}, &args)
	check(err, "cannot create parser")
	p.MustParse(osArgs)

	// get the replacer
	if args.ToUnicode && args.ToLatex {
		return nil, errors.New("cannot specify both -to-unicode and -to-latex")
	}
	if !args.ToUnicode && !args.ToLatex {
		return nil, errors.New("must specify either -to-unicode or -to-latex")
	}
	params.ToUnicode = args.ToUnicode

	// get the input
	if args.Input != "" && args.Text != "" {
		return nil, errors.New("cannot specify both a file and a string")
	}
	if args.Input != "" {
		f, err := os.Open(args.Input)
		check(err, "cannot open input file")
		params.In = f
	} else if args.Text != "" {
		params.In = io.NopCloser(strings.NewReader(args.Text))
	} else {
		// check if there is data on stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			params.In = io.NopCloser(os.Stdin)
		}
	}
	if params.In == nil {
		return nil, errors.New("no input data provided")
	}

	// get the output
	if args.Output == "" {
		params.Out = nopWriteCloser{os.Stdout}
	} else {
		f, err := os.Create(args.Output)
		check(err, "cannot create output file")
		params.Out = f
	}

	return params, nil
}
