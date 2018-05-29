package httphandler

import (
	. "http/comm"
)


func Upload(handler Handler)(map[string]interface{}){

	queryParams := handler.QueryParams

	Debug(handler,"test:%s,k:%s",
		GetValue(queryParams, "test"), GetValue(queryParams,"k"))

	resMap := make(map[string]interface{})
	resMap["ret"]=0
	resMap["msg"]="handler1 ok"

	return resMap
}
