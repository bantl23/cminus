package main

import (
	"github.com/bantl23/cminus/log"
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
		cli.BoolFlag{
			Name:        "trace",
			Usage:       "Turn on code tracing",
			Destination: &trace,
		},
		cli.BoolFlag{
			Name:        "echo",
			Usage:       "Print source code",
			Destination: &echo,
		},
	}
	app.Action = func(c *cli.Context) {

		log.InitLogger(trace, echo)

		if len(c.Args()) == 0 {
			log.Error.Println("Must supply filename(s)")
			os.Exit(1)
		}
		if analyze == false {
			code = false
		}
		if parse == false {
			analyze = false
			code = false
		}

		log.Trace.Printf("[parse=%t, analyze=%t, code=%t, echo=%t, trace=%t]\n",
			parse, analyze, code, echo, trace)

		for _, ifilename := range c.Args() {
			if strings.HasSuffix(ifilename, ".cm") == false {
				ifilename = ifilename + ".cm"
			}
			ofilename := strings.TrimSuffix(ifilename, ".cm") + ".tm"

			log.Trace.Printf("compiling %s\n", ifilename)
			ifile, err := os.Open(ifilename)
			if err == nil {
				yyParse(NewLexer(ifile))

				log.Trace.Printf("scanning\n")
				if parse == true {
					log.Trace.Printf("parsing\n")
					if analyze == true {
						log.Trace.Printf("analyzing\n")
						if code == true {
							log.Trace.Printf("code generation\n")
							log.Trace.Printf("creating %s\n", ofilename)
						}
					}
				}
			} else {
				log.Error.Printf("File open %s\n", err)
			}
		}
	}
	app.Run(os.Args)
}
