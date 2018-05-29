package httphandler

import (
	. "http/comm"
)


func Handler1(handler Handler)(map[string]interface{}){

	r := handler.Request

	Debug(handler,"test:%s", r.FormValue("test"))

	resMap := make(map[string]interface{})
	resMap["ret"]=0
	resMap["msg"]="handler1 ok"

	return resMap
}
