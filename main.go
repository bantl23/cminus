package main

import (
	"fmt"
	"github.com/bantl23/cminus/gen"
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/opt"
	"github.com/bantl23/cminus/symtbl"
	"github.com/bantl23/cminus/syntree"
	"github.com/codegangsta/cli"
	"os"
	"strings"
)

func main() {
	parse := true
	analyze := true
	optimize := true
	code := true
	print_source_code := false
	print_parse_tree := false
	print_symbol_table := false
	print_symbol_map := false
	print_machine_code := false
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
			Name:        "optimize",
			Usage:       "Enable or disable code optimization",
			Destination: &optimize,
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
			Name:        "print-source-code",
			Usage:       "Print source code",
			Destination: &print_source_code,
		},
		cli.BoolFlag{
			Name:        "print-parse-tree",
			Usage:       "Prints parse tree",
			Destination: &print_parse_tree,
		},
		cli.BoolFlag{
			Name:        "print-symbol-table",
			Usage:       "Prints symbol table",
			Destination: &print_symbol_table,
		},
		cli.BoolFlag{
			Name:        "print-symbol-map",
			Usage:       "Prints symbol map",
			Destination: &print_symbol_map,
		},
		cli.BoolFlag{
			Name:        "print-machine-code",
			Usage:       "Print machine code",
			Destination: &print_machine_code,
		},
	}
	app.Action = func(c *cli.Context) {

		log.InfoLog = log.InitLog(true)
		log.ErrorLog = log.InitLog(true)
		log.SrcLog = log.InitLog(print_source_code)
		log.DstLog = log.InitLog(print_machine_code)
		log.ScanLog = log.InitLog(trace_scan)
		log.ParseLog = log.InitLog(trace_parse)
		log.AnalyzeLog = log.InitLog(trace_analyze)
		log.CodeLog = log.InitLog(trace_codegen)

		if len(c.Args()) == 0 {
			log.ErrorLog.Println(">>>>> Error must supply filename(s)")
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

			log.InfoLog.Printf("compiling %s", ifilename)
			ifile, err := os.Open(ifilename)
			if err == nil {
				if parse == true {
					log.InfoLog.Printf("scanning and parsing")
					log.InfoLog.Printf("====================")
					yyParse(NewLexer(ifile))
					if print_parse_tree == true {
						log.InfoLog.Printf("parse tree")
						log.InfoLog.Printf("==========")
						fmt.Println(">>>>")
						syntree.PrintNode(rootNode, 0)
						fmt.Println("<<<<")
					}
					if analyze == true {
						log.InfoLog.Printf("building symbol table")
						log.InfoLog.Printf("=====================")
						symtbl.BuildTableList(rootNode)
						if print_symbol_table == true {
							log.InfoLog.Printf("symbol table")
							log.InfoLog.Printf("============")
							fmt.Println(">>>>")
							symtbl.GlbPrintTableList()
							fmt.Println("<<<<")
						}
						if print_symbol_map == true {
							fmt.Println("<<<<")
							log.InfoLog.Printf("symbol map")
							log.InfoLog.Printf("==========")
							fmt.Println(">>>>")
							symtbl.GlbPrintTableMap()
							fmt.Println("<<<<")
							if print_parse_tree == true {
								log.InfoLog.Printf("parse tree with symbol map keys")
								log.InfoLog.Printf("===============================")
								fmt.Println("<<<<")
								syntree.PrintNodeWithSymKey(rootNode, 0)
								fmt.Println(">>>>")
							}
						}
						log.InfoLog.Printf("analyzing")
						log.InfoLog.Printf("=========")
						symtbl.Analyze(rootNode)
						if optimize == true {
							log.InfoLog.Printf("optimzing")
							opt.Optimize(rootNode)
						}
						if code == true {
							log.InfoLog.Printf("generating code")
							log.InfoLog.Printf("===============")
							log.InfoLog.Printf("creating %s", ofilename)
							gen.Generate(rootNode, ofilename)
						}
					}
				}
			} else {
				log.ErrorLog.Printf(">>>>> Error file open %s", err)
			}
		}
	}
	app.Run(os.Args)
}
