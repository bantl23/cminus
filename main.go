package main

import (
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/syntree"
	"github.com/codegangsta/cli"
	"os"
	"strings"
)

func main() {
	parse := true
	analyze := true
	code := true
	echo := false
	print_parse_tree := false
	trace_scan := false
	trace_parse := false
	trace_analyze := false
	trace_codegen := false

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
			Name:        "trace-scan",
			Usage:       "Turn on tracing for scanning phase",
			Destination: &trace_scan,
		},
		cli.BoolFlag{
			Name:        "trace-parse",
			Usage:       "Turn on tracing for parsing phase",
			Destination: &trace_parse,
		},
		cli.BoolFlag{
			Name:        "trace-analyze",
			Usage:       "Turn on tracing for analysis phase",
			Destination: &trace_analyze,
		},
		cli.BoolFlag{
			Name:        "trace-codegen",
			Usage:       "Turn on tracing for code generation phase",
			Destination: &trace_codegen,
		},
		cli.BoolFlag{
			Name:        "echo",
			Usage:       "Print source code",
			Destination: &echo,
		},
		cli.BoolFlag{
			Name:        "print-parse-tree",
			Usage:       "Prints parse tree",
			Destination: &print_parse_tree,
		},
	}
	app.Action = func(c *cli.Context) {

		log.InfoLog = log.InitLog(true)
		log.ErrorLog = log.InitLog(true)
		log.EchoLog = log.InitLog(echo)
		log.ScanLog = log.InitLog(trace_scan)
		log.ParseLog = log.InitLog(trace_parse)
		log.AnalyzeLog = log.InitLog(trace_analyze)
		log.CodeLog = log.InitLog(trace_codegen)

		if len(c.Args()) == 0 {
			log.ErrorLog.Println("Must supply filename(s)")
			os.Exit(1)
		}
		if analyze == false {
			code = false
		}
		if parse == false {
			analyze = false
			code = false
		}

		for _, ifilename := range c.Args() {
			if strings.HasSuffix(ifilename, ".cm") == false {
				ifilename = ifilename + ".cm"
			}
			ofilename := strings.TrimSuffix(ifilename, ".cm") + ".tm"

			log.InfoLog.Printf("compiling %s\n", ifilename)
			ifile, err := os.Open(ifilename)
			if err == nil {
				log.InfoLog.Printf("scanning\n")
				if parse == true {
					log.InfoLog.Printf("parsing\n")
					yyParse(NewLexer(ifile))
					if print_parse_tree == true {
						syntree.Print(root, 0)
					}
					if analyze == true {
						log.InfoLog.Printf("analyzing\n")
						if code == true {
							log.InfoLog.Printf("code generation\n")
							log.InfoLog.Printf("creating %s\n", ofilename)
						}
					}
				}
			} else {
				log.ErrorLog.Printf("File open %s\n", err)
			}
		}
	}
	app.Run(os.Args)
}
