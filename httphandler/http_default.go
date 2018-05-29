package httphandler

import (
	. "http/comm"
)


func HandlerDefault(handler HandlerHead)(map[string]interface{}){
	querys := handler.Query

	Debug(handler, "default page:%v\n", querys)

	Debug(handler,"test:%s", GetValue(querys, "test"))

	resMap := make(map[string]interface{})
	resMap["ret"]=999
	resMap["msg"]="err interface"

	return resMap
}
