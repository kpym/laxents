# laxents
Accents translater between UTF-8 and LaTeX.

## Usage

```bash
$ laxents.exe --help
convert between LaTeX accents and Unicode characters
Usage: laxents.exe [--to-unicode] [--to-latex] [--input INPUT] [--output OUTPUT] [TEXT]

Positional arguments:
  TEXT                   string to convert

Options:
  --to-unicode, -u       convert from LaTeX to Unicode
  --to-latex, -l         convert from Unicode to LaTeX
  --input INPUT, -i INPUT
                         input file
  --output OUTPUT, -o OUTPUT
                         output file
  --help, -h             display this help and exit

Examples:
        laxents -to-latex "déçû"
        laxents -to-unicode "d\\'e\\c{c}\\^{u}"
        laxents -to-unicode -i input.tex -o output.tex
        laxents -to-latex -i input.tex -o output.tex
        cat input.tex | laxents -to-unicode
```

## Installation

Dowload it from the [releases page](https://github.com/kpym/esplus/releases) and put it in your path.
Or build it yourself:

```bash
go install github.com/kpym/esplus@latest
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

