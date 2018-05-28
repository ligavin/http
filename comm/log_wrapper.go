package comm

import (
	"log"
	"os"
	"fmt"
	"runtime"
	"strings"
)

var globalLogger *log.Logger

func CreateLog(fileName string)(*log.Logger, error){

	if (fileName == fileName){
		fileName = "go_log.txt"
	}

	logFile,err  := os.OpenFile(fileName, os.O_RDWR| os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalln("open file(%s) error !", fileName)

	}
	logger := log.New(logFile,"",log.Lshortfile|log.LstdFlags)

	globalLogger = logger

	return logger, err
}

func Debug(handler HandlerHead, format string, v ...interface{}) {
	format = fmt.Sprintf("seq=%d,%s", handler.Seq, format)
	OutputLog(globalLogger,3,  format, v...)
}

func DebugBase(format string, v ...interface{}) {
	OutputLog(globalLogger, 3, format, v...)
}

func OutputLog(logger *log.Logger, calldepth int, format string, v ...interface{}){


	pc, _, _, ok := runtime.Caller(calldepth-1)
	fn := "???"

	if ok{
		f := runtime.FuncForPC(pc)
		fn = f.Name()
	}

	if strings.Contains(fn, ".") {
		callers := strings.Split(fn, ".")
		fn = callers[len(callers)-1]
	}

	format = fmt.Sprintf("caller:%s,%s", fn, format)

	logger.Output(calldepth, fmt.Sprintf(format, v...))
}