package main

import(
	"net/http"
	"log"
)

import (
	. "http/comm"
	"strconv"
	. "http/handlebase"
)

var logger *log.Logger
var handlerMap  map[string]interface{}
var funcs Funcs

var visitTimes int


func httpHandlerInit(){

	funcs = NewFuncs(100)
	handlerMap = HandlerMapConfig()
	FuncBind(handlerMap, funcs)
}


//DefaultRequestHandler
func defaultHandler(w http.ResponseWriter,r *http.Request){

	logger.Println("request over\n")

	seq := GetRandUint32()
	query := r.URL.Query()
	handler := HandlerHead{logger, w,r, query, seq}

	defer panicHandler(handler)

	var handlePath string
	if len(r.URL.Path[1:]) > 0 {
		handlePath = r.URL.Path[1:]
	}

	if _,exist := handlerMap[handlePath]; !exist{
		handlePath = "default"
	}



	go funcs.Call(handlePath, handler)

	Debug(handler, "original path:%s,handlePath:%s,visitTimes:%s,req_map:%v",
		r.URL.Path,handlePath,strconv.Itoa(visitTimes),query)

	visitTimes++
}


func panicHandler(handler HandlerHead) {

	if err := recover(); err != nil {

		Debug(handler,"recover msg: ", err)

	} else {

		Debug(handler,"request over\n")

	}
}

func main(){

	visitTimes = 1;

	httpHandlerInit()

	var err error
	logger,err = CreateLog("")

	http.HandleFunc("/",defaultHandler)

	DebugBase("start svr ok\n\n\n\n")

	err = http.ListenAndServe(":8080",nil)

	if err != nil {
		DebugBase("http svr err:%v", err)
	}
}


