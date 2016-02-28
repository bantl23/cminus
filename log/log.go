package log

import (
	"io/ioutil"
	"log"
	"os"
)

var EchoLog *log.Logger
var InfoLog *log.Logger
var ErrorLog *log.Logger
var ScanLog *log.Logger
var ParseLog *log.Logger
var AnalyzeLog *log.Logger
var CodeLog *log.Logger

func InitLog(on bool) *log.Logger {
	out := ioutil.Discard
	if on == true {
		out = os.Stdout
	}
	return log.New(out, "", log.Lshortfile)
}
