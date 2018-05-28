package httphandler

import (
	. "http/comm"
)


func Handler1(handler HandlerHead)(map[string]interface{} ,error){
	querys := handler.Query

	r := handler.Request

	Debug(handler, "default page:%v\n", querys)

	Debug(handler,"test:%s", r.FormValue("test"))

	resMap := make(map[string]interface{})
	resMap["ret"]=0
	resMap["msg"]="handler1 ok"

	return resMap,nil
}
