package main

import(
	"net/http"
	"log"
)

import (
	. "http/comm"
	"strconv"
	. "http/handlebase"
	"fmt"
	"encoding/json"
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

	Debug(handler, "original path:%s,handlePath:%s,visitTimes:%s,req_map:%v",
		r.URL.Path,handlePath,strconv.Itoa(visitTimes),query)

	callHandler(handlePath, handler)

	visitTimes++
}

func callHandler(handlePath string, handler HandlerHead){
	resValue, err := funcs.Call(handlePath, handler)

	resMap := resValue[0].Interface().(map[string]interface{})

	data, _ := json.Marshal(resMap)

	res :=string(data)

	Debug(handler,"resmap:%v,err:%v,res:%s", resMap, err, res)
	fmt.Fprintf(handler.Writer,"%s", res)

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


