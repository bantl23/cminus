package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"strings"
)

func main() {
	parse := true
	analyze := true
	code := true
	echo := false
	trace := false

	app := cli.NewApp()
	app.Name = "cminus"
	app.Usage = "cminus [flags] <filename0>...<filenameN>"
	app.Version = "1.0.0-alpha0"
	app.Flags = []cli.Flag{
		cli.BoolTFlag{
			Name:        "parse",
			Usage:       "Enable or disable code parsing",
			Destination: &parse,
		},
		cli.BoolTFlag{
			Name:        "analyze",
			Usage:       "Enable or disable code analysis",
			Destination: &analyze,
		},
		cli.BoolTFlag{
			Name:        "code",
			Usage:       "Enable or disable code generation",
			Destination: &code,
		},
		cli.BoolTFlag{
			Name:        "trace",
			Usage:       "Turn on code tracing",
			Destination: &trace,
		},
		cli.BoolTFlag{
			Name:        "echo",
			Usage:       "Print source code",
			Destination: &echo,
		},
	}
	app.Action = func(c *cli.Context) {
		if len(c.Args()) == 0 {
			fmt.Println("Error: must supply filename to compile")
			os.Exit(1)
		}
		if analyze == false {
			code = false
		}
		if parse == false {
			analyze = false
			code = false
		}

		fmt.Printf("options: [parse=%t, analyze=%t, code=%t, echo=%t, trace=%t]\n",
			parse, analyze, code, echo, trace)

		for _, ifilename := range c.Args() {
			if strings.HasSuffix(ifilename, ".cm") == false {
				ifilename = ifilename + ".cm"
			}
			ofilename := strings.TrimSuffix(ifilename, ".cm") + ".tm"

			fmt.Println("compiling", ifilename)
			ifile, err := os.Open(ifilename)
			if err == nil {
				yyParse(NewLexer(ifile))

				fmt.Println("\nScanning")
				fmt.Println("========")

				if parse == true {
					fmt.Println("\nParsing")
					fmt.Println("=======")
					if analyze == true {
						fmt.Println("\nAnalyzing")
						fmt.Println("=========")
						if code == true {
							fmt.Println("\nCode Generation")
							fmt.Println("===============")
							fmt.Println(ofilename)
						}
					}
				}
			} else {
				fmt.Println("Error opening file", err)
			}
		}
	}
	app.Run(os.Args)
}
