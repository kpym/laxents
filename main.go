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
	"errors"
	"fmt"
	"os"

	"github.com/kpym/laxents/api"
	"github.com/kpym/laxents/parameters"
)

func main() {
	var (
		err    error
		params *parameters.Parameters
	)
	params, err = parameters.Get(os.Args[1:])
	if err != nil {
		if errors.Is(err, parameters.ErrMissingParameters) {
			parameters.PrintHelp()
			os.Exit(0)
		}
		fmt.Println(err)
		os.Exit(1)
	}
	defer params.In.Close()
	defer params.Out.Close()

	// convert the input
	if params.ToUnicode {
		err = api.ToUnicode(params.Out, params.In)
	} else {
		err = api.ToLaTeX(params.Out, params.In)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
