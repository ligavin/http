package httphandler

import (
	"fmt"
)
import . "http/comm"


func HandlerDefault(handler HandlerHead){
	w := handler.Writer
	r := handler.Request
	querys := handler.Query
	path:=r.URL.Path[1:]

	fmt.Fprintf(w,"<head><title>404</title></head><h1>Hello,%s!</h1>", path)
	Debug(handler, "default page:%v\n", querys)

	Debug(handler,"test:%s", GetValue(querys, "test"))

}
