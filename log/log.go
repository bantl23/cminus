package log

import (
	"io/ioutil"
	"log"
	"os"
)

var Trace *log.Logger
var Error *log.Logger
var Echo *log.Logger
var Info *log.Logger

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
	infoOutput := os.Stdout

	Trace = log.New(traceOutput, "TRACE: ", log.Lshortfile)
	Error = log.New(errorOutput, "ERROR: ", log.Lshortfile)
	Echo = log.New(echoOutput, "ECHO: ", log.Lshortfile)
	Info = log.New(infoOutput, "INFO: ", log.Lshortfile)
}
