package httphandler

import (
	. "http/comm"
)


func HandlerDefault(handler Handler)(map[string]interface{}){


	Debug(handler,"test:%s", handler.Request.FormValue("test"))

	resMap := make(map[string]interface{})
	resMap["ret"]=999
	resMap["msg"]="err interface"

	return resMap
}
