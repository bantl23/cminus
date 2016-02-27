package log

import (
	"io/ioutil"
	"log"
	"os"
)

var Trace *log.Logger
var Error *log.Logger
var Echo *log.Logger

func InitLogger(trace bool, echo bool) {
	traceOutput := ioutil.Discard
	if trace == true {
		traceOutput = os.Stdout
	}
	echoOutput := ioutil.Discard
	if echo == true {
		echoOutput = os.Stdout
	}
	errorOutput := os.Stderr

	Trace = log.New(traceOutput, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorOutput, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	Echo = log.New(echoOutput, "ECHO: ", log.Ldate|log.Ltime|log.Lshortfile)
}
